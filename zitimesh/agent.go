package zitimesh

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gage-technologies/gigo-lib/logging"
	"github.com/openziti/sdk-golang/ziti/edge"
	"github.com/sourcegraph/conc"
	"io"
	"net"
	"sync"
	"time"

	"github.com/openziti/edge-api/rest_model"
	"github.com/openziti/sdk-golang/ziti"
)

// ForwardContext
//
// Context used for forwarding a port on the local network to the Ziti network.
type ForwardContext struct {
	// Ctx is the context used to manage the port forwarding
	Ctx context.Context

	// Cancel is the function used to cancel the port forwarding
	Cancel context.CancelFunc

	// Service is the name of the Ziti service that is being forwarded
	Service string

	// Done is the channel that is closed when the port forwarding is complete
	Done chan struct{}
}

// Agent
//
// An agent is a Ziti network binder that exposes the local node ports to the Ziti network.
type Agent struct {
	// Name if name of the agent on the ziti network
	Name string

	// token is the JWT token used to authenticate the agent with the Ziti controller
	token string

	// ctx is the context used to manage the agent
	ctx context.Context

	// cancel is the function used to cancel the agent
	cancel context.CancelFunc

	// zitiCtx is the Ziti context used to connect to the Ziti network
	zitiCtx ziti.Context

	// forwardCtxs is a map of port forwarding contexts
	forwardCtxs map[string]*ForwardContext

	// mu is the mutex used to synchronize access to the forwarding contexts
	mu *sync.Mutex

	// removeListeners is a list of listeners that are used to callbacks from the Ziti context
	removeListeners []func()

	// wg is the wait group used to wait for the agent to close
	wg *conc.WaitGroup

	// logger is the logger used by the agent
	logger logging.Logger
}

// NewAgent
//
// Create a new Ziti agent.
func NewAgent(ctx context.Context, name string, token string, logger logging.Logger) (*Agent, error) {
	// enroll the identity into a configuration
	zitiConfig, err := EnrollIdentity(token)
	if err != nil {
		return nil, fmt.Errorf("failed to enroll identity: %w", err)
	}

	// create a new Ziti context
	zitiCtx, err := ziti.NewContext(zitiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ziti context: %w", err)
	}

	// create a new context with a cancel function
	ctx, cancel := context.WithCancel(ctx)

	// create the agent
	agent := &Agent{
		Name:        name,
		token:       token,
		ctx:         ctx,
		cancel:      cancel,
		zitiCtx:     zitiCtx,
		forwardCtxs: make(map[string]*ForwardContext),
		mu:          &sync.Mutex{},
		wg:          conc.NewWaitGroup(),
		logger:      logger,
	}

	zitiCtx.Events().AddServiceAddedListener(func(_ ziti.Context, event *rest_model.ServiceDetail) {
		if event == nil {
			return
		}
		agent.serviceAdded(*event)
	})

	// add the service watcher
	agent.wg.Go(agent.serviceWatcher)

	return agent, nil
}

func (a *Agent) Close() {
	// cancel the agent context
	a.cancel()

	// close the ziti context
	a.zitiCtx.Close()

	// wait for the agent to close
	a.wg.Wait()
}

func (a *Agent) serviceWatcher() {
	// loop until the context is cancelled
	for {
		// refresh the nodes services
		err := a.zitiCtx.RefreshServices()
		if err != nil {
			a.logger.Warnf("(ziti-agent: %s) failed to refresh ziti services: %v", a.Name, err)
			continue
		}

		// wait for the context to be cancelled
		select {
		case <-a.ctx.Done():
			return
		case <-time.After(time.Second):
		}
	}
}

func (a *Agent) serviceAdded(service rest_model.ServiceDetail) {
	// acquire the mutex
	a.mu.Lock()
	defer a.mu.Unlock()

	// check if there is an existing forwarding context for this port and network
	// if there is then we need to terminate the forwarding before we create a new one
	if fctx, ok := a.forwardCtxs[*service.Name]; ok {
		// kill the context and wait for it to exit
		fctx.Cancel()
		<-fctx.Done

		// remove the forwarding context from the list of active forwarding contexts
		delete(a.forwardCtxs, *service.Name)
	}

	// create a new forwarding context
	innerCtx, innerCancel := context.WithCancel(context.Background())
	fctx := &ForwardContext{
		Ctx:     innerCtx,
		Cancel:  innerCancel,
		Service: *service.Name,
		Done:    make(chan struct{}),
	}

	// add the forwarding context to the list of active forwarding contexts
	a.forwardCtxs[*service.Name] = fctx

	// start the port listener
	a.wg.Go(func() {
		a.startPortListener(*fctx)
	})
}

func (a *Agent) startPortListener(fctx ForwardContext) {
	// defer the closing of the done channel
	defer close(fctx.Done)

	// create a listener on the zitinet for the service
	listener, err := a.zitiCtx.ListenWithOptions(fctx.Service, &ziti.ListenOptions{
		ConnectTimeout: 5 * time.Minute,
		MaxConnections: 64,
	})
	if err != nil {
		a.logger.Warnf("(ziti-agent: %s) failed to listen on service %s: %v", a.Name, fctx.Service, err)
		return
	}
	defer listener.Close()

	// accept connections forever forwarding them to
	// the local port
	a.wg.Go(func() {
		for {
			zitiConn, err := listener.AcceptEdge()
			if err != nil {
				// exit if the listener is closed
				if listener.IsClosed() {
					return
				}
				a.logger.Warnf("(ziti-agent: %s) failed to accept connection: %v", a.Name, err)
				continue
			}

			a.wg.Go(func() {
				a.handleConnection(zitiConn, fctx)
			})
		}
	})

	// wait for the context to be cancelled
	<-fctx.Ctx.Done()
}

func (a *Agent) handleConnection(zitiConn edge.Conn, fctx ForwardContext) {
	defer zitiConn.Close()

	// load the local config from the connection
	buf := zitiConn.GetAppData()
	if len(buf) == 0 {
		a.logger.Warnf("(ziti-agent: %s) failed to read local config from connection", a.Name)
		return
	}
	var localConfig AgentService
	err := json.Unmarshal(buf, &localConfig)
	if err != nil {
		a.logger.Warnf("(ziti-agent: %s) failed to unmarshal local config: %v", a.Name, err)
		return
	}

	// dial the localport
	localConn, err := net.Dial(string(localConfig.Network), fmt.Sprintf("localhost:%d", localConfig.Port))
	if err != nil {
		a.logger.Warnf("(ziti-agent: %s) failed to dial local port: %v", a.Name, err)
		return
	}
	defer localConn.Close()

	// copy input from the zitinet to the local port
	a.wg.Go(func() {
		io.Copy(localConn, zitiConn)
	})

	// copy input from the local port to the zitinet
	a.wg.Go(func() {
		io.Copy(zitiConn, localConn)
	})

	// wait for the context to be cancelled
	<-fctx.Ctx.Done()
}
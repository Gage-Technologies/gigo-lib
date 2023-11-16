package zitimesh

import (
	"context"
	"fmt"
	"github.com/gage-technologies/gigo-lib/logging"
	"github.com/openziti/sdk-golang/ziti"
	"net"
)

// Server
//
// A server is a Ziti network dialer that can dial resource on any agent in the Ziti network.
type Server struct {
	// Name is the name of the server on the Ziti network
	Name string

	// token is the JWT token used to authenticate the server with the Ziti controller
	token string

	// ctx is the context used to manage the server
	ctx context.Context

	// cancel is the function used to cancel the server
	cancel context.CancelFunc

	// zitiCtx is the Ziti context used to connect to the Ziti network
	zitiCtx ziti.Context

	// logger is the logger used by the server
	logger logging.Logger
}

// NewServer
//
// Creates a new server instance.
func NewServer(ctx context.Context, name string, token string, logger logging.Logger) (*Server, error) {
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

	// create a new server instance
	server := &Server{
		Name:    name,
		token:   token,
		ctx:     ctx,
		cancel:  cancel,
		zitiCtx: zitiCtx,
		logger:  logger,
	}

	return server, nil
}

func (s *Server) DialAgent(agentId int64, network NetworkType, port int) (net.Conn, error) {
	// refresh the services
	_ = s.zitiCtx.RefreshServices()

	// dial the agent
	return s.zitiCtx.DialWithOptions(fmt.Sprintf("gigo-workspace-access-%d", agentId), &ziti.DialOptions{
		AppData: []byte(fmt.Sprintf(`{"network":"%s","port":%d}`, network, port)),
	})
}

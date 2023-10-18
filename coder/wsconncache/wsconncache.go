// Package wsconncache caches workspace agent connections by UUID.
package wsconncache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/gage-technologies/gigo-lib/coder/agentsdk"
	"github.com/gage-technologies/gigo-lib/logging"
	"github.com/gage-technologies/gigo-lib/mq"
	"github.com/gage-technologies/gigo-lib/mq/models"
	"github.com/gage-technologies/gigo-lib/mq/streams"
	"github.com/nats-io/nats.go"
	"net/http"
	"sync"
	"time"

	"go.uber.org/atomic"
	"golang.org/x/sync/singleflight"
	"golang.org/x/xerrors"
)

type CacheParams struct {
	// Dialer creates a new agent connection by ID.
	Dialer          Dialer
	InactiveTimeout time.Duration
	Logger          logging.Logger
	Js              *mq.JetstreamClient
}

// New creates a new workspace connection cache that closes
// connections after the inactive timeout provided.
//
// Agent connections are cached due to WebRTC negotiation
// taking a few hundred milliseconds.
func New(params CacheParams) (*Cache, error) {
	if params.InactiveTimeout == 0 {
		params.InactiveTimeout = 5 * time.Minute
	}

	cache := &Cache{
		closed:          make(chan struct{}),
		dialer:          params.Dialer,
		inactiveTimeout: params.InactiveTimeout,
		logger:          params.Logger,
		js:              params.Js,
	}

	subscription, err := params.Js.Subscribe(
		streams.SubjectWsConnCacheForget,
		cache.handleForgetMsg,
	)
	if err != nil {
		return nil, fmt.Errorf("could not subscribe to forget messages: %v", err)
	}

	cache.subscriber = subscription

	return cache, nil
}

// Dialer creates a new agent connection by ID.
type Dialer func(r *http.Request, id int64) (*agentsdk.AgentConn, error)

// Conn wraps an agent connection with a reusable HTTP transport.
type Conn struct {
	*agentsdk.AgentConn

	locks         atomic.Uint64
	timeoutMutex  sync.Mutex
	timeout       *time.Timer
	timeoutCancel context.CancelFunc
	transport     *http.Transport
}

func (c *Conn) HTTPTransport() *http.Transport {
	return c.transport
}

// CloseWithError ends the HTTP transport if exists, and closes the agent.
func (c *Conn) CloseWithError(err error) error {
	if c.transport != nil {
		c.transport.CloseIdleConnections()
	}
	c.timeoutMutex.Lock()
	defer c.timeoutMutex.Unlock()
	if c.timeout != nil {
		c.timeout.Stop()
	}
	return c.AgentConn.CloseWithError(err)
}

type Cache struct {
	closed          chan struct{}
	closeMutex      sync.Mutex
	closeGroup      sync.WaitGroup
	connGroup       singleflight.Group
	connMap         sync.Map
	dialer          Dialer
	inactiveTimeout time.Duration
	logger          logging.Logger
	js              *mq.JetstreamClient
	subscriber      *nats.Subscription
}

// Acquire gets or establishes a connection with the dialer using the ID provided.
// If a connection is in-progress, that connection or error will be returned.
//
// The returned function is used to release a lock on the connection. Once zero
// locks exist on a connection, the inactive timeout will begin to tick down.
// After the time expires, the connection will be cleared from the cache.
func (c *Cache) Acquire(r *http.Request, id int64) (*Conn, func(), error) {
	rawConn, found := c.connMap.Load(id)
	// If the connection isn't found, establish a new one!
	if !found {
		c.logger.Debugf("establishing new connection for %d", id)
		var err error
		// A singleflight group is used to allow for concurrent requests to the
		// same identifier to resolve.
		rawConn, err, _ = c.connGroup.Do(fmt.Sprintf("%d", id), func() (interface{}, error) {
			c.closeMutex.Lock()
			select {
			case <-c.closed:
				c.closeMutex.Unlock()
				return nil, xerrors.New("closed")
			default:
			}
			c.closeGroup.Add(1)
			c.closeMutex.Unlock()
			c.logger.Debugf("dialing workspace tunnel %d", id)
			agentConn, err := c.dialer(r, id)
			if err != nil {
				c.closeGroup.Done()
				return nil, xerrors.Errorf("dial: %w", err)
			}
			c.logger.Debugf("dialing workspace tunnel %d succeeded", id)
			timeoutCtx, timeoutCancelFunc := context.WithCancel(context.Background())
			defaultTransport, valid := http.DefaultTransport.(*http.Transport)
			if !valid {
				panic("dev error: default transport is the wrong type")
			}
			transport := defaultTransport.Clone()
			transport.DialContext = agentConn.DialContext
			conn := &Conn{
				AgentConn:     agentConn,
				timeoutCancel: timeoutCancelFunc,
				transport:     transport,
			}
			c.logger.Debugf("created connection to workspace %d", id)
			go func() {
				defer c.closeGroup.Done()
				var err error
				select {
				case <-timeoutCtx.Done():
					err = xerrors.New("cache timeout")
				case <-c.closed:
					err = xerrors.New("cache closed")
				case <-conn.Closed():
				}

				c.connMap.Delete(id)
				c.connGroup.Forget(fmt.Sprintf("%d", id))
				_ = conn.CloseWithError(err)
			}()
			return conn, nil
		})
		if err != nil {
			return nil, nil, err
		}
		c.connMap.Store(id, rawConn)
	} else {
		c.logger.Debugf("found existing connection for %d", id)
	}

	conn, _ := rawConn.(*Conn)
	conn.timeoutMutex.Lock()
	defer conn.timeoutMutex.Unlock()
	if conn.timeout != nil {
		conn.timeout.Stop()
	}
	conn.locks.Inc()
	return conn, func() {
		conn.timeoutMutex.Lock()
		defer conn.timeoutMutex.Unlock()
		if conn.timeout != nil {
			conn.timeout.Stop()
		}
		conn.locks.Dec()
		if conn.locks.Load() == 0 {
			conn.timeout = time.AfterFunc(c.inactiveTimeout, conn.timeoutCancel)
		}
	}, nil
}

// ForgetAndClose closes the connection and removes it from the cache.
func (c *Cache) ForgetAndClose(id int64) {
	rawConn, found := c.connMap.Load(id)
	if !found {
		return
	}
	conn, _ := rawConn.(*Conn)
	_ = conn.CloseWithError(nil)
	c.connMap.Delete(id)
	c.connGroup.Forget(fmt.Sprintf("%d", id))
}

func (c *Cache) Close() error {
	c.closeMutex.Lock()
	defer c.closeMutex.Unlock()
	select {
	case <-c.closed:
		return nil
	default:
	}
	_ = c.subscriber.Unsubscribe()
	close(c.closed)
	c.closeGroup.Wait()
	return nil
}

// handleForgetMsg watches the Jetstream for forget instructions to close connections
func (c *Cache) handleForgetMsg(msg *nats.Msg) {
	// ack message immediately
	_ = msg.Ack()

	// gob decode the message
	var forgetMsg models.ForgetConnMsg
	err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(&forgetMsg)
	if err != nil {
		c.logger.Errorf("error decoding forget message: %v", err)
		return
	}
	c.logger.Debugf("received forget message for %d - %d", forgetMsg.WorkspaceID, forgetMsg.AgentID)
	c.ForgetAndClose(forgetMsg.AgentID)
}

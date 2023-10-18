package websocket

import (
	"context"
	"go.uber.org/atomic"
	"golang.org/x/sync/singleflight"
	"net/http"
	"sync"
	"time"
)

// Socket
//
//	wraps a websocket connection with a reusable HTTP transport.
type Socket struct {
	conn *SocketConn

	locks         atomic.Uint64
	timeoutMutex  sync.Mutex
	timeout       *time.Timer
	timeoutCancel context.CancelFunc
	transport     *http.Transport
}

type Cache struct {
	closed          chan struct{}
	closeMutex      sync.Mutex
	closeGroup      sync.WaitGroup
	connGroup       singleflight.Group
	connMap         sync.Map
	dialer          Dialer
	inactiveTimeout time.Duration
}

// Dialer creates a new agent connection by ID.
type Dialer func(r *http.Request, id int64) (*SocketConn, error)

// New creates a new workspace connection cache that closes
// connections after the inactive timeout provided.
//
// Agent connections are cached due to WebRTC negotiation
// taking a few hundred milliseconds.
func New(dialer Dialer, inactiveTimeout time.Duration) *Cache {
	if inactiveTimeout == 0 {
		inactiveTimeout = 5 * time.Minute
	}
	return &Cache{
		closed:          make(chan struct{}),
		dialer:          dialer,
		inactiveTimeout: inactiveTimeout,
	}
}

func (s *Socket) HTTPTransport() *http.Transport {
	return s.transport
}

//// CloseWithError ends the HTTP transport if exists, and closes the agent.
//func (s *Socket) CloseWithError(err error) error {
//	if s.transport != nil {
//		s.transport.CloseIdleConnections()
//	}
//	s.timeoutMutex.Lock()
//	defer s.timeoutMutex.Unlock()
//	if s.timeout != nil {
//		s.timeout.Stop()
//	}
//	return s.conn.Close()
//}

//func (c *Cache) Acquire(r *http.Request, id int64) (*Socket, func(), error) {
//	rawConn, found := c.connMap.Load(id)
//	// If the connection isn't found, establish a new one!
//	if !found {
//		var err error
//		// A singleflight group is used to allow for concurrent requests to the
//		// same identifier to resolve.
//		rawConn, err, _ = c.connGroup.Do(fmt.Sprintf("%d", id), func() (interface{}, error) {
//			c.closeMutex.Lock()
//			select {
//			case <-c.closed:
//				c.closeMutex.Unlock()
//				return nil, xerrors.New("closed")
//			default:
//			}
//			c.closeGroup.Add(1)
//			c.closeMutex.Unlock()
//			socketConn, err := c.dialer(r, id)
//			if err != nil {
//				c.closeGroup.Done()
//				return nil, xerrors.Errorf("dial: %w", err)
//			}
//			timeoutCtx, timeoutCancelFunc := context.WithCancel(context.Background())
//			defaultTransport, valid := http.DefaultTransport.(*http.Transport)
//			if !valid {
//				panic("dev error: default transport is the wrong type")
//			}
//			transport := defaultTransport.Clone()
//			//transport.DialContext = agentConn.DialContext
//			conn := &Socket{
//				conn:     socketConn,
//				timeoutCancel: timeoutCancelFunc,
//				transport:     transport,
//			}
//			go func() {
//				defer c.closeGroup.Done()
//				var err error
//				select {
//				case <-timeoutCtx.Done():
//					err = xerrors.New("cache timeout")
//				case <-c.closed:
//					err = xerrors.New("cache closed")
//				case <-conn.Closed():
//				}
//
//				c.connMap.Delete(id)
//				c.connGroup.Forget(fmt.Sprintf("%d", id))
//				_ = conn.CloseWithError(err)
//			}()
//			return conn, nil
//		})
//		if err != nil {
//			return nil, nil, err
//		}
//		c.connMap.Store(id, rawConn)
//	}
//
//	conn, _ := rawConn.(*Conn)
//	conn.timeoutMutex.Lock()
//	defer conn.timeoutMutex.Unlock()
//	if conn.timeout != nil {
//		conn.timeout.Stop()
//	}
//	conn.locks.Inc()
//	return conn, func() {
//		conn.timeoutMutex.Lock()
//		defer conn.timeoutMutex.Unlock()
//		if conn.timeout != nil {
//			conn.timeout.Stop()
//		}
//		conn.locks.Dec()
//		if conn.locks.Load() == 0 {
//			conn.timeout = time.AfterFunc(c.inactiveTimeout, conn.timeoutCancel)
//		}
//	}, nil
//}

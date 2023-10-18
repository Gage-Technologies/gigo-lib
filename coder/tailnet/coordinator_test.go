package tailnet_test

import (
	"context"
	"encoding/json"
	"github.com/gage-technologies/gigo-lib/cluster"
	"github.com/gage-technologies/gigo-lib/config"
	"github.com/gage-technologies/gigo-lib/logging"
	"github.com/gage-technologies/gigo-lib/mq"
	"github.com/gage-technologies/gigo-lib/mq/streams"
	etcd "go.etcd.io/etcd/client/v3"
	"golang.org/x/xerrors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gage-technologies/gigo-lib/coder/tailnet"
)

type MockTailnetConn struct{}

func (m *MockTailnetConn) ConnectToCoordinatorNoOp(c net.Conn) chan error {
	return make(chan error)
}

func (m *MockTailnetConn) ConnectToCoordinatorHook(conn net.Conn, hook func([]*tailnet.Node) error, logger logging.Logger) chan error {
	errChan := make(chan error, 1)
	sendErr := func(err error) {
		select {
		case errChan <- err:
		default:
		}
	}
	go func() {
		decoder := json.NewDecoder(conn)
		for {
			logger.Debugf("(coordinator handler) waiting for nodes")
			var nodes []*tailnet.Node
			err := decoder.Decode(&nodes)
			if err != nil {
				logger.Errorf("(coordinator handler) exiting")
				sendErr(xerrors.Errorf("read: %w", err))
				return
			}
			logger.Debugf("(coordinator handler) received %d nodes", len(nodes))
			for _, n := range nodes {
				logger.Debugf("(coordinator handler) adding node as peer: %d - %+v", int64(n.ID), n)
			}

			err = hook(nodes)
			if err != nil {
				sendErr(xerrors.Errorf("update nodes: %w", err))
			}
		}
	}()

	// m.SetNodeCallback(func(node *Node) {
	// 	data, err := json.Marshal(node)
	// 	if err != nil {
	// 		sendErr(xerrors.Errorf("marshal node: %w", err))
	// 		return
	// 	}
	// 	c.logger.Debugf("(coordinator handler) sending node: %d - %+v", int64(node.ID), node)
	// 	_, err = conn.Write(data)
	// 	if err != nil {
	// 		sendErr(xerrors.Errorf("write: %w", err))
	// 	}
	// })

	return errChan
}

func sendNodesManual(conn net.Conn, node *tailnet.Node, errChan chan error) {
	data, err := json.Marshal(node)
	if err != nil {
		errChan <- xerrors.Errorf("marshal node: %w", err)
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		errChan <- xerrors.Errorf("write: %w", err)
	}
}

func TestCoordinator(t *testing.T) {
	leaderRoutine := func(ctx context.Context) error {
		return nil
	}
	followerRoutine := func(ctx context.Context) error {
		return nil
	}

	t.Run("ClientWithoutAgent", func(t *testing.T) {
		logger, err := logging.CreateBasicLogger(logging.NewDefaultBasicLoggerOptions("/tmp/gigo-test-coord-cwa.test"))
		require.NoError(t, err)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		clusterNode1 := cluster.NewStandaloneNode(ctx, 69, "test", leaderRoutine, followerRoutine, time.Second, logger)
		js, err := mq.NewJetstreamClient(config.JetstreamConfig{
			Host:        "mq://gigo-dev-nats:4222",
			MaxPubQueue: 256,
		}, logger)
		require.NoError(t, err)
		defer js.DeleteStream(streams.StreamTailscale)
		coordinator, err := tailnet.NewCoordinator(clusterNode1, js, logger)
		assert.NoError(t, err)
		client, server := net.Pipe()
		tailnetConn := &MockTailnetConn{}
		errChan := tailnetConn.ConnectToCoordinatorNoOp(client)
		closeChan := make(chan struct{})
		go func() {
			err := coordinator.ServeClient(server, 69, 420)
			assert.NoError(t, err)
			close(closeChan)
		}()
		sendNodesManual(client, &tailnet.Node{ID: 69}, errChan)
		require.Eventually(t, func() bool {
			s, err := coordinator.Server(69)
			assert.NoError(t, err)
			return s != nil
		}, 5*time.Second, 50*time.Millisecond)
		require.NoError(t, client.Close())
		require.NoError(t, server.Close())
		<-closeChan
	})

	t.Run("AgentWithoutClients", func(t *testing.T) {
		logger, err := logging.CreateBasicLogger(logging.NewDefaultBasicLoggerOptions("/tmp/gigo-test-coord-awc.test"))
		require.NoError(t, err)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		clusterNode2 := cluster.NewStandaloneNode(ctx, 69, "test", leaderRoutine, followerRoutine, time.Second, logger)
		js, err := mq.NewJetstreamClient(config.JetstreamConfig{
			Host:        "mq://gigo-dev-nats:4222",
			MaxPubQueue: 256,
		}, logger)
		require.NoError(t, err)
		defer js.DeleteStream(streams.StreamTailscale)
		coordinator, err := tailnet.NewCoordinator(clusterNode2, js, logger)
		assert.NoError(t, err)
		client, server := net.Pipe()
		tailnetConn := &MockTailnetConn{}
		errChan := tailnetConn.ConnectToCoordinatorNoOp(client)
		closeChan := make(chan struct{})
		go func() {
			err := coordinator.ServeAgent(server, 420)
			assert.NoError(t, err)
			close(closeChan)
		}()
		sendNodesManual(client, &tailnet.Node{ID: 420}, errChan)
		require.Eventually(t, func() bool {
			s, err := coordinator.Agent(420)
			assert.NoError(t, err)
			return s != nil
		}, 5*time.Second, 50*time.Millisecond)
		err = client.Close()
		require.NoError(t, err)
		<-closeChan
	})

	t.Run("AgentWithClient", func(t *testing.T) {
		logger, err := logging.CreateBasicLogger(logging.NewDefaultBasicLoggerOptions("/tmp/gigo-test-coord-full.test"))
		require.NoError(t, err)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		clusterNode3 := cluster.NewStandaloneNode(ctx, 69, "test", leaderRoutine, followerRoutine, time.Second, logger)
		js, err := mq.NewJetstreamClient(config.JetstreamConfig{
			Host:        "mq://gigo-dev-nats:4222",
			MaxPubQueue: 256,
		}, logger)
		require.NoError(t, err)
		defer js.DeleteStream(streams.StreamTailscale)
		coordinator, err := tailnet.NewCoordinator(clusterNode3, js, logger)
		agentWS, agentServerWS := net.Pipe()
		defer agentWS.Close()
		agentNodeChan := make(chan []*tailnet.Node)
		agentTailnetConn := &MockTailnetConn{}
		agentErrChan := agentTailnetConn.ConnectToCoordinatorHook(agentWS, func(nodes []*tailnet.Node) error {
			agentNodeChan <- nodes
			return nil
		}, logger)
		logger.Debug("(test) served coord")
		closeAgentChan := make(chan struct{})
		go func() {
			err := coordinator.ServeAgent(agentServerWS, 420)
			assert.NoError(t, err)
			logger.Debug("(test) agent serve loop exited")
			close(closeAgentChan)
		}()
		logger.Debug("(test) served agent")
		sendNodesManual(agentWS, &tailnet.Node{ID: 420}, agentErrChan)
		require.Eventually(t, func() bool {
			s, err := coordinator.Agent(420)
			assert.NoError(t, err)
			return s != nil
		}, 5*time.Second, 50*time.Millisecond)
		logger.Debug("(test) agent received")

		clientWS, clientServerWS := net.Pipe()
		defer clientWS.Close()
		defer clientServerWS.Close()
		clientNodeChan := make(chan []*tailnet.Node, 5)
		clientTailnetConn := &MockTailnetConn{}
		clientErrChan := clientTailnetConn.ConnectToCoordinatorHook(clientWS, func(nodes []*tailnet.Node) error {
			clientNodeChan <- nodes
			return nil
		}, logger)
		logger.Debug("(test) served coord 2")
		closeClientChan := make(chan struct{})
		go func() {
			err := coordinator.ServeClient(clientServerWS, 69, 420)
			assert.NoError(t, err)
			close(closeClientChan)
		}()
		logger.Debug("(test) served client")
		agentNodes := <-clientNodeChan
		logger.Debug("(test) client received agent nodes")
		require.Len(t, agentNodes, 1)
		sendNodesManual(clientWS, &tailnet.Node{ID: 69}, clientErrChan)
		clientNodes := <-agentNodeChan
		logger.Debug("(test) agent received client nodes")
		require.Len(t, clientNodes, 1)

		// Ensure an update to the agent node reaches the client!
		sendNodesManual(agentWS, &tailnet.Node{ID: 420}, agentErrChan)
		agentNodes = <-clientNodeChan
		logger.Debug("(test) client received agent nodes 2")
		require.Len(t, agentNodes, 1)

		// Close the agent WebSocket so a new one can connect.
		err = agentWS.Close()
		require.NoError(t, err)
		logger.Debug("(test) agent closed")
		// <-agentErrChan
		logger.Debug("(test) agent error received")
		<-closeAgentChan
		logger.Debug("(test) agent channels closed")

		// Create a new agent connection. This is to simulate a reconnect!
		agentWS, agentServerWS = net.Pipe()
		defer agentWS.Close()
		agentNodeChan = make(chan []*tailnet.Node)
		agentTailnetConn = &MockTailnetConn{}
		agentErrChan = agentTailnetConn.ConnectToCoordinatorHook(agentWS, func(nodes []*tailnet.Node) error {
			agentNodeChan <- nodes
			return nil
		}, logger)
		logger.Debug("(test) serve coord 3")
		closeAgentChan = make(chan struct{})
		go func() {
			err := coordinator.ServeAgent(agentServerWS, 420)
			assert.NoError(t, err)
			close(closeAgentChan)
		}()
		logger.Debug("(test) serve agent 2")

		// Ensure the existing listening client sends it's node immediately!
		clientNodes = <-agentNodeChan
		logger.Debug("(test) agent received client nodes 2")
		require.Len(t, clientNodes, 1)

		err = agentWS.Close()
		logger.Debug("(test) agent closed 2")
		require.NoError(t, err)
		<-agentErrChan
		<-closeAgentChan
		logger.Debug("(test) agent channels closed 2")

		err = clientWS.Close()
		logger.Debug("(test) client closed")
		require.NoError(t, err)
		// we have to clear the buffer since the agent watcher can cause
		// duplicate writes - this isn't a problem in a deployment because
		// the duplications are just consumed by the servers
		logger.Debug("(test) client waiting for agent chan to close")
		for len(clientNodeChan) > 0 {
			<-clientNodeChan
		}
		logger.Debug("(test) client waiting for agent chan to close 2")
		<-clientErrChan
		logger.Debug("(test) client waiting for agent chan to close 3")
		<-closeClientChan
		logger.Debug("(test) client channels closed")
	})

	t.Run("AgentWithClientClusterNode", func(t *testing.T) {
		logger, err := logging.CreateBasicLogger(logging.NewDefaultBasicLoggerOptions("/tmp/gigo-test-coord-full-clsuter.test"))
		require.NoError(t, err)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		clusterNode1, err := cluster.NewClusterNode(cluster.ClusterNodeOptions{
			ctx,
			69,
			"test",
			time.Second,
			"coord-test",
			etcd.Config{
				Endpoints: []string{"gigo-dev-etcd:2379"},
			},
			leaderRoutine,
			followerRoutine,
			time.Second,
			logger,
		})
		clusterNode1.Start()
		defer func() {
			clusterNode1.Stop()
			clusterNode1.Close()
		}()
		clusterNode2, err := cluster.NewClusterNode(cluster.ClusterNodeOptions{
			ctx,
			692,
			"test",
			time.Second,
			"coord-test",
			etcd.Config{
				Endpoints: []string{"gigo-dev-etcd:2379"},
			},
			leaderRoutine,
			followerRoutine,
			time.Second,
			logger,
		})
		clusterNode2.Start()
		defer func() {
			clusterNode2.Stop()
			clusterNode2.Close()
		}()
		time.Sleep(time.Second * 2)
		js, err := mq.NewJetstreamClient(config.JetstreamConfig{
			Host:        "mq://gigo-dev-nats:4222",
			MaxPubQueue: 256,
		}, logger)
		require.NoError(t, err)
		defer js.DeleteStream(streams.StreamTailscale)
		coordinator1, err := tailnet.NewCoordinator(clusterNode1, js, logger)
		require.NoError(t, err)
		coordinator2, err := tailnet.NewCoordinator(clusterNode2, js, logger)
		require.NoError(t, err)
		logger.Debug("(test) coord created 1")
		agentWS, agentServerWS := net.Pipe()
		defer agentWS.Close()
		agentNodeChan := make(chan []*tailnet.Node)
		agentTailnetConn := &MockTailnetConn{}
		agentErrChan := agentTailnetConn.ConnectToCoordinatorHook(agentWS, func(nodes []*tailnet.Node) error {
			agentNodeChan <- nodes
			return nil
		}, logger)
		logger.Debug("(test) coord served 1")
		closeAgentChan := make(chan struct{})
		go func() {
			err := coordinator1.ServeAgent(agentServerWS, 420)
			assert.NoError(t, err)
			close(closeAgentChan)
		}()
		logger.Debug("(test) agent served 1")
		sendNodesManual(agentWS, &tailnet.Node{ID: 420}, agentErrChan)
		logger.Debug("(test) agent sent 1")
		require.Eventually(t, func() bool {
			s, err := coordinator1.Agent(420)
			assert.NoError(t, err)
			return s != nil
		}, 5*time.Second, 50*time.Millisecond)
		logger.Debug("(test) agent registered 1")

		clientWS, clientServerWS := net.Pipe()
		defer clientWS.Close()
		defer clientServerWS.Close()
		clientNodeChan := make(chan []*tailnet.Node, 5)
		clientTailnetConn := &MockTailnetConn{}
		clientErrChan := clientTailnetConn.ConnectToCoordinatorHook(clientWS, func(nodes []*tailnet.Node) error {
			clientNodeChan <- nodes
			return nil
		}, logger)
		logger.Debug("(test) coord served 2")
		closeClientChan := make(chan struct{})
		go func() {
			err := coordinator2.ServeClient(clientServerWS, 692, 420)
			assert.NoError(t, err)
			close(closeClientChan)
		}()
		logger.Debug("(test) client served 1")
		agentNodes := <-clientNodeChan
		logger.Debug("(test) client received agents")
		require.Len(t, agentNodes, 1)
		sendNodesManual(clientWS, &tailnet.Node{ID: 692}, clientErrChan)
		// writing client node does not make it back to the agent
		clientNodes := <-agentNodeChan
		logger.Debug("(test) agent received clients")
		require.Len(t, clientNodes, 1)

		// Ensure an update to the agent node reaches the client!
		sendNodesManual(agentWS, &tailnet.Node{ID: 420}, agentErrChan)
		agentNodes = <-clientNodeChan
		logger.Debug("(test) client received agents")
		require.Len(t, agentNodes, 1)

		// Close the agent WebSocket so a new one can connect.
		err = agentWS.Close()
		require.NoError(t, err)
		logger.Debug("(test) closing agent 1")
		<-agentErrChan
		<-closeAgentChan

		// Create a new agent connection. This is to simulate a reconnect!
		agentWS, agentServerWS = net.Pipe()
		defer agentWS.Close()
		agentNodeChan = make(chan []*tailnet.Node)
		agentTailnetConn = &MockTailnetConn{}
		agentErrChan = agentTailnetConn.ConnectToCoordinatorHook(agentWS, func(nodes []*tailnet.Node) error {
			agentNodeChan <- nodes
			return nil
		}, logger)
		logger.Debug("(test) coord served 3")
		closeAgentChan = make(chan struct{})
		go func() {
			err := coordinator1.ServeAgent(agentServerWS, 420)
			assert.NoError(t, err)
			close(closeAgentChan)
		}()
		logger.Debug("(test) agent served 2")

		// Ensure the existing listening client sends it's node immediately!
		clientNodes = <-agentNodeChan
		logger.Debug("(test) agent received clients")
		require.Len(t, clientNodes, 1)

		err = agentWS.Close()
		require.NoError(t, err)
		logger.Debug("(test) closing agent 2")
		<-agentErrChan
		<-closeAgentChan

		logger.Debug("(test) closing client 1")

		err = clientWS.Close()
		require.NoError(t, err)

		logger.Debug("(test) closing client node chan")
		// we have to clear the buffer since the agent watcher can cause
		// duplicate writes - this isn't a problem in a deployment because
		// the duplications are just consumed by the servers
		for len(clientNodeChan) > 0 {
			<-clientNodeChan
		}

		logger.Debug("(test) closing client err chan")
		<-clientErrChan

		logger.Debug("(test) closing client close chan")
		<-closeClientChan

		logger.Debug("(test) done")
	})
}

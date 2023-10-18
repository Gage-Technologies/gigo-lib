package tailnet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gage-technologies/gigo-lib/mq"
	"github.com/gage-technologies/gigo-lib/mq/streams"
	"github.com/nats-io/nats.go"
	"io"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/coder/retry"
	"github.com/gage-technologies/gigo-lib/cluster"
	"github.com/gage-technologies/gigo-lib/logging"
	"github.com/sourcegraph/conc"

	"golang.org/x/xerrors"
	"tailscale.com/tailcfg"
	"tailscale.com/types/key"
)

type AgentConnMsgType int

const (
	AgentConnMsgTypeAdd AgentConnMsgType = iota
	AgentConnMsgTypeRemove
)

const (
	// AgentPrefix is the prefix key that is used to save an agent's tailnet node information to the cluster state
	AgentPrefix = "ts/coord/agent"
	// ServerPrefix is the prefix key that is used to save a server's tailnet node information to the cluster state
	ServerPrefix = "ts/coord/server"
	// ConnectionPrefix is the prefix key that is used to save a connection's tailnet information to the cluster state
	ConnectionPrefix = "ts/coord/connection"
)

// Node represents a node in the network.
type Node struct {
	// ID is used to identify the connection.
	ID tailcfg.NodeID `json:"id"`
	// AsOf is the time the node was created.
	AsOf time.Time `json:"as_of"`
	// Key is the Wireguard public key of the node.
	Key key.NodePublic `json:"key"`
	// DiscoKey is used for discovery messages over DERP to establish peer-to-peer connections.
	DiscoKey key.DiscoPublic `json:"disco"`
	// PreferredDERP is the DERP server that peered connections
	// should meet at to establish.
	PreferredDERP int `json:"preferred_derp"`
	// DERPLatency is the latency in seconds to each DERP server.
	DERPLatency map[string]float64 `json:"derp_latency"`
	// Addresses are the IP address ranges this connection exposes.
	Addresses []netip.Prefix `json:"addresses"`
	// AllowedIPs specify what addresses can dial the connection.
	// We allow all by default.
	AllowedIPs []netip.Prefix `json:"allowed_ips"`
	// Endpoints are ip:port combinations that can be used to establish
	// peer-to-peer connections.
	Endpoints []string `json:"endpoints"`
}

// AgentConnMsg is the metadata of an agent
type AgentConnMsg struct {
	Type    AgentConnMsgType `json:"type"`
	AgentID int64            `json:"agent_id"`
	Node    Node             `json:"node"`
}

// ConnectionMetadata is the metadata associated with a connection
type ConnectionMetadata struct {
	// ID is the unique ID of the connection
	ID int64 `json:"_id"`
	// AgentID is the unique ID of the agent
	AgentID int64 `json:"agent_id"`
	// ServerID is the unique ID of the server
	ServerID int64 `json:"server_id"`
	// CreatedAt is the time the connection was created
	CreatedAt time.Time `json:"created_at"`
}

// Coordinator exchanges nodes with agents to establish connections.
// ┌──────────────────┐   ┌────────────────────┐   ┌───────────────────┐   ┌──────────────────┐
// │tailnet.Coordinate├──►│tailnet.AcceptClient│◄─►│tailnet.AcceptAgent│◄──┤tailnet.Coordinate│
// └──────────────────┘   └────────────────────┘   └───────────────────┘   └──────────────────┘
type Coordinator struct {
	ctx    context.Context
	cancel context.CancelFunc
	mutex  sync.Mutex
	closed bool
	wg     *conc.WaitGroup

	// clusterNode node within the Gigo cluster that is the owner of this coordinator
	clusterNode cluster.Node
	// jsClient jetstream client used to publish changes to the agent and connection states to other nodes
	jsClient *mq.JetstreamClient
	// agentNodeSubscription jetstream subscription used to propagate updates to the servers hosting client connections to the agents on any changes to the agent
	agentNodeSubscription *nats.Subscription
	// connectionSubscription jetstream subscription used to propagate updates to the servers hosting agent connections on what other nodes want the connect to the agent
	connectionSubscription *nats.Subscription
	// agentSockets maps agent IDs to their open websocket.
	agentSockets map[int64]net.Conn
	// serverSockets maps agent IDs to the server socket that correspond to the server<->agent connection
	// so that we can update the server connections with the latest agent tailscale nodes
	serverSockets map[int64]map[int64]net.Conn
	logger        logging.Logger
}

// NewCoordinator constructs a new in-memory connection Coordinator. This
// Coordinator is incompatible with multiple Coder replicas as all node data is
// in-memory.
func NewCoordinator(clusterNode cluster.Node, jsClient *mq.JetstreamClient, logger logging.Logger) (*Coordinator, error) {
	// create context for the coordinator
	ctx, cancel := context.WithCancel(context.Background())

	// create a new coordinator
	coord := &Coordinator{
		ctx:           ctx,
		cancel:        cancel,
		wg:            conc.NewWaitGroup(),
		closed:        false,
		clusterNode:   clusterNode,
		jsClient:      jsClient,
		agentSockets:  map[int64]net.Conn{},
		serverSockets: map[int64]map[int64]net.Conn{},
		logger:        logger,
	}

	// create an async consumer for the agent node changes
	agentSub, err := jsClient.Subscribe(streams.SubjectTailscaleAgent, coord.watchAgentNodeChanges)
	if err != nil {
		return nil, xerrors.Errorf("subscribe agent node changes: %w", err)
	}
	coord.agentNodeSubscription = agentSub

	// create an async consumer for the agent connection changes
	connSub, err := jsClient.Subscribe(streams.SubjectTailscaleConnection, coord.watchAgentConnectionChanges)
	if err != nil {
		return nil, xerrors.Errorf("subscribe agent connection changes: %w", err)
	}
	coord.connectionSubscription = connSub

	return coord, nil
}

// Agent
//
//	Returns the tailscale node for the passed agent ID.
//	If the agent does not exist, nil is returned.
func (c *Coordinator) Agent(id int64) (*Node, error) {
	// check all nodes in the cluster for the agent
	clusterAgents, err := c.clusterNode.GetCluster(fmt.Sprintf("%s/%d", AgentPrefix, id))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve node as agent %d: %v", id, err)
	}

	// create nil node to hold tailnet node for agent
	var agent *Node

	// iterate cluster agent data selecting the most recent agent node
	for _, data := range clusterAgents {
		for _, agentData := range data {
			// unmarshal node data
			node := &Node{}
			err = json.Unmarshal([]byte(agentData.Value), node)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal agent node %d: %v", id, err)
			}

			// if this is the first agent accept it and continue
			if agent == nil {
				agent = node
			}

			// if there is an existing agent we only want to keep this one
			// if it was created more recently than the agent we already have
			if node.AsOf.After(agent.AsOf) {
				agent = node
			}
		}
	}

	return agent, nil
}

// Server
//
//	Returns the tailscale node for the passed server ID.
//	If the server does not exist, nil is returned.
func (c *Coordinator) Server(id int64) (*Node, error) {
	// request the server's node from itself since the only
	// server that should be able to publish its node data
	// is itself
	nodeData, err := c.clusterNode.GetAsNode(id, ServerPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve node as server %d: %v", id, err)
	}

	// return nil if the node is not found
	if nodeData == "" {
		return nil, nil
	}

	// unmarshal node data
	node := &Node{}
	err = json.Unmarshal([]byte(nodeData), node)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal server node %d: %v", id, err)
	}

	return node, nil
}

// watchAgentNode
//
//	Watches for changes to the agent's tailscale node and updates any
//	servers subscribed to the agent's changes
func (c *Coordinator) watchAgentNodeChanges(msg *nats.Msg) {
	// format logger with the watcher flag
	logger := c.logger.WithName("gigo-core-coordinator-agent-node-watcher")

	// ignore empty events - sanity check
	if len(msg.Data) == 0 {
		return
	}

	// gob decode the message
	var meta AgentConnMsg
	err := json.Unmarshal(msg.Data, &meta)
	if err != nil {
		logger.Errorf("failed to decode agent metadata: %v", err)
		return
	}

	logger.Debugf("agent tailnet node update: %d", meta.AgentID)

	// update the agent node in all of the known server sockets asynchronously
	c.mutex.Lock()
	defer c.mutex.Unlock()
	sockets := c.serverSockets[meta.AgentID]
	if len(sockets) > 0 {
		for connectionId, serverSocket := range sockets {
			connId := connectionId
			socket := serverSocket
			c.wg.Go(func() {
				// update the server connection to make its coordinator aware of the agent
				// so we can DERP a wireguard connection
				data, err := json.Marshal([]*Node{&meta.Node})
				if err != nil {
					logger.Errorf("failed to marshal agent node during agent update for connection %d: %v", connId, err)
					return
				}
				_, err = socket.Write(data)
				if err != nil {
					logger.Errorf("failed to write agent node during agent update for connection %d: %v", connId, err)
					return
				}
			})
		}
	}

	// ack message
	err = msg.Ack()
	if err != nil {
		logger.Errorf("failed to ack agent node update: %v", err)
		return
	}
}

// watchAgentConnectionChanges
//
//	  Watches for changes to connections related to agents and updates the
//		 agent if we are the owner of the agent socket.
func (c *Coordinator) watchAgentConnectionChanges(msg *nats.Msg) {
	// format logger with the watcher flag
	logger := c.logger.WithName("gigo-core-coordinator-agent-connection-watcher")

	// ignore empty events - sanity check
	if len(msg.Data) == 0 {
		return
	}

	// gob decode the message
	var meta ConnectionMetadata
	err := json.Unmarshal(msg.Data, &meta)
	if err != nil {
		logger.Errorf("failed to decode connection metadata: %v", err)
		return
	}

	logger.Debugf("connection update event: %d - %d - %d", meta.ID, meta.AgentID, meta.ServerID)

	// attempt to retrieve agent socket and skip if we are not the
	// owner of the agent socket
	c.mutex.Lock()
	defer c.mutex.Unlock()
	socket, ok := c.agentSockets[meta.AgentID]
	if !ok {
		return
	}

	// retrieve all of the active agent connections
	agentConnections, err := c.clusterNode.GetCluster(fmt.Sprintf("%s/%d", ConnectionPrefix, meta.AgentID))
	if err != nil {
		logger.Errorf("failed to retrieve active agent connections for agent %d: %d", meta.AgentID, meta.ID)
		return
	}

	// create slice to hold nodes attempting to connect to the agent
	serverNodes := make([]*Node, 0)

	// iterate agent connection for all nodes and prep the agent with
	// the servers that are expecting to be connected to the agent
	for serverId, connections := range agentConnections {
		// skip if there are no connections to this agent from this server
		if len(connections) == 0 {
			continue
		}

		// TODO: we may need to retry here since it is likely that the server will not have
		// put its node on the first connection

		// retry up to 10 times since the server node could be written to the cluster
		// after the connection on the first connection of that server
		for localRetrier := retry.New(time.Millisecond*50, time.Second); localRetrier.Wait(c.ctx); {
			logger.Debugf("attempting to retrieve server node %d for connection to agent %d", serverId, meta.AgentID)

			// we only care to add the first connection for this server since
			// we'll get all the node data we need from the first - unlike agents
			// servers don't change IPs
			serverNode, err := c.Server(serverId)
			if err != nil {
				logger.Errorf("failed to get server node %d while updating agent connection: %v", serverId, err)
				continue
			}

			// skip if we can't find the server
			if serverNode == nil {
				logger.Warnf("failed to locate server node %d while updating agent connection", serverId)
				continue
			}

			// append to the server nodes slice to be sent to the agent
			serverNodes = append(serverNodes, serverNode)

			logger.Debugf("added server node %d to agent %d peers update", serverId, meta.AgentID)

			// exit retry loop
			break
		}
	}

	// marshall and write the server tailnet node to the agent socket
	data, err := json.Marshal(serverNodes)
	if err != nil {
		logger.Errorf("failed to marshall server nodes while updating agent connection: %v", err)
		return
	}
	_, err = socket.Write(data)
	if err != nil {

		logger.Errorf("failed write nodes to agent connection while updating agent connection: %v", err)
		return
	}
}

// ServeClient accepts a WebSocket connection that wants to connect to an agent
// with the specified ID.
func (c *Coordinator) ServeClient(conn net.Conn, connectionID int64, agent int64) error {
	c.mutex.Lock()
	closed := c.closed
	c.mutex.Unlock()

	if closed {
		return xerrors.New("coordinator is closed")
	}

	// When a new connection is requested, we update it with the latest
	// node of the agent. This allows the connection to establish.

	// retrieve the agent from the cluster state
	node, err := c.Agent(agent)
	if err != nil {
		return fmt.Errorf("failed to retrieve agent node %d: %v", agent, err)
	}

	// handle an existing agent connection
	if node != nil {
		c.logger.Debugf("init client connection with existing agent: %d", agent)
		// update the server connection to make its coordinator aware of the agent
		// so we can DERP a wireguard connection
		data, err := json.Marshal([]*Node{node})
		if err != nil {
			return xerrors.Errorf("marshal node: %w", err)
		}
		_, err = conn.Write(data)
		if err != nil {
			return xerrors.Errorf("write nodes: %w", err)
		}
	}

	// retrieve our local connections for the agent
	c.mutex.Lock()
	connectionSockets, ok := c.serverSockets[agent]
	if !ok {
		// initialize a set of connections for the agent if we don't
		// have any connections to the agent from this node
		connectionSockets = map[int64]net.Conn{}
		c.serverSockets[agent] = connectionSockets
	}
	// save this connection to the local server connections under the agent
	// and connection id so that we can update the server with changes to
	// the agent's tailnet node
	connectionSockets[connectionID] = conn
	c.mutex.Unlock()

	// create a new connection metadata and add it to the cluster state
	// so that we can track what servers are trying to connect with
	// what agents across the cluster
	connectionMeta := ConnectionMetadata{
		ID:        connectionID,
		AgentID:   agent,
		ServerID:  c.clusterNode.GetSelfMetadata().ID,
		CreatedAt: time.Now(),
	}
	buf, err := json.Marshal(connectionMeta)
	if err != nil {
		return fmt.Errorf("failed to marshal connection metadata: %v", err)
	}
	err = c.clusterNode.Put(fmt.Sprintf("%s/%d/%d", ConnectionPrefix, agent, connectionID), string(buf))
	if err != nil {
		return fmt.Errorf("failed to save connection metadata: %v", err)
	}

	// propagate the connection metadata to the cluster
	_, err = c.jsClient.Publish(streams.SubjectTailscaleConnection, buf)
	if err != nil {
		return fmt.Errorf("failed to publish connection metadata: %v", err)
	}

	// defer a cleanup function for this connection
	defer func() {
		// optimistically remove the metadata for this cluster connection
		err := c.clusterNode.Delete(fmt.Sprintf("%s/%d/%d", ConnectionPrefix, agent, connectionID))
		if err != nil {
			c.logger.Error(fmt.Sprintf(
				"failed to remove metadata for connection %d: %d -> %d",
				connectionID, c.clusterNode.GetSelfMetadata().ID, agent,
			))
		}

		c.mutex.Lock()
		defer c.mutex.Unlock()

		// remove this connection from the local server sockets
		connectionSockets, ok := c.serverSockets[agent]
		if !ok {
			return
		}
		delete(connectionSockets, connectionID)
	}()

	decoder := json.NewDecoder(conn)
	for {
		err := c.handleNextClientMessage(agent, decoder, connectionMeta)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return xerrors.Errorf("handle next client message: %w", err)
		}
	}
}

// writes server to itself

// handleNextClientMessage
//
//	Updates the agent with the Gigo core server's node information.
//	This is where the agent is told what nodes (Gigo core servers) it
//	should establish wiregaurd connections with.
func (c *Coordinator) handleNextClientMessage(agent int64, decoder *json.Decoder, connectionMeta ConnectionMetadata) error {
	// this will be a Gigo core server node
	var node Node
	err := decoder.Decode(&node)
	if err != nil {
		return xerrors.Errorf("read json: %w", err)
	}

	c.logger.Debugf("server node %d to be forwarded to agent %d", int64(node.ID), agent)

	// update the server in the cluster data
	c.logger.Debugf("adding new server node %d to cluster", int64(node.ID))
	buf, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshal node: %v", err)
	}

	// we put the server directly to the cluster node without its id
	// because this will only ever put the local server node it will
	// never have a remote server node here
	err = c.clusterNode.Put(ServerPrefix, string(buf))
	if err != nil {
		return fmt.Errorf("failed to put server node: %v", err)
	}

	// update connection metadata to trigger the connection watcher routine
	// to modify the connection
	buf, err = json.Marshal(connectionMeta)
	if err != nil {
		return fmt.Errorf("failed to marshal connection metadata: %v", err)
	}
	err = c.clusterNode.Put(fmt.Sprintf("%s/%d/%d", ConnectionPrefix, agent, connectionMeta.ID), string(buf))
	if err != nil {
		return fmt.Errorf("failed to save connection metadata: %v", err)
	}

	// propagate the connection metadata update to the cluster
	_, err = c.jsClient.Publish(streams.SubjectTailscaleConnection, buf)
	if err != nil {
		return fmt.Errorf("failed to publish connection metadata: %v", err)
	}

	return nil
}

// ServeAgent accepts a WebSocket connection to an agent that
// listens to incoming connections and publishes node updates.
func (c *Coordinator) ServeAgent(conn net.Conn, agent int64) error {
	c.logger.Debugf("(coordinator: %d) serving agent", agent)
	c.mutex.Lock()
	closed := c.closed
	c.mutex.Unlock()

	if closed {
		return xerrors.New("coordinator is closed")
	}

	c.logger.Debugf("(coordinator: %d) retrieving agent connections", agent)

	// retrieve all servers that are currently attempting to connect to an agent
	agentConnections, err := c.clusterNode.GetCluster(fmt.Sprintf("%s/%d", ConnectionPrefix, agent))
	if err != nil {
		return xerrors.Errorf("failed to get cluster node: %w", err)
	}

	c.logger.Debugf("(coordinator: %d) retrieved agent connections: %d", agent, len(agentConnections))

	// create slice to hold the server nodes that will initialize the agent
	// connection
	serverNodes := make([]*Node, 0)

	// iterate agent connection for all nodes and prep the agent with
	// the servers that are expecting to be connected to the agent
	for serverId, connections := range agentConnections {
		// skip if there are no connections to this agent from this server
		if len(connections) == 0 {
			continue
		}

		// we only care to add the first connection for this server since
		// we'll get all the node data we need from the first - unlike agents
		// servers don't change IPs
		serverNode, err := c.Server(serverId)
		if err != nil {
			c.logger.Errorf("failed to get server node %d while initializing agent connection: %v", serverId, err)
			continue
		}

		// skip if we can't find the server
		if serverNode == nil {
			c.logger.Warnf("failed to locate server node %d while initializing agent connection", serverId)
			continue
		}

		// append to the server nodes slice to be sent to the agent
		serverNodes = append(serverNodes, serverNode)
	}

	// marshall and write the server tailnet node to the agent socket
	if len(serverNodes) > 0 {
		c.logger.Debugf("(coordinator: %d) writing server nodes to agent socket: %d", agent, len(serverNodes))
		data, err := json.Marshal(serverNodes)
		if err != nil {
			return xerrors.Errorf("marshal json: %w", err)
		}
		_, err = conn.Write(data)
		if err != nil {
			return xerrors.Errorf("write nodes: %w", err)
		}
		c.logger.Debugf("(coordinator: %d) wrote server nodes to agent socket", agent)
	}

	c.mutex.Lock()
	// If an old agent socket is connected, we close it
	// to avoid any leaks. This shouldn't ever occur because
	// we expect one agent to be running.
	oldAgentSocket, ok := c.agentSockets[agent]
	if ok {
		_ = oldAgentSocket.Close()
	}
	c.agentSockets[agent] = conn
	c.mutex.Unlock()

	c.logger.Debugf("(coordinator: %d) agent socket established", agent)

	defer func() {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		// we perform the agent node cleanup for the cluster optimistically
		// assuming the connection was established and the node was written
		// to the cluster
		_ = c.clusterNode.Delete(fmt.Sprintf("%s/%d", AgentPrefix, agent))

		delete(c.agentSockets, agent)
	}()

	decoder := json.NewDecoder(conn)
	for {
		c.logger.Debugf("(coordinator: %d) waiting for agent message", agent)
		err := c.handleNextAgentMessage(agent, decoder)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
				return nil
			}
			return xerrors.Errorf("handle next agent message: %w", err)
		}
	}
}

// writes the agent to itself - in turn publishes the agent to all servers that are talking to the agent

// handleNextAgentMessage
//
//	Updates the coordinators node state with messages received from the
//	agent. This is where the agent announces itself to the Coordinator so
//	that the wireguard connection can be established between the Gigo core
//	server and the workspace agent.
func (c *Coordinator) handleNextAgentMessage(agent int64, decoder *json.Decoder) error {
	// this will be the tailnet node of the agent
	var node Node
	err := decoder.Decode(&node)
	if err != nil {
		return xerrors.Errorf("read json: %w", err)
	}

	c.logger.Debugf("(coordinator: %d) received agent message: %d - %+v\n", agent, int64(node.ID), node)

	// add the new agent to the cluster state so the servers watching for
	// the agent node become aware of it
	buf, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshall agent node while attempting to write to etcd: %v", err)
	}
	err = c.clusterNode.Put(fmt.Sprintf("%s/%d", AgentPrefix, agent), string(buf))
	if err != nil {
		return fmt.Errorf("failed to write agent node while attempting to write to etcd: %v", err)
	}

	// propagate the agent to the cluster
	streamBuf, err := json.Marshal(AgentConnMsg{
		Type:    AgentConnMsgTypeAdd,
		AgentID: agent,
		Node:    node,
	})
	if err != nil {
		return fmt.Errorf("failed to marshall agent metadata: %v", err)
	}
	_, err = c.jsClient.Publish(streams.SubjectTailscaleAgent, streamBuf)
	if err != nil {
		return fmt.Errorf("failed to publish agent metadata: %v", err)
	}

	return nil
}

// Close closes all of the open connections in the Coordinator and stops the
// Coordinator from accepting new connections.
func (c *Coordinator) Close() error {
	c.mutex.Lock()
	if c.closed {
		c.mutex.Unlock()
		return nil
	}
	c.cancel()
	c.closed = true

	c.wg.Go(func() {
		_ = c.agentNodeSubscription.Unsubscribe()
	})
	c.wg.Go(func() {
		_ = c.connectionSubscription.Unsubscribe()
	})

	for _, socket := range c.agentSockets {
		socket := socket
		c.wg.Go(func() {
			_ = socket.Close()
		})
	}

	for _, connections := range c.serverSockets {
		for _, socket := range connections {
			socket := socket
			c.wg.Go(func() {
				_ = socket.Close()
			})
		}
	}

	c.mutex.Unlock()

	c.wg.Wait()
	return nil
}

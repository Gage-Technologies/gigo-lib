package agentsdk

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"net"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"tailscale.com/ipn/ipnstate"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"golang.org/x/xerrors"
	"tailscale.com/net/speedtest"

	"github.com/gage-technologies/gigo-lib/coder/tailnet"
)

var (
	// TailnetIP is a static IPv6 address with the Tailscale prefix that is used to route
	// connections from clients to this node. A dynamic address is not required because a Tailnet
	// client only dials a single agent at a time.
	TailnetIP = netip.MustParseAddr("fd7a:115c:a1e0:49d6:b259:b7ac:b1b2:48f4")
)

const (
	TailnetSSHPort             = 1
	TailnetReconnectingPTYPort = 2
	TailnetSpeedtestPort       = 3
	// TailnetStatisticsPort serves a HTTP server with endpoints for gathering
	// agent statistics.
	TailnetStatisticsPort = 4

	// MinimumListeningPort is the minimum port that the listening-ports
	// endpoint will return to the client, and the minimum port that is accepted
	// by the proxy applications endpoint. Coder consumes ports 1-4 at the
	// moment, and we reserve some extra ports for future use. Port 9 and up are
	// available for the user.
	//
	// This is not enforced in the CLI intentionally as we don't really care
	// *that* much. The user could bypass this in the CLI by using SSH instead
	// anyways.
	MinimumListeningPort = 9
)

// IgnoredListeningPorts contains a list of ports in the global ignore list.
// This list contains common TCP ports that are not HTTP servers, such as
// databases, SSH, FTP, etc.
//
// This is implemented as a map for fast lookup.
var IgnoredListeningPorts = map[uint16]struct{}{
	0: {},
	// Ports 1-8 are reserved for future use by the Coder agent.
	1: {},
	2: {},
	3: {},
	4: {},
	5: {},
	6: {},
	7: {},
	8: {},
	// ftp
	20: {},
	21: {},
	// ssh
	22: {},
	// telnet
	23: {},
	// smtp
	25: {},
	// dns over TCP
	53: {},
	// pop3
	110: {},
	// imap
	143: {},
	// bgp
	179: {},
	// ldap
	389: {},
	636: {},
	// vnc
	631: {},
	// smtps
	465: {},
	// smtp
	587: {},
	// ftps
	989: {},
	990: {},
	// imaps
	993: {},
	// pop3s
	995: {},
	// mysql
	3306: {},
	// rdp
	3389: {},
	// postgres
	5432: {},
	// vnc
	5890: {},
	5990: {},
	// pprof
	6060: {},
	// redis
	6379: {},
	// code-server
	13337: {},
	// novnc
	13338: {},
	// mongodb
	27017: {},
	27018: {},
	27019: {},
	28017: {},
}

func init() {
	// Add a thousand more ports to the ignore list during tests so it's easier
	// to find an available port.
	if strings.HasSuffix(os.Args[0], ".test") {
		for i := 63000; i < 64000; i++ {
			IgnoredListeningPorts[uint16(i)] = struct{}{}
		}
	}
}

// ReconnectingPTYRequest is sent from the client to the server
// to pipe data to a PTY.
type ReconnectingPTYRequest struct {
	Data   string `json:"data"`
	Height uint16 `json:"height"`
	Width  uint16 `json:"width"`
}

type AgentConn struct {
	*tailnet.Conn
	CloseFunc func()
}

func (c *AgentConn) AwaitReachable(ctx context.Context) bool {
	return c.Conn.AwaitReachable(ctx, TailnetIP)
}

// Ping pings the agent and returns the round-trip time.
// The bool returns true if the ping was made P2P.
func (c *AgentConn) Ping(ctx context.Context) (time.Duration, bool, *ipnstate.PingResult, error) {
	return c.Conn.Ping(ctx, TailnetIP)
}

func (c *AgentConn) CloseWithError(_ error) error {
	return c.Close()
}

func (c *AgentConn) Close() error {
	if c.CloseFunc != nil {
		c.CloseFunc()
	}
	return c.Conn.Close()
}

type ReconnectingPTYInit struct {
	ID      uuid.UUID
	Height  uint16
	Width   uint16
	Command string
}

func (c *AgentConn) ReconnectingPTY(ctx context.Context, id uuid.UUID, height, width uint16, command string) (net.Conn, error) {
	conn, err := c.DialContextTCP(ctx, netip.AddrPortFrom(TailnetIP, TailnetReconnectingPTYPort))
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(ReconnectingPTYInit{
		ID:      id,
		Height:  height,
		Width:   width,
		Command: command,
	})
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	data = append(make([]byte, 2), data...)
	binary.LittleEndian.PutUint16(data, uint16(len(data)-2))

	_, err = conn.Write(data)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}

func (c *AgentConn) SSH(ctx context.Context) (net.Conn, error) {
	return c.DialContextTCP(ctx, netip.AddrPortFrom(TailnetIP, TailnetSSHPort))
}

// SSHClient calls SSH to create a client that uses a weak cipher
// for high throughput.
func (c *AgentConn) SSHClient(ctx context.Context) (*ssh.Client, error) {
	netConn, err := c.SSH(ctx)
	if err != nil {
		return nil, xerrors.Errorf("ssh: %w", err)
	}
	sshConn, channels, requests, err := ssh.NewClientConn(netConn, "localhost:22", &ssh.ClientConfig{
		// SSH host validation isn't helpful, because obtaining a peer
		// connection already signifies user-intent to dial a workspace.
		// #nosec
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, xerrors.Errorf("ssh conn: %w", err)
	}
	return ssh.NewClient(sshConn, channels, requests), nil
}

func (c *AgentConn) Speedtest(ctx context.Context, direction speedtest.Direction, duration time.Duration) ([]speedtest.Result, error) {
	speedConn, err := c.DialContextTCP(ctx, netip.AddrPortFrom(TailnetIP, TailnetSpeedtestPort))
	if err != nil {
		return nil, xerrors.Errorf("dial speedtest: %w", err)
	}
	results, err := speedtest.RunClientWithConn(direction, duration, speedConn)
	if err != nil {
		return nil, xerrors.Errorf("run speedtest: %w", err)
	}
	return results, err
}

func (c *AgentConn) DialContext(ctx context.Context, network string, addr string) (net.Conn, error) {
	if network == "unix" {
		return nil, xerrors.New("network must be tcp or udp")
	}
	_, rawPort, _ := net.SplitHostPort(addr)
	port, _ := strconv.ParseUint(rawPort, 10, 16)
	ipp := netip.AddrPortFrom(TailnetIP, uint16(port))
	if network == "udp" {
		return c.Conn.DialContextUDP(ctx, ipp)
	}
	return c.Conn.DialContextTCP(ctx, ipp)
}

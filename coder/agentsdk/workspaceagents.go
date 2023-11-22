package agentsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/gage-technologies/gigo-lib/db/models"

	"golang.org/x/xerrors"
	"nhooyr.io/websocket"

	"cdr.dev/slog"
	"github.com/coder/retry"
)

// InitializeWorkspaceAgent fetches metadata for the currently authenticated workspace agent.
func (c *AgentClient) InitializeWorkspaceAgent(ctx context.Context, isVnc bool) (WorkspaceAgentMetadata, error) {
	isVNCReq := InitializeWorkspaceAgentRequest{IsVNC: isVnc}
	res, err := c.Request(ctx, http.MethodPost, "/internal/v1/ws/initialize", isVNCReq)
	if err != nil {
		return WorkspaceAgentMetadata{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return WorkspaceAgentMetadata{}, readBodyAsError(res)
	}
	var agentMetadata WorkspaceAgentMetadata
	err = json.NewDecoder(res.Body).Decode(&agentMetadata)
	if err != nil {
		return WorkspaceAgentMetadata{}, err
	}

	// accessingPort := c.URL.Port()
	// if accessingPort == "" {
	// 	accessingPort = "80"
	// 	if c.URL.Scheme == "https" {
	// 		accessingPort = "443"
	// 	}
	// }
	// accessPort, err := strconv.Atoi(accessingPort)
	// if err != nil {
	// 	return WorkspaceAgentMetadata{}, xerrors.Errorf("convert accessing port %q: %w", accessingPort, err)
	// }
	// // Agents can provide an arbitrary access URL that may be different
	// // that the globally configured one. This breaks the built-in DERP,
	// // which would continue to reference the global access URL.
	// //
	// // This converts all built-in DERPs to use the access URL that the
	// // metadata request was performed with.
	// for _, region := range agentMetadata.DERPMap.Regions {
	// 	if !region.EmbeddedRelay {
	// 		continue
	// 	}
	//
	// 	for _, node := range region.Nodes {
	// 		if node.STUNOnly {
	// 			continue
	// 		}
	// 		node.HostName = c.URL.Hostname()
	// 		node.DERPPort = accessPort
	// 		node.ForceHTTP = c.URL.Scheme == "http"
	// 	}
	// }

	return agentMetadata, nil
}

func (c *AgentClient) ListenWorkspaceAgent(ctx context.Context) (net.Conn, error) {
	coordinateURL, err := c.URL.Parse("/internal/v1/ws/coordinate")
	if err != nil {
		return nil, xerrors.Errorf("parse url: %w", err)
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, xerrors.Errorf("create cookie jar: %w", err)
	}
	auth := c.SessionAuth()
	jar.SetCookies(coordinateURL, []*http.Cookie{
		{
			Name:  WorkspaceIDHeader,
			Value: fmt.Sprintf("%d", auth.WorkspaceID),
		},
		{
			Name:  AgentTokenHeader,
			Value: auth.Token,
		},
	})
	httpClient := &http.Client{
		Jar:       jar,
		Transport: c.HTTPClient.Transport,
	}
	// nolint:bodyclose
	conn, res, err := websocket.Dial(ctx, coordinateURL.String(), &websocket.DialOptions{
		HTTPClient: httpClient,
	})
	if err != nil {
		if res == nil {
			return nil, err
		}
		return nil, readBodyAsError(res)
	}

	// Ping once every 30 seconds to ensure that the websocket is alive. If we
	// don't get a response within 30s we kill the websocket and reconnect.
	// See: https://github.com/coder/coder/pull/5824
	go func() {
		tick := 30 * time.Second
		ticker := time.NewTicker(tick)
		defer ticker.Stop()
		defer func() {
			c.Logger.Debug(ctx, "coordinate pinger exited")
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case start := <-ticker.C:
				ctx, cancel := context.WithTimeout(ctx, tick)

				err := conn.Ping(ctx)
				if err != nil {
					c.Logger.Error(ctx, "workspace agent coordinate ping", slog.Error(err))

					err := conn.Close(websocket.StatusGoingAway, "Ping failed")
					if err != nil {
						c.Logger.Error(ctx, "close workspace agent coordinate websocket", slog.Error(err))
					}

					cancel()
					return
				}

				c.Logger.Debug(ctx, "got coordinate pong", slog.F("took", time.Since(start)))
				cancel()
			}
		}
	}()

	return websocket.NetConn(ctx, conn, websocket.MessageBinary), nil
}

func (c *AgentClient) PostWorkspaceAgentVersion(ctx context.Context, version string) error {
	versionReq := PostWorkspaceAgentVersionRequest{Version: version}
	res, err := c.Request(ctx, http.MethodPost, "/internal/v1/ws/version", versionReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return readBodyAsError(res)
	}
	return nil
}

func (c *AgentClient) PostAgentStats(ctx context.Context, stats *AgentStats) (AgentStatsResponse, error) {
	res, err := c.Request(ctx, http.MethodPost, "/internal/v1/ws/stats", stats)
	if err != nil {
		return AgentStatsResponse{}, xerrors.Errorf("send request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return AgentStatsResponse{}, readBodyAsError(res)
	}

	var interval AgentStatsResponse
	err = json.NewDecoder(res.Body).Decode(&interval)
	if err != nil {
		return AgentStatsResponse{}, xerrors.Errorf("decode stats response: %w", err)
	}

	return interval, nil
}

// AgentReportStats begins a stat streaming connection with the Coder server.
// It is resilient to network failures and intermittent coderd issues.
func (c *AgentClient) AgentReportStats(
	ctx context.Context,
	log slog.Logger,
	statsChan <-chan *AgentStats,
	setInterval func(time.Duration),
) (io.Closer, error) {
	var interval time.Duration
	ctx, cancel := context.WithCancel(ctx)
	exited := make(chan struct{})

	postStat := func(stat *AgentStats) {
		var nextInterval time.Duration
		for r := retry.New(100*time.Millisecond, time.Minute); r.Wait(ctx); {
			resp, err := c.PostAgentStats(ctx, stat)
			if err != nil {
				if !xerrors.Is(err, context.Canceled) {
					log.Error(ctx, "report stats", slog.Error(err))
				}
				continue
			}

			nextInterval = resp.ReportInterval
			break
		}

		if nextInterval != 0 && interval != nextInterval {
			setInterval(nextInterval)
		}
		interval = nextInterval
	}

	// Send an empty stat to get the interval.
	postStat(&AgentStats{ConnsByProto: map[string]int64{}})

	go func() {
		defer close(exited)

		for {
			select {
			case <-ctx.Done():
				return
			case stat, ok := <-statsChan:
				if !ok {
					return
				}

				postStat(stat)
			}
		}
	}()

	return closeFunc(func() error {
		cancel()
		<-exited
		return nil
	}), nil
}

func (c *AgentClient) PostWorkspaceAgentState(ctx context.Context, state models.WorkspaceAgentState) error {
	stateReq := PostWorkspaceAgentState{State: state}
	res, err := c.Request(ctx, http.MethodPost, "/internal/v1/ws/state", stateReq)
	if err != nil {
		return xerrors.Errorf("agent state post request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		return readBodyAsError(res)
	}

	return nil
}

func (c *AgentClient) PostAgentPorts(ctx context.Context, req *AgentPorts) error {
	res, err := c.Request(ctx, http.MethodPost, "/internal/v1/ws/ports", req)
	if err != nil {
		return xerrors.Errorf("send request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		return readBodyAsError(res)
	}
	return nil
}

package zitimesh

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/gage-technologies/gigo-lib/config"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/stretchr/testify/assert"
	"gopkg.in/natefinch/lumberjack.v2"
)

func testLaunchLocalDevServer() {
	// launch a local dev server running at localhost:42435
	// that will echo back the request body
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.Copy(w, r.Body)
		})
		_ = http.ListenAndServe("localhost:42435", nil)
	}()
}

func TestAgent(t *testing.T) {
	testLaunchLocalDevServer()

	// Set up Ziti configuration
	zitiConfig := config.ZitiConfig{
		ManagementUser: "gigo-dev",
		ManagementPass: "gigo-dev",
		EdgeHost:       "gigo-dev-ziti-controller:1280",
		EdgeBasePath:   "/",
		EdgeSchemes:    []string{"https"},
	}

	// create a ziti manager
	manager, err := NewManager(zitiConfig)
	assert.NoError(t, err)

	// setup zitinet
	defer manager.DeleteWorkspaceServicePolicy()
	err = manager.CreateWorkspaceServicePolicy()
	assert.NoError(t, err)

	// Test creating a workspace service
	defer manager.DeleteWorkspaceService()
	svcName, err := manager.CreateWorkspaceService()
	assert.NoError(t, err)

	defer manager.DeleteServer(169)
	serverId, serverToken, err := manager.CreateServer(169)
	assert.NoError(t, err)
	assert.NotEmpty(t, serverToken)
	assert.NotEmpty(t, serverId)

	defer manager.DeleteAgent(1420)
	agentId, agentToken, err := manager.CreateAgent(1420)
	assert.NoError(t, err)
	assert.NotEmpty(t, agentToken)
	assert.NotEmpty(t, agentId)

	if agentId == "" || agentToken == "" || serverId == "" || serverToken == "" {
		return
	}

	// create a new agent instance
	logWriter := &lumberjack.Logger{
		Filename: filepath.Join(os.TempDir(), "gigo-agent-ziti-test.log"),
		MaxSize:  5, // MB
	}
	defer logWriter.Close()
	logger := slog.Make(sloghuman.Sink(os.Stdout), sloghuman.Sink(logWriter)).Leveled(slog.LevelDebug)
	_, identity, err := NewAgentFromToken(context.TODO(), agentId, agentToken, logger)
	assert.NoError(t, err)
	agent, err := NewAgent(context.TODO(), agentId, identity, logger)
	assert.NoError(t, err)

	// create a new ziti context for the server side
	// enroll the identity into a configuration
	zitiCtxConfig, err := EnrollIdentity(serverToken)
	assert.NoError(t, err)

	// create a new Ziti context
	zitiCtx, err := ziti.NewContext(zitiCtxConfig)
	assert.NoError(t, err)

	// err = ctx.RefreshServices()
	assert.NoError(t, err)

	conn, err := zitiCtx.DialWithOptions(svcName, &ziti.DialOptions{
		AppData:  []byte(`{"network":"tcp","port":42435}`),
		Identity: "gigo-ws-agent-1420",
	})
	assert.NoError(t, err)

	// make a http request to the server over the ziti connection
	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (netConn net.Conn, e error) {
				return conn, nil
			},
		},
	}
	resp, err := client.Post("http://localhost:42435/", "text/plain", bytes.NewBuffer([]byte("hello ziti")))
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// read the response body
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "hello ziti", buf.String())

	// retrieve the network stats
	stats := agent.GetNetworkStats()
	b, _ := json.Marshal(stats)
	fmt.Printf("%s\n", string(b))

	// clear the stats
	agent.ClearStats()

	// retrieve the network stats
	stats = agent.GetNetworkStats()
	b, _ = json.Marshal(stats)
	fmt.Printf("%s\n", string(b))
}

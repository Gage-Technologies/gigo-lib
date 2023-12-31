package tailnettest

import (
	"crypto/tls"
	"github.com/gage-technologies/gigo-lib/logging"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"tailscale.com/derp"
	"tailscale.com/derp/derphttp"
	"tailscale.com/net/stun/stuntest"
	"tailscale.com/tailcfg"
	"tailscale.com/types/key"
	tslogger "tailscale.com/types/logger"
	"tailscale.com/types/nettype"
)

// RunDERPAndSTUN creates a DERP mapping for tests.
func RunDERPAndSTUN(t *testing.T) *tailcfg.DERPMap {
	logger, err := logging.CreateBasicLogger(
		logging.NewDefaultBasicLoggerOptions("/tmp/tailnettest-run-derp&stun.log"),
	)
	if err != nil {
		t.Fatal(err)
	}
	d := derp.NewServer(key.NewNode(), logger.TailscaleDebugLogger(""))
	server := httptest.NewUnstartedServer(derphttp.Handler(d))
	server.Config.ErrorLog = tslogger.StdLogger(logger.TailscaleDebugLogger(""))
	server.Config.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
	server.StartTLS()

	stunAddr, stunCleanup := stuntest.ServeWithPacketListener(t, nettype.Std{})
	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
		d.Close()
		stunCleanup()
	})
	tcpAddr, ok := server.Listener.Addr().(*net.TCPAddr)
	if !ok {
		t.FailNow()
	}

	return &tailcfg.DERPMap{
		Regions: map[int]*tailcfg.DERPRegion{
			1: {
				RegionID:   1,
				RegionCode: "test",
				RegionName: "Test",
				Nodes: []*tailcfg.DERPNode{
					{
						Name:             "t2",
						RegionID:         1,
						IPv4:             "127.0.0.1",
						IPv6:             "none",
						STUNPort:         stunAddr.Port,
						DERPPort:         tcpAddr.Port,
						InsecureForTests: true,
					},
				},
			},
		},
	}
}

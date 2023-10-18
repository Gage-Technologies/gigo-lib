package tailnet_test

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/gage-technologies/gigo-lib/logging"
	"net/netip"
	"testing"

	"github.com/gage-technologies/gigo-lib/coder/tailnet"
	"github.com/gage-technologies/gigo-lib/coder/tailnet/tailnettest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: the go leak check fails on the logger - fix this
// func TestMain(m *testing.M) {
// 	goleak.VerifyTestMain(m)
// }

func TestTailnet(t *testing.T) {
	sf, err := snowflake.NewNode(0)
	if err != nil {
		t.Fatal(err)
	}
	derpMap := tailnettest.RunDERPAndSTUN(t)
	t.Run("InstantClose", func(t *testing.T) {
		logger, err := logging.CreateBasicLogger(logging.NewDefaultBasicLoggerOptions("/tmp/tailnet-test-instant-close.log"))
		if err != nil {
			t.Fatal(err)
		}
		conn, err := tailnet.NewConn(tailnet.ConnTypeServer, &tailnet.Options{
			NodeID:    sf.Generate().Int64(),
			Addresses: []netip.Prefix{netip.PrefixFrom(tailnet.IP(), 128)},
			DERPMap:   derpMap,
		}, logger)
		require.NoError(t, err)
		err = conn.Close()
		require.NoError(t, err)
	})
	t.Run("Connect", func(t *testing.T) {
		logger, err := logging.CreateBasicLogger(logging.NewDefaultBasicLoggerOptions("/tmp/tailnet-test-connect.log"))
		if err != nil {
			t.Fatal(err)
		}
		w1IP := tailnet.IP()
		w1, err := tailnet.NewConn(tailnet.ConnTypeAgent, &tailnet.Options{
			NodeID:    sf.Generate().Int64(),
			Addresses: []netip.Prefix{netip.PrefixFrom(w1IP, 128)},
			DERPMap:   derpMap,
		}, logger)
		require.NoError(t, err)

		w2, err := tailnet.NewConn(tailnet.ConnTypeServer, &tailnet.Options{
			NodeID:    sf.Generate().Int64(),
			Addresses: []netip.Prefix{netip.PrefixFrom(tailnet.IP(), 128)},
			DERPMap:   derpMap,
		}, logger)
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = w1.Close()
			_ = w2.Close()
		})
		w1.SetNodeCallback(func(node *tailnet.Node) {
			err := w2.UpdateNodes([]*tailnet.Node{node})
			require.NoError(t, err)
		})
		w2.SetNodeCallback(func(node *tailnet.Node) {
			err := w1.UpdateNodes([]*tailnet.Node{node})
			require.NoError(t, err)
		})
		require.True(t, w2.AwaitReachable(context.Background(), w1IP))
		conn := make(chan struct{})
		go func() {
			listener, err := w1.Listen("tcp", ":35565")
			assert.NoError(t, err)
			defer listener.Close()
			nc, err := listener.Accept()
			if !assert.NoError(t, err) {
				return
			}
			_ = nc.Close()
			conn <- struct{}{}
		}()

		nc, err := w2.DialContextTCP(context.Background(), netip.AddrPortFrom(w1IP, 35565))
		require.NoError(t, err)
		_ = nc.Close()
		<-conn

		w1.Close()
		w2.Close()
	})
}

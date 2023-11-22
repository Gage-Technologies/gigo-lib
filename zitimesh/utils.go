package zitimesh

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"

	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/enroll"
)

func EnrollIdentity(token string) (*ziti.Config, error) {
	// parse the identity token
	tkn, _, err := enroll.ParseToken(strings.TrimSpace(token))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// enroll the identity into a configuration
	conf, err := enroll.Enroll(enroll.EnrollmentFlags{
		Token:  tkn,
		KeyAlg: "RSA",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to enroll identity: %w", err)
	}

	return conf, nil
}

func recordStats(bytes, packets int, in bool, stat *NetworkStats) {
	if in {
		stat.BytesIn += int64(bytes)
		stat.PacketsIn += int64(packets)
	} else {
		stat.BytesOut += int64(bytes)
		stat.PacketsOut += int64(packets)
	}
}

func copyWithStats(dst io.Writer, src io.Reader, srcRemote bool, netType NetworkType, port int, mu *sync.Mutex, stats *GlobalStats) (int64, error) {
	defer func() {
		r := recover()
		if r != nil {
			// retrieve the stack trace
			stack := make([]byte, 1024*8)
			n := runtime.Stack(stack[:], true)
			fmt.Printf("recovered from panic: %v\n%s\n", r, stack[:n])
		}
	}()

	buf := make([]byte, 32*1024) // 32KB buffer
	var totalBytes int64
	for {
		n, err := src.Read(buf)
		if n > 0 {
			totalBytes += int64(n)
			mu.Lock()
			// record global stats
			recordStats(n, 1, srcRemote, stats.Total)
			// record stats for the network type
			if netStats, ok := stats.ByNetwork[netType]; ok {
				recordStats(n, 1, srcRemote, netStats)
			} else {
				// create a new stat for this network type
				netStats := &NetworkStats{}
				stats.ByNetwork[netType] = netStats
				recordStats(n, 1, srcRemote, netStats)
			}
			// record stats for the port
			if portStats, ok := stats.ByPort[port]; ok {
				recordStats(n, 1, srcRemote, portStats.Stats)
			} else {
				// create a new stat for this network type
				portStats := &PortStats{
					Port:        port,
					NetworkType: netType,
					Stats:       &NetworkStats{},
				}
				stats.ByPort[port] = portStats
				recordStats(n, 1, srcRemote, portStats.Stats)
			}
			mu.Unlock()
			nw, writeErr := dst.Write(buf[:n])
			if writeErr != nil {
				return totalBytes, writeErr
			}
			if n != nw {
				return totalBytes, io.ErrShortWrite
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return totalBytes, err
		}
	}
	return totalBytes, nil
}

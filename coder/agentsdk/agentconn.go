package agentsdk

import (
	"os"
	"strings"
)

const (
	ZitiInitConnPort        = 42069
	ZitiSSHPort             = 42070
	ZitiReconnectingPTYPort = 42071
	ZitiSpeedtestPort       = 42072
	ZitiStatisticsPort      = 42073
	ZitiSerialPort          = 42074

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

var GigoReservedPorts = map[uint16]struct{}{
	ZitiInitConnPort:        {},
	ZitiSSHPort:             {},
	ZitiReconnectingPTYPort: {},
	ZitiSpeedtestPort:       {},
	ZitiStatisticsPort:      {},
	ZitiSerialPort:          {},
	42075:                   {},
	42076:                   {},
	42077:                   {},
	42078:                   {},
}

// IgnoredListeningPorts contains a list of ports in the global ignore list.
// This list contains common TCP ports that are not HTTP servers, such as
// databases, SSH, FTP, etc.
//
// This is implemented as a map for fast lookup.
var IgnoredListeningPorts = map[uint16]struct{}{
	0: {},
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
	// gigo agent operational ports
	42069: {},
	42070: {},
	42071: {},
	42072: {},
	42073: {},
	42074: {},
	42075: {},
	42076: {},
	42077: {},
	42078: {},
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

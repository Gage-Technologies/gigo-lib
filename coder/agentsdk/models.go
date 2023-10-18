package agentsdk

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gage-technologies/gigo-lib/db/models"
	"github.com/gage-technologies/gigo-lib/logging"
	"github.com/gage-technologies/gigo-lib/workspace_config"
	"tailscale.com/tailcfg"
	"time"
)

type PostWorkspaceAgentVersionRequest struct {
	Version string `json:"version"`
}

type PostWorkspaceAgentState struct {
	State models.WorkspaceAgentState `json:"state"`
}

// WorkspaceAgentConnectionInfo returns required information for establishing
// a connection with a workspace.
type WorkspaceAgentConnectionInfo struct {
	DERPMap *tailcfg.DERPMap `json:"derp_map"`
}

type DialWorkspaceAgentOptions struct {
	Logger        logging.Logger
	SnowflakeNode *snowflake.Node
	// BlockEndpoints forced a direct connection through DERP.
	BlockEndpoints bool
}

// Stats records the Agent's network connection statistics for use in
// user-facing metrics and debugging.
type AgentStats struct {
	// ConnsByProto is a count of connections by protocol.
	ConnsByProto map[string]int64 `json:"conns_by_proto"`
	// NumConns is the number of connections received by an agent.
	NumConns int64 `json:"num_comms"`
	// RxPackets is the number of received packets.
	RxPackets int64 `json:"rx_packets"`
	// RxBytes is the number of received bytes.
	RxBytes int64 `json:"rx_bytes"`
	// TxPackets is the number of transmitted bytes.
	TxPackets int64 `json:"tx_packets"`
	// TxBytes is the number of transmitted bytes.
	TxBytes int64 `json:"tx_bytes"`
}

type AgentStatsResponse struct {
	// ReportInterval is the duration after which the agent should send stats
	// again.
	ReportInterval time.Duration `json:"report_interval"`
}

type ListeningPortNetwork string

const (
	ListeningPortNetworkTCP ListeningPortNetwork = "tcp"
)

type ListeningPort struct {
	ProcessName string               `json:"process_name"` // may be empty
	Network     ListeningPortNetwork `json:"network"`      // only "tcp" at the moment
	Port        uint16               `json:"port"`
}

func (l *ListeningPort) String() string {
	return fmt.Sprintf("%s:%d:%s", l.Network, l.Port, l.ProcessName)
}

type AgentPorts struct {
	Ports []ListeningPort `json:"ports"`
}

type WorkspaceAgentMetadata struct {
	WorkspaceID        int64                                     `json:"workspace_id"`
	WorkspaceIDString  string                                    `json:"workspace_id_string"`
	Repo               string                                    `json:"repo"`
	Commit             string                                    `json:"commit"`
	GitToken           string                                    `json:"git_token"`
	GitEmail           string                                    `json:"git_email"`
	GitName            string                                    `json:"git_name"`
	Expiration         int64                                     `json:"expiration"`
	OwnerID            int64                                     `json:"owner_id"`
	OwnerIDString      string                                    `json:"owner_id_string"`
	WorkspaceSettings  *models.WorkspaceSettings                 `json:"workspace_settings"`
	VSCodePortProxyURI string                                    `json:"vscode_port_proxy_uri"`
	DERPMap            *tailcfg.DERPMap                          `json:"derpmap"`
	LastInitState      models.WorkspaceInitState                 `json:"last_init_state"`
	WorkspaceState     models.WorkspaceState                     `json:"workspace_state"`
	GigoConfig         workspace_config.GigoWorkspaceConfigAgent `json:"gigo_config"`
	UserStatus         models.UserStatus                         `json:"user_status"`
	HolidaySeason      Holiday                                   `json:"holiday_season"`
	ChallengeType      models.ChallengeType                      `json:"challenge_type"`
	UserHolidayTheme   bool                                      `json:"user_holiday_theme"`
	Hosts              map[string]string                         `json:"hosts"`
}

type Holiday int

const (
	NoHoliday Holiday = iota
	Halloween
	Christmas
	NewYears
	Valentine
	Easter
	Independence
)

func (h Holiday) String() string {
	switch h {
	case NoHoliday:
		return "no-holiday"
	case Halloween:
		return "gigo-halloween"
	case Christmas:
		return "gigo-christmas"
	case NewYears:
		return "gigo-new-years"
	case Valentine:
		return "gigo-valentines"
	case Easter:
		return "gigo-easter"
	case Independence:
		return "gigo-independence"
	default:
		return "no-holiday"
	}
}

type PostWorkspaceInitStateCompleted struct {
	State models.WorkspaceInitState `json:"state"`
}

type PostWorkspaceInitFailure struct {
	State   models.WorkspaceInitState `json:"state"`
	Command string                    `json:"command"`
	Status  int                       `json:"status"`
	Stdout  string                    `json:"stdout"`
	Stderr  string                    `json:"stderr"`
}

package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/kisielk/sqlstruct"
	"time"
)

type WorkspaceAgentStats struct {
	ID          int64     `json:"_id" sql:"_id"`
	AgentID     int64     `json:"agent_id" sql:"agent_id"`
	WorkspaceID int64     `json:"workspace_id" sql:"workspace_id"`
	Timestamp   time.Time `json:"timestamp" sql:"timestamp"`
	// ConnsByProto is a count of connections by protocol.
	ConnsByProto map[string]int64 `json:"conns_by_proto" sql:"conns_by_proto"`
	// NumConns is the number of connections received by an agent.
	NumConns int64 `json:"num_comms" sql:"num_comms"`
	// RxPackets is the number of received packets.
	RxPackets int64 `json:"rx_packets" sql:"rx_packets"`
	// RxBytes is the number of received bytes.
	RxBytes int64 `json:"rx_bytes" sql:"rx_bytes"`
	// TxPackets is the number of transmitted bytes.
	TxPackets int64 `json:"tx_packets" sql:"tx_packets"`
	// TxBytes is the number of transmitted bytes.
	TxBytes int64 `json:"tx_bytes" sql:"tx_bytes"`
}

type WorkspaceAgentStatsSQL struct {
	ID           int64     `sql:"_id"`
	AgentID      int64     `sql:"agent_id"`
	WorkspaceID  int64     `sql:"workspace_id"`
	Timestamp    time.Time `sql:"timestamp"`
	ConnsByProto []byte    `sql:"conns_by_proto"`
	NumConns     int64     `sql:"num_comms"`
	RxPackets    int64     `sql:"rx_packets"`
	RxBytes      int64     `sql:"rx_bytes"`
	TxPackets    int64     `sql:"tx_packets"`
	TxBytes      int64     `sql:"tx_bytes"`
}

func CreateWorkspaceAgentStats(id int64, agent int64, ws int64, ts time.Time, connsByProto map[string]int64, conns,
	rPacks, rBytes, tPacks, tBytes int64) *WorkspaceAgentStats {
	return &WorkspaceAgentStats{
		ID:           id,
		AgentID:      agent,
		WorkspaceID:  ws,
		Timestamp:    ts,
		ConnsByProto: connsByProto,
		NumConns:     conns,
		RxPackets:    rPacks,
		RxBytes:      rBytes,
		TxPackets:    tPacks,
		TxBytes:      tBytes,
	}
}

func WorkspaceAgentStatsFromSQLNative(rows *sql.Rows) (*WorkspaceAgentStats, error) {
	var stats WorkspaceAgentStatsSQL
	err := sqlstruct.Scan(&stats, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan workspace agent stats from rows: %v", err)
	}

	// create map to scan connections-by-protocol into
	var connsByProto map[string]int64
	err = json.Unmarshal(stats.ConnsByProto, &connsByProto)
	if err != nil {
		return nil, fmt.Errorf("failed to scan workspace agent stats from rows: %v", err)
	}

	return &WorkspaceAgentStats{
		ID:           stats.ID,
		AgentID:      stats.AgentID,
		WorkspaceID:  stats.WorkspaceID,
		Timestamp:    stats.Timestamp,
		ConnsByProto: connsByProto,
		NumConns:     stats.NumConns,
		RxPackets:    stats.RxPackets,
		RxBytes:      stats.RxBytes,
		TxPackets:    stats.TxPackets,
		TxBytes:      stats.TxBytes,
	}, nil
}

func (s *WorkspaceAgentStats) ToSQLNative() ([]*SQLInsertStatement, error) {
	// attempt to marshall connections-by-protocol to bytes
	connsByProtoBytes, err := json.Marshal(s.ConnsByProto)
	if err != nil {
		return nil, fmt.Errorf("failed to martial connections-by-protocol: %v", err)
	}

	// create slice to hold insertion statements for this workspace config and initialize the slice with the main insertion statement
	return []*SQLInsertStatement{
		{
			Statement: "insert ignore into workspace_agent_stats(_id, agent_id, workspace_id, timestamp, conns_by_proto, num_comms, rx_packets, rx_bytes, tx_packets, tx_bytes) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
			Values: []interface{}{
				s.ID, s.AgentID, s.WorkspaceID, s.Timestamp, connsByProtoBytes, s.NumConns, s.RxPackets, s.RxBytes, s.TxPackets, s.TxBytes,
			},
		},
	}, nil
}

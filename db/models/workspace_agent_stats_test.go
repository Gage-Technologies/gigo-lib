package models

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	ti "github.com/gage-technologies/gigo-lib/db"
)

func TestCreateWorkspaceAgentStats(t *testing.T) {
	id := int64(420)
	agent := int64(69)
	workspace := int64(710)
	ts := time.Now()
	cbp := map[string]int64{
		"tcp": 27483,
		"udp": 3384,
	}

	agentStats := CreateWorkspaceAgentStats(id, agent, workspace, ts, cbp, 274849, 234, 5, 345, 3466)

	if agentStats.ID != id {
		t.Fatalf("\nCreateWorkspaceAgentStats() ID does not match")
	}

	if agentStats.AgentID != agent {
		t.Fatalf("\nCreateWorkspaceAgentStats() Agent does not match")
	}

	if agentStats.WorkspaceID != workspace {
		t.Fatalf("\nCreateWorkspaceAgentStats() Workspace does not match")
	}

	if agentStats.Timestamp.Unix() != ts.Unix() {
		t.Fatalf("\nCreateWorkspaceAgentStats() Timestamp does not match")
	}

	if agentStats.ConnsByProto["tcp"] != cbp["tcp"] {
		t.Fatalf("\nCreateWorkspaceAgentStats() CBP TCP does not match")
	}

	if agentStats.ConnsByProto["udp"] != cbp["udp"] {
		t.Fatalf("\nCreateWorkspaceAgentStats() CBP UDP does not match")
	}

	if agentStats.RxPackets != 234 {
		t.Fatalf("\nCreateWorkspaceAgentStats() RxPackets does not match")
	}

	if agentStats.TxPackets != 345 {
		t.Fatalf("\nCreateWorkspaceAgentStats() TxPackets does not match")
	}

	if agentStats.RxBytes != 5 {
		t.Fatalf("\nCreateWorkspaceAgentStats() RxBytes does not match")
	}

	if agentStats.TxBytes != 3466 {
		t.Fatalf("\nCreateWorkspaceAgentStats() TxBytes does not match")
	}

	t.Log("\nCreateWorkspaceAgentStats() succeeded")
}

func TestWorkspaceAgentStatsFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nToSQlNative WorkspaceAgentStats Table\n    Error: ", err)
		return
	}

	defer db.DB.Exec("delete from workspace_agent_stats")

	id := int64(420)
	agent := int64(69)
	workspace := int64(710)
	ts := time.Now()
	cbp := map[string]int64{
		"tcp": 27483,
		"udp": 3384,
	}

	agentStats := CreateWorkspaceAgentStats(id, agent, workspace, ts, cbp, 274849, 234, 5, 345, 3466)

	statements, err := agentStats.ToSQLNative()
	if err != nil {
		t.Error("\nToSQLNative WorkspaceAgentStats failed\n    Error: ", err)
		return
	}

	if statements[0].Statement != "insert ignore into workspace_agent_stats(_id, agent_id, workspace_id, timestamp, conns_by_proto, num_comms, rx_packets, rx_bytes, tx_packets, tx_bytes) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);" {
		t.Errorf("failed workspace agent stats to sql native, err: workspace insert statement was incorrect: %v", statements[0].Statement)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		t.Error("\nTo sql native failed\n    Error: incorrect values returned for user badges table")
		return
	}

	for _, s := range statements {
		_, err = tx.Exec(s.Statement, s.Values...)
		if err != nil {
			_ = tx.Rollback()
			t.Errorf("\nTo sql native failed\n    Error: %v", err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		t.Errorf(fmt.Sprintf("failed to commit transaction, err: %v", err))
		return
	}

	res, err := db.DB.Query("select * from workspace_agent_stats where _id = ?", agentStats.ID)
	if err != nil {
		t.Errorf("\nWorkspaceAgentStatsFromSQLNative failed\n    Error: failed to query for workspace: %v", err)
	}

	res.Next()

	a2, err := WorkspaceAgentStatsFromSQLNative(res)
	if err != nil {
		t.Errorf("\nWorkspaceAgentStatsFromSQLNative failed\n    Error: failed to load workspace: %v", err)
	}

	if a2 == nil {
		t.Error("\nWorkspaceAgentStatsFromSQLNative failed\n    Error: load from sql is nil")
		return
	}

	if math.Abs(float64(agentStats.Timestamp.Unix()-a2.Timestamp.Unix())) > 3 {
		t.Fatalf("\nWorkspaceAgentStatsFromSQLNative failed\n    Error: timestamp does not match")
	}

	agentStats.Timestamp = time.Unix(0, 0)
	a2.Timestamp = time.Unix(0, 0)

	if !reflect.DeepEqual(*agentStats, *a2) {
		t.Errorf("\nWorkspaceAgentStatsFromSQLNative failed\n    Error: incorrect data returned\n%+v\n%+v", agentStats, a2)
		return
	}

	t.Logf("WorkspaceAgentStatsFromSQLNative Succeeded")
}

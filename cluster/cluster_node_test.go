package cluster

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/gage-technologies/gigo-lib/logging"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/atomic"
)

func TestClusterNode(t *testing.T) {
	etcdCfg := etcd.Config{
		Endpoints: []string{"gigo-dev-etcd:2379"},
	}

	logger, err := logging.CreateBasicLogger(logging.NewDefaultBasicLoggerOptions("/tmp/gigo-lib-cluster-test.log"))
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	leaderExecCount := atomic.NewInt32(0)
	followerExecCount := atomic.NewInt32(0)

	leaderRoutine := func(ctx context.Context) error {
		leaderExecCount.Inc()
		return nil
	}

	followerRoutine := func(ctx context.Context) error {
		followerExecCount.Inc()
		return nil
	}

	node1, err := NewClusterNode(ClusterNodeOptions{
		ctx,
		rand.Int63(),
		"node1",
		time.Second,
		"test",
		etcdCfg,
		leaderRoutine,
		followerRoutine,
		time.Millisecond * 200,
		logger,
	})
	if err != nil {
		t.Fatal(err)
	}

	node2, err := NewClusterNode(ClusterNodeOptions{
		ctx,
		rand.Int63(),
		"node2",
		time.Second,
		"test",
		etcdCfg,
		leaderRoutine,
		followerRoutine,
		time.Millisecond * 200,
		logger,
	})
	if err != nil {
		t.Fatal(err)
	}

	node3, err := NewClusterNode(ClusterNodeOptions{
		ctx,
		rand.Int63(),
		"node3",
		time.Second,
		"test",
		etcdCfg,
		leaderRoutine,
		followerRoutine,
		time.Millisecond * 200,
		logger,
	})
	if err != nil {
		t.Fatal(err)
	}

	node1.Start()
	node2.Start()
	node3.Start()

	state1 := NodeRoleUnknown
	state2 := NodeRoleUnknown
	state3 := NodeRoleUnknown

	testCtx, testCancel := context.WithTimeout(ctx, time.Second*5)
	defer testCancel()

	for state1 == NodeRoleUnknown || state2 == NodeRoleUnknown || state3 == NodeRoleUnknown {
		select {
		case <-testCtx.Done():
			fmt.Println("testCtx done: ", testCtx.Err())
			time.Sleep(time.Hour * 30)
			t.Fatal(testCtx.Err())
		default:
		}

		if state1 == NodeRoleUnknown {
			node1.lock.Lock()
			state1 = node1.Role
			node1.lock.Unlock()
		}

		if state2 == NodeRoleUnknown {
			node2.lock.Lock()
			state2 = node2.Role
			node2.lock.Unlock()
		}

		if state3 == NodeRoleUnknown {
			node3.lock.Lock()
			state3 = node3.Role
			node3.lock.Unlock()
		}
	}

	if state1 != NodeRoleLeader && state2 != NodeRoleLeader && state3 != NodeRoleLeader {
		t.Fatal("no leader")
	}

	t.Logf("Cluster Role:\n    ClusterNode 1: %v\n    ClusterNode 2: %v\n    ClusterNode 3: %v\n", state1, state2, state3)

	time.Sleep(time.Second)

	if followerExecCount.Load() < 6 {
		t.Fatal("follower should have been executed at least 6 times: ", followerExecCount.Load())
	}

	if leaderExecCount.Load() < 3 {
		t.Fatal("leader should have been executed at least 3 times: ", leaderExecCount.Load())
	}

	// /////////////////////////////// GetNodes test
	nodesTruth := []int64{node1.ID, node2.ID, node3.ID}

	sort.Slice(nodesTruth, func(i, j int) bool {
		return nodesTruth[i] < nodesTruth[j]
	})

	nodesQuery1, err := node1.GetNodes()
	if err != nil {
		t.Fatal(err)
	}
	nodesQuery1Ids := make([]int64, len(nodesQuery1))
	for i, node := range nodesQuery1 {
		nodesQuery1Ids[i] = node.ID
	}

	sort.Slice(nodesQuery1Ids, func(i, j int) bool {
		return nodesQuery1Ids[i] < nodesQuery1Ids[j]
	})

	if !reflect.DeepEqual(nodesQuery1Ids, nodesTruth) {
		t.Fatal("nodes query is not equal to nodes truth")
	}

	nodesQuery2, err := node2.GetNodes()
	if err != nil {
		t.Fatal(err)
	}
	nodesQuery2Ids := make([]int64, len(nodesQuery2))
	for i, node := range nodesQuery2 {
		nodesQuery2Ids[i] = node.ID
	}

	sort.Slice(nodesQuery2Ids, func(i, j int) bool {
		return nodesQuery2Ids[i] < nodesQuery2Ids[j]
	})

	if !reflect.DeepEqual(nodesQuery2Ids, nodesTruth) {
		t.Fatal("nodes query is not equal to nodes truth")
	}

	nodesQuery3, err := node3.GetNodes()
	if err != nil {
		t.Fatal(err)
	}
	nodesQuery3Ids := make([]int64, len(nodesQuery3))
	for i, node := range nodesQuery3 {
		nodesQuery3Ids[i] = node.ID
	}

	sort.Slice(nodesQuery3Ids, func(i, j int) bool {
		return nodesQuery3Ids[i] < nodesQuery3Ids[j]
	})

	if !reflect.DeepEqual(nodesQuery3Ids, nodesTruth) {
		t.Fatal("nodes query is not equal to nodes truth")
	}
	// ////////////////////////////////////////////

	leader := node1
	followers := []*ClusterNode{node2, node3}
	if node2.Role == NodeRoleLeader {
		leader = node2
		followers = []*ClusterNode{node1, node3}
	}
	if node3.Role == NodeRoleLeader {
		leader = node3
		followers = []*ClusterNode{node1, node2}
	}

	// /////////////////////////////// GetLeader test
	leaderQuery, err := followers[0].GetLeader()
	if err != nil {
		t.Fatal(err)
	}

	if leaderQuery != leader.ID {
		t.Fatal("leader query is not equal to leader")
	}
	// ////////////////////////////////////////////

	// /////////////////////////////// KV test
	err = node1.Put("test", "1")
	if err != nil {
		t.Fatal(err)
	}
	err = node2.Put("test", "2")
	if err != nil {
		t.Fatal(err)
	}
	err = node3.Put("test", "3")
	if err != nil {
		t.Fatal(err)
	}

	v1, err := node1.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	if v1 != "1" {
		t.Fatal("kv value is not equal to 1")
	}

	v2, err := node2.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	if v2 != "2" {
		t.Fatal("kv value is not equal to 2")
	}

	v3, err := node3.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	if v3 != "3" {
		t.Fatal("kv value is not equal to 3")
	}

	vc, err := node3.GetCluster("test")
	if err != nil {
		t.Fatal(err)
	}
	if vc[node1.ID][0].Value != "1" {
		t.Fatal("cluster value is not equal to 1")
	}
	if vc[node2.ID][0].Value != "2" {
		t.Fatal("cluster value is not equal to 2")
	}
	if vc[node3.ID][0].Value != "3" {
		t.Fatal("cluster value is not equal to 3")
	}
	if len(vc) != 3 {
		t.Fatal("cluster value length is not equal to 3")
	}

	err = node1.Delete("test")
	if err != nil {
		t.Fatal(err)
	}

	err = node2.Delete("test")
	if err != nil {
		t.Fatal(err)
	}

	err = node3.Delete("test")
	if err != nil {
		t.Fatal(err)
	}

	v1, err = node1.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	if v1 != "" {
		t.Fatal("kv value was not deleted")
	}

	v2, err = node2.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	if v2 != "" {
		t.Fatal("kv value was not deleted")
	}

	v3, err = node3.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	if v3 != "" {
		t.Fatal("kv value was not deleted")
	}
	// ////////////////////////////////////////////

	t.Log("Killing Leader: ", leader.Address)

	leader.Stop()
	_ = leader.Close()

	testCtx, testCancel = context.WithTimeout(ctx, time.Second*5)
	defer testCancel()

	for followers[0].Role != NodeRoleLeader && followers[1].Role != NodeRoleLeader {
		select {
		case <-testCtx.Done():
			t.Fatal(testCtx.Err())
		default:
		}
	}

	t.Logf("Cluster Role:\n    ClusterNode 1: %v\n    ClusterNode 2: %v\n", followers[0].Role, followers[1].Role)

	node4, err := NewClusterNode(ClusterNodeOptions{
		ctx,
		rand.Int63(),
		"node4",
		time.Second,
		"test",
		etcdCfg,
		leaderRoutine,
		followerRoutine,
		time.Millisecond * 200,
		logger,
	})
	if err != nil {
		t.Fatal(err)
	}

	node4.Start()

	testCtx, testCancel = context.WithTimeout(ctx, time.Second*5)
	defer testCancel()

	for {
		node4.lock.Lock()
		if node4.Role == NodeRoleFollower {
			node4.lock.Unlock()
			break
		}
		node4.lock.Unlock()

		select {
		case <-testCtx.Done():
			t.Fatal(testCtx.Err())
		default:
		}
	}

	t.Logf("Cluster Role:\n    ClusterNode 1: %v\n    ClusterNode 2: %v\n    ClusterNode 4: %v\n", followers[0].Role, followers[1].Role, node4.Role)

	// /////////////////////////////// GetNodes test
	nodesTruth = []int64{followers[0].ID, followers[1].ID, node4.ID}

	sort.Slice(nodesTruth, func(i, j int) bool {
		return nodesTruth[i] < nodesTruth[j]
	})

	nodesQuery4, err := node4.GetNodes()
	if err != nil {
		t.Fatal(err)
	}

	nodesQuery4Ids := make([]int64, len(nodesQuery4))
	for i, node := range nodesQuery4 {
		nodesQuery4Ids[i] = node.ID
	}

	sort.Slice(nodesQuery4Ids, func(i, j int) bool {
		return nodesQuery4Ids[i] < nodesQuery4Ids[j]
	})

	if !reflect.DeepEqual(nodesQuery4Ids, nodesTruth) {
		fmt.Println(nodesQuery4Ids)
		fmt.Println(nodesTruth)
		t.Fatal("nodes query is not equal to nodes truth")
	}
	// ////////////////////////////////////////////

	// /////////////////////////////// WatchClusterKey test
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watchCh, err := node4.WatchKeyCluster(ctx, "test-watch")
	if err != nil {
		t.Fatal(err)
	}

	err = followers[0].Put("test-watch", "add")
	if err != nil {
		t.Fatal(err)
	}

	event := <-watchCh
	if event.Type != EventTypeAdded {
		t.Fatal("invalid event type:", event.Type)
	}
	if event.Key != fmt.Sprintf("/test/state-data/test-watch/%d", followers[0].ID) {
		t.Fatal("invalid event key:", event.Key)
	}
	if event.Value != "add" {
		t.Fatal("invalid event value:", event.Value)
	}
	if event.OldValue != "" {
		t.Fatal("invalid event old value:", event.OldValue)
	}
	if event.OldKey != "" {
		t.Fatal("invalid event old key:", event.OldKey)
	}
	if event.NodeID != followers[0].ID {
		t.Fatal("invalid event node id:", event.NodeID)
	}

	err = followers[0].Put("test-watch", "mod")
	if err != nil {
		t.Fatal(err)
	}

	event = <-watchCh
	if event.Type != EventTypeModified {
		t.Fatal("invalid event type:", event.Type)
	}
	if event.Key != fmt.Sprintf("/test/state-data/test-watch/%d", followers[0].ID) {
		t.Fatal("invalid event key:", event.Key)
	}
	if event.Value != "mod" {
		t.Fatal("invalid event value:", event.Value)
	}
	if event.OldValue != "add" {
		t.Fatal("invalid event old value:", event.OldValue)
	}
	if event.OldKey != fmt.Sprintf("/test/state-data/test-watch/%d", followers[0].ID) {
		t.Fatal("invalid event old key:", event.OldKey)
	}
	if event.NodeID != followers[0].ID {
		t.Fatal("invalid event node id:", event.NodeID)
	}

	err = followers[0].Delete("test-watch")
	if err != nil {
		t.Fatal(err)
	}

	event = <-watchCh
	if event.Type != EventTypeDeleted {
		t.Fatal("invalid event type:", event.Type)
	}
	if event.Key != "" {
		t.Fatal("invalid event key:", event.Key)
	}
	if event.Value != "" {
		t.Fatal("invalid event value:", event.Value)
	}
	if event.OldValue != "mod" {
		t.Fatal("invalid event old value:", event.OldValue)
	}
	if event.OldKey != fmt.Sprintf("/test/state-data/test-watch/%d", followers[0].ID) {
		t.Fatal("invalid event old key:", event.OldKey)
	}
	if event.NodeID != followers[0].ID {
		t.Fatal("invalid event node id:", event.NodeID)
	}
	// ////////////////////////////////////////////

}

// func TestClusterNodeStall(t *testing.T) {
// 	etcdCfg := etcd.Config{
// 		Endpoints: []string{"gigo-dev-etcd:2379"},
// 	}

// 	logger, err := logging.CreateBasicLogger(logging.NewDefaultBasicLoggerOptions("/tmp/gigo-lib-cluster-test.log"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	ctx := context.Background()

// 	leaderExecCount := atomic.NewInt32(0)
// 	followerExecCount := atomic.NewInt32(0)

// 	leaderRoutine := func(ctx context.Context) error {
// 		leaderExecCount.Inc()
// 		return nil
// 	}

// 	followerRoutine := func(ctx context.Context) error {
// 		followerExecCount.Inc()
// 		return nil
// 	}

// 	node1, err := NewClusterNode(ClusterNodeOptions{
// 		ctx,
// 		rand.Int63(),
// 		"node1",
// 		time.Second,
// 		"test",
// 		etcdCfg,
// 		leaderRoutine,
// 		followerRoutine,
// 		time.Minute,
// 		logger,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	node2, err := NewClusterNode(ClusterNodeOptions{
// 		ctx,
// 		rand.Int63(),
// 		"node2",
// 		time.Second,
// 		"test",
// 		etcdCfg,
// 		leaderRoutine,
// 		followerRoutine,
// 		time.Minute,
// 		logger,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	node3, err := NewClusterNode(ClusterNodeOptions{
// 		ctx,
// 		rand.Int63(),
// 		"node3",
// 		time.Second,
// 		"test",
// 		etcdCfg,
// 		leaderRoutine,
// 		followerRoutine,
// 		time.Minute,
// 		logger,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	node1.Start()
// 	node2.Start()
// 	node3.Start()

// 	nodes, err := node1.GetNodes()
// 	assert.NoError(t, err)
// 	b, _ := json.Marshal(nodes)
// 	fmt.Println("existing nodes: ", string(b))

// 	state1 := NodeRoleUnknown
// 	state2 := NodeRoleUnknown
// 	state3 := NodeRoleUnknown

// 	testCtx, testCancel := context.WithTimeout(ctx, time.Second*5)
// 	defer testCancel()

// 	for state1 == NodeRoleUnknown || state2 == NodeRoleUnknown || state3 == NodeRoleUnknown {
// 		select {
// 		case <-testCtx.Done():
// 			fmt.Println("testCtx done: ", testCtx.Err())
// 			time.Sleep(time.Hour * 30)
// 			t.Fatal(testCtx.Err())
// 		default:
// 		}

// 		if state1 == NodeRoleUnknown {
// 			node1.lock.Lock()
// 			state1 = node1.Role
// 			node1.lock.Unlock()
// 		}

// 		if state2 == NodeRoleUnknown {
// 			node2.lock.Lock()
// 			state2 = node2.Role
// 			node2.lock.Unlock()
// 		}

// 		if state3 == NodeRoleUnknown {
// 			node3.lock.Lock()
// 			state3 = node3.Role
// 			node3.lock.Unlock()
// 		}
// 	}

// 	t.Logf("Cluster Role:\n    ClusterNode 1: %v\n    ClusterNode 2: %v\n    ClusterNode 3: %v\n", state1, state2, state3)

// 	idx := 0
// 	for {
// 		fmt.Printf("===================================== %d\n", idx)
// 		fmt.Println("Node 1: ", node1.Role)
// 		fmt.Println("Node 2: ", node2.Role)
// 		fmt.Println("Node 3: ", node3.Role)
// 		time.Sleep(time.Second)
// 		idx++
// 	}
// }

/*
Copyright 2023 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package discovery

import (
	"context"
	"encoding/hex"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"mdibaiee/vitess/go/test/utils"
	"mdibaiee/vitess/go/vt/topo"
	"mdibaiee/vitess/go/vt/topo/faketopo"

	querypb "mdibaiee/vitess/go/vt/proto/query"
	topodatapb "mdibaiee/vitess/go/vt/proto/topodata"
	vschemapb "mdibaiee/vitess/go/vt/proto/vschema"
)

func TestSrvKeyspaceWithNilNewKeyspace(t *testing.T) {
	ctx := utils.LeakCheckContext(t)
	cell := "cell"
	keyspace := "testks"
	factory := faketopo.NewFakeTopoFactory()
	factory.AddCell(cell)
	ts := faketopo.NewFakeTopoServer(ctx, factory)
	ts2 := &fakeTopoServer{}
	hc := NewHealthCheck(ctx, 1*time.Millisecond, time.Hour, ts, cell, "")
	defer hc.Close()
	kew := NewKeyspaceEventWatcher(ctx, ts2, hc, cell)
	kss := &keyspaceState{
		kew:      kew,
		keyspace: keyspace,
		shards:   make(map[string]*shardState),
	}
	kss.lastKeyspace = &topodatapb.SrvKeyspace{}
	require.True(t, kss.onSrvKeyspace(nil, nil))
}

// TestKeyspaceEventConcurrency confirms that the keyspace event watcher
// does not fail to broadcast received keyspace events to subscribers.
// This verifies that no events are lost when there's a high number of
// concurrent keyspace events.
func TestKeyspaceEventConcurrency(t *testing.T) {
	cell := "cell1"
	factory := faketopo.NewFakeTopoFactory()
	factory.AddCell(cell)
	sts := &fakeTopoServer{}
	hc := NewFakeHealthCheck(make(chan *TabletHealth))
	defer hc.Close()
	kew := &KeyspaceEventWatcher{
		hc:        hc,
		ts:        sts,
		localCell: cell,
		keyspaces: make(map[string]*keyspaceState),
		subs:      make(map[chan *KeyspaceEvent]struct{}),
	}

	// Subscribe to the watcher's broadcasted keyspace events.
	receiver := kew.Subscribe()

	updates := atomic.Uint32{}
	updates.Store(0)
	wg := sync.WaitGroup{}
	concurrency := 100
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-receiver:
				updates.Add(1)
			}
		}
	}()
	// Start up concurent go-routines that will broadcast keyspace events.
	for i := 1; i <= concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			kew.broadcast(&KeyspaceEvent{})
		}()
	}
	wg.Wait()
	for {
		select {
		case <-ctx.Done():
			require.Equal(t, concurrency, int(updates.Load()), "expected %d updates, got %d", concurrency, updates.Load())
			return
		default:
			if int(updates.Load()) == concurrency { // Pass
				cancel()
				return
			}
		}
	}
}

// TestKeyspaceEventTypes confirms that the keyspace event watcher determines
// that the unavailability event is caused by the correct scenario. We should
// consider it to be caused by a resharding operation when the following
// conditions are present:
// 1. The keyspace is inconsistent (in the middle of an availability event)
// 2. The target tablet is a primary
// 3. The keyspace has overlapping shards
// 4. The overlapping shard's tablet is serving
// And we should consider the cause to be a primary not serving when the
// following conditions exist:
// 1. The keyspace is inconsistent (in the middle of an availability event)
// 2. The target tablet is a primary
// 3. The target tablet is not serving
// 4. The shard's externallyReparented time is not 0
// 5. The shard's currentPrimary state is not nil
// We should never consider both as a possible cause given the same
// keyspace state.
func TestKeyspaceEventTypes(t *testing.T) {
	utils.EnsureNoLeaks(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cell := "cell"
	keyspace := "testks"
	factory := faketopo.NewFakeTopoFactory()
	factory.AddCell(cell)
	ts := faketopo.NewFakeTopoServer(ctx, factory)
	ts2 := &fakeTopoServer{}
	hc := NewHealthCheck(ctx, 1*time.Millisecond, time.Hour, ts, cell, "")
	defer hc.Close()
	kew := NewKeyspaceEventWatcher(ctx, ts2, hc, cell)

	type testCase struct {
		name                    string
		kss                     *keyspaceState
		shardToCheck            string
		expectResharding        bool
		expectPrimaryNotServing bool
	}

	testCases := []testCase{
		{
			name: "one to two resharding in progress",
			kss: &keyspaceState{
				kew:      kew,
				keyspace: keyspace,
				shards: map[string]*shardState{
					"-": {
						target: &querypb.Target{
							Keyspace:   keyspace,
							Shard:      "-",
							TabletType: topodatapb.TabletType_PRIMARY,
						},
						serving: false,
					},
					"-80": {
						target: &querypb.Target{
							Keyspace:   keyspace,
							Shard:      "-80",
							TabletType: topodatapb.TabletType_PRIMARY,
						},
						serving: true,
					},
					"80-": {
						target: &querypb.Target{
							Keyspace:   keyspace,
							Shard:      "80-",
							TabletType: topodatapb.TabletType_PRIMARY,
						},
						serving: false,
					},
				},
				consistent: false,
			},
			shardToCheck:            "-",
			expectResharding:        true,
			expectPrimaryNotServing: false,
		},
		{
			name: "two to four resharding in progress",
			kss: &keyspaceState{
				kew:      kew,
				keyspace: keyspace,
				shards: map[string]*shardState{
					"-80": {
						target: &querypb.Target{
							Keyspace:   keyspace,
							Shard:      "-80",
							TabletType: topodatapb.TabletType_PRIMARY,
						},
						serving: false,
					},
					"80-": {
						target: &querypb.Target{
							Keyspace:   keyspace,
							Shard:      "80-",
							TabletType: topodatapb.TabletType_PRIMARY,
						},
						serving: true,
					},
					"-40": {
						target: &querypb.Target{
							Keyspace:   keyspace,
							Shard:      "-40",
							TabletType: topodatapb.TabletType_PRIMARY,
						},
						serving: true,
					},
					"40-80": {
						target: &querypb.Target{
							Keyspace:   keyspace,
							Shard:      "40-80",
							TabletType: topodatapb.TabletType_PRIMARY,
						},
						serving: true,
					},
					"80-c0": {
						target: &querypb.Target{
							Keyspace:   keyspace,
							Shard:      "80-c0",
							TabletType: topodatapb.TabletType_PRIMARY,
						},
						serving: false,
					},
					"c0-": {
						target: &querypb.Target{
							Keyspace:   keyspace,
							Shard:      "c0-",
							TabletType: topodatapb.TabletType_PRIMARY,
						},
						serving: false,
					},
				},
				consistent: false,
			},
			shardToCheck:            "-80",
			expectResharding:        true,
			expectPrimaryNotServing: false,
		},
		{
			name: "unsharded primary not serving",
			kss: &keyspaceState{
				kew:      kew,
				keyspace: keyspace,
				shards: map[string]*shardState{
					"-": {
						target: &querypb.Target{
							Keyspace:   keyspace,
							Shard:      "-",
							TabletType: topodatapb.TabletType_PRIMARY,
						},
						serving:              false,
						externallyReparented: time.Now().UnixNano(),
						currentPrimary: &topodatapb.TabletAlias{
							Cell: cell,
							Uid:  100,
						},
					},
				},
				consistent: false,
			},
			shardToCheck:            "-",
			expectResharding:        false,
			expectPrimaryNotServing: true,
		},
		{
			name: "sharded primary not serving",
			kss: &keyspaceState{
				kew:      kew,
				keyspace: keyspace,
				shards: map[string]*shardState{
					"-80": {
						target: &querypb.Target{
							Keyspace:   keyspace,
							Shard:      "-80",
							TabletType: topodatapb.TabletType_PRIMARY,
						},
						serving:              false,
						externallyReparented: time.Now().UnixNano(),
						currentPrimary: &topodatapb.TabletAlias{
							Cell: cell,
							Uid:  100,
						},
					},
					"80-": {
						target: &querypb.Target{
							Keyspace:   keyspace,
							Shard:      "80-",
							TabletType: topodatapb.TabletType_PRIMARY,
						},
						serving: true,
					},
				},
				consistent: false,
			},
			shardToCheck:            "-80",
			expectResharding:        false,
			expectPrimaryNotServing: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kew.mu.Lock()
			kew.keyspaces[keyspace] = tc.kss
			kew.mu.Unlock()

			require.NotNil(t, tc.kss.shards[tc.shardToCheck], "the specified shardToCheck of %q does not exist in the shardState", tc.shardToCheck)

			resharding := kew.TargetIsBeingResharded(ctx, tc.kss.shards[tc.shardToCheck].target)
			require.Equal(t, resharding, tc.expectResharding, "TargetIsBeingResharded should return %t", tc.expectResharding)

			_, primaryDown := kew.PrimaryIsNotServing(ctx, tc.kss.shards[tc.shardToCheck].target)
			require.Equal(t, primaryDown, tc.expectPrimaryNotServing, "PrimaryIsNotServing should return %t", tc.expectPrimaryNotServing)
		})
	}
}

type fakeTopoServer struct {
}

// GetTopoServer returns the full topo.Server instance.
func (f *fakeTopoServer) GetTopoServer() (*topo.Server, error) {
	return nil, nil
}

// GetSrvKeyspaceNames returns the list of keyspaces served in
// the provided cell.
func (f *fakeTopoServer) GetSrvKeyspaceNames(ctx context.Context, cell string, staleOK bool) ([]string, error) {
	return []string{"ks1"}, nil
}

// GetSrvKeyspace returns the SrvKeyspace for a cell/keyspace.
func (f *fakeTopoServer) GetSrvKeyspace(ctx context.Context, cell, keyspace string) (*topodatapb.SrvKeyspace, error) {
	zeroHexBytes, _ := hex.DecodeString("")
	eightyHexBytes, _ := hex.DecodeString("80")
	ks := &topodatapb.SrvKeyspace{
		Partitions: []*topodatapb.SrvKeyspace_KeyspacePartition{
			{
				ServedType: topodatapb.TabletType_PRIMARY,
				ShardReferences: []*topodatapb.ShardReference{
					{Name: "-80", KeyRange: &topodatapb.KeyRange{Start: zeroHexBytes, End: eightyHexBytes}},
					{Name: "80-", KeyRange: &topodatapb.KeyRange{Start: eightyHexBytes, End: zeroHexBytes}},
				},
			},
		},
	}
	return ks, nil
}

// GetSrvVSchema returns the SrvVSchema for a cell.
func (f *fakeTopoServer) GetSrvVSchema(ctx context.Context, cell string) (*vschemapb.SrvVSchema, error) {
	vs := &vschemapb.SrvVSchema{
		Keyspaces: map[string]*vschemapb.Keyspace{
			"ks1": {
				Sharded: true,
			},
		},
		RoutingRules: &vschemapb.RoutingRules{
			Rules: []*vschemapb.RoutingRule{
				{
					FromTable: "db1.t1",
					ToTables:  []string{"db1.t1"},
				},
			},
		},
	}
	return vs, nil
}

func (f *fakeTopoServer) WatchSrvKeyspace(ctx context.Context, cell, keyspace string, callback func(*topodatapb.SrvKeyspace, error) bool) {
	ks, err := f.GetSrvKeyspace(ctx, cell, keyspace)
	callback(ks, err)
}

// WatchSrvVSchema starts watching the SrvVSchema object for
// the provided cell.  It will call the callback when
// a new value or an error occurs.
func (f *fakeTopoServer) WatchSrvVSchema(ctx context.Context, cell string, callback func(*vschemapb.SrvVSchema, error) bool) {
	sv, err := f.GetSrvVSchema(ctx, cell)
	callback(sv, err)
}

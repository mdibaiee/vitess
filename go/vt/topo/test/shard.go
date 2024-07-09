/*
Copyright 2019 The Vitess Authors.

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

package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"mdibaiee/vitess/oracle/go/vt/topo"

	topodatapb "mdibaiee/vitess/oracle/go/vt/proto/topodata"
)

// checkShard verifies the Shard operations work correctly
func checkShard(t *testing.T, ctx context.Context, ts *topo.Server) {
	if err := ts.CreateKeyspace(ctx, "test_keyspace", &topodatapb.Keyspace{}); err != nil {
		t.Fatalf("CreateKeyspace: %v", err)
	}

	// Check GetShardNames returns [], nil for existing keyspace with no shards.
	if names, err := ts.GetShardNames(ctx, "test_keyspace"); err != nil || len(names) != 0 {
		t.Errorf("GetShardNames(keyspace with no shards) didn't return [] nil: %v %v", names, err)
	}

	if err := ts.CreateShard(ctx, "test_keyspace", "b0-c0"); err != nil {
		t.Fatalf("CreateShard: %v", err)
	}
	if err := ts.CreateShard(ctx, "test_keyspace", "b0-c0"); !topo.IsErrType(err, topo.NodeExists) {
		t.Errorf("CreateShard called second time, got: %v", err)
	}

	// Delete shard and see if we can re-create it.
	if err := ts.DeleteShard(ctx, "test_keyspace", "b0-c0"); err != nil {
		t.Fatalf("DeleteShard: %v", err)
	}
	if err := ts.DeleteShard(ctx, "test_keyspace", "b0-c0"); !topo.IsErrType(err, topo.NoNode) {
		t.Errorf("DeleteShard(again): %v", err)
	}
	if err := ts.CreateShard(ctx, "test_keyspace", "b0-c0"); err != nil {
		t.Fatalf("CreateShard: %v", err)
	}

	// Test getting an invalid shard returns ErrNoNode.
	if _, err := ts.GetShard(ctx, "test_keyspace", "666"); !topo.IsErrType(err, topo.NoNode) {
		t.Errorf("GetShard(666): %v", err)
	}

	// Test UpdateShardFields works.
	other := &topodatapb.TabletAlias{Cell: "ny", Uid: 82873}
	_, err := ts.UpdateShardFields(ctx, "test_keyspace", "b0-c0", func(si *topo.ShardInfo) error {
		si.PrimaryAlias = other
		return nil
	})
	if err != nil {
		t.Fatalf("UpdateShardFields error: %v", err)
	}

	si, err := ts.GetShard(ctx, "test_keyspace", "b0-c0")
	if err != nil {
		t.Fatalf("GetShard: %v", err)
	}
	if !proto.Equal(si.Shard.PrimaryAlias, other) {
		t.Fatalf("shard.PrimaryAlias = %v, want %v", si.Shard.PrimaryAlias, other)
	}

	// Test FindAllShardsInKeyspace.
	require.NoError(t, err)
	_, err = ts.FindAllShardsInKeyspace(ctx, "test_keyspace", nil)
	require.NoError(t, err)

	// Test GetServingShards.
	require.NoError(t, err)
	_, err = ts.GetServingShards(ctx, "test_keyspace")
	require.NoError(t, err)

	// test GetShardNames
	shardNames, err := ts.GetShardNames(ctx, "test_keyspace")
	if err != nil {
		t.Errorf("GetShardNames: %v", err)
	}
	if len(shardNames) != 1 || shardNames[0] != "b0-c0" {
		t.Errorf(`GetShardNames: want [ "b0-c0" ], got %v`, shardNames)
	}

	if _, err := ts.GetShardNames(ctx, "test_keyspace666"); !topo.IsErrType(err, topo.NoNode) {
		t.Errorf("GetShardNames(666): %v", err)
	}
}

// checkShardWithLock verifies that `TryLockShard` will keep failing with `NodeExists` error if there is
// a lock already taken for given shard. Once we unlock that shard, then subsequent call to `TryLockShard`
// should succeed.
func checkShardWithLock(t *testing.T, ctx context.Context, ts *topo.Server) {
	if err := ts.CreateKeyspace(ctx, "test_keyspace", &topodatapb.Keyspace{}); err != nil {
		t.Fatalf("CreateKeyspace: %v", err)
	}

	unblock := make(chan struct{})
	finished := make(chan struct{})

	// Check GetShardNames returns [], nil for existing keyspace with no shards.
	if names, err := ts.GetShardNames(ctx, "test_keyspace"); err != nil || len(names) != 0 {
		t.Errorf("GetShardNames(keyspace with no shards) didn't return [] nil: %v %v", names, err)
	}

	if err := ts.CreateShard(ctx, "test_keyspace", "b0-c0"); err != nil {
		t.Fatalf("CreateShard: %v", err)
	}

	_, unlock1, err := ts.LockShard(ctx, "test_keyspace", "b0-c0", "lock")
	if err != nil {
		t.Errorf("CreateShard called second time, got: %v", err)
	}

	duration := 10 * time.Second
	waitUntil := time.Now().Add(duration)
	// As soon as we're unblocked, we try to lock the keyspace.
	go func() {
		<-unblock
		var isUnLocked1 = false
		for time.Now().Before(waitUntil) {
			_, unlock2, err := ts.TryLockShard(ctx, "test_keyspace", "b0-c0", "lock")
			// TryLockShard will fail since we already have acquired lock for `test-keyspace`
			if err != nil {
				if !topo.IsErrType(err, topo.NodeExists) {
					require.Fail(t, "expected node exists during tryLockShard", err.Error())
				}
				var finalErr error
				// unlock `test-keyspace` shard. Now the subsequent call to `TryLockShard` will succeed.
				unlock1(&finalErr)
				isUnLocked1 = true
				if finalErr != nil {
					require.Fail(t, "Unlock(test_keyspace) failed", finalErr.Error())
				}
			} else {
				// unlock shard acquired through `TryLockShard`
				unlock2(&err)
				if err != nil {
					require.Fail(t, "Unlock(test_keyspace) failed", err.Error())
				}
				// true value of 'isUnLocked1' signify that we at-least hit 'NodeExits' once.
				if isUnLocked1 {
					close(finished)
				} else {
					require.Fail(t, "Test was expecting to hit `NodeExists` error at-least once")
				}
				break
			}
		}
	}()

	// unblock the go routine
	close(unblock)

	timeout := time.After(duration * 2)
	select {
	case <-finished:
	case <-timeout:
		t.Fatalf("Unlock(test_keyspace) timed out")
	}
}

/*
Copyright 2024 The Vitess Authors.

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

package topo_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mdibaiee/vitess/go/vt/topo"
	"github.com/mdibaiee/vitess/go/vt/topo/memorytopo"

	vschemapb "github.com/mdibaiee/vitess/go/vt/proto/vschema"
)

// lower the lock timeout for testing
const testLockTimeout = 3 * time.Second

// TestTopoLockTimeout tests that the lock times out after the specified duration.
func TestTopoLockTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts := memorytopo.NewServer(ctx, "zone1")
	defer ts.Close()

	err := ts.CreateKeyspaceRoutingRules(ctx, &vschemapb.KeyspaceRoutingRules{})
	require.NoError(t, err)

	currentTopoLockTimeout := topo.LockTimeout
	topo.LockTimeout = testLockTimeout
	defer func() {
		topo.LockTimeout = currentTopoLockTimeout
	}()

	// acquire the lock
	origCtx := ctx
	_, unlock, err := ts.LockRoutingRules(origCtx, "ks1")
	require.NoError(t, err)
	defer unlock(&err)

	// re-acquiring the lock should fail
	_, _, err2 := ts.LockRoutingRules(origCtx, "ks1")
	require.Errorf(t, err2, "deadline exceeded")
}

// TestTopoLockBasic tests basic lock operations.
func TestTopoLockBasic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts := memorytopo.NewServer(ctx, "zone1")
	defer ts.Close()

	err := ts.CreateKeyspaceRoutingRules(ctx, &vschemapb.KeyspaceRoutingRules{})
	require.NoError(t, err)

	origCtx := ctx
	ctx, unlock, err := ts.LockRoutingRules(origCtx, "ks1")
	require.NoError(t, err)

	// locking the same key again, without unlocking, should return an error
	_, _, err2 := ts.LockRoutingRules(ctx, "ks1")
	require.ErrorContains(t, err2, "already held")

	// confirm that the lock can be re-acquired after unlocking
	unlock(&err)
	_, unlock, err = ts.LockRoutingRules(origCtx, "ks1")
	require.NoError(t, err)
	defer unlock(&err)
}

// TestKeyspaceRoutingRulesLock tests that the lock is acquired and released correctly.
func TestKeyspaceRoutingRulesLock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts := memorytopo.NewServer(ctx, "zone1")
	defer ts.Close()

	currentTopoLockTimeout := topo.LockTimeout
	topo.LockTimeout = testLockTimeout
	defer func() {
		topo.LockTimeout = currentTopoLockTimeout
	}()

	err := ts.CreateKeyspaceRoutingRules(ctx, &vschemapb.KeyspaceRoutingRules{})
	require.NoError(t, err)

	_, unlock, err := ts.LockRoutingRules(ctx, "ks1")
	require.NoError(t, err)

	// re-acquiring the lock should fail
	_, _, err = ts.LockRoutingRules(ctx, "ks1")
	require.Error(t, err)

	unlock(&err)

	// re-acquiring the lock should succeed
	_, _, err = ts.LockRoutingRules(ctx, "ks1")
	require.NoError(t, err)
}

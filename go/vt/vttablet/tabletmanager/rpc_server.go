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

package tabletmanager

import (
	"fmt"

	"github.com/mdibaiee/vitess/go/vt/vterrors"

	"context"

	"github.com/mdibaiee/vitess/go/tb"
	"github.com/mdibaiee/vitess/go/vt/callinfo"
	"github.com/mdibaiee/vitess/go/vt/log"
	"github.com/mdibaiee/vitess/go/vt/topo/topoproto"
)

// This file contains the RPC method helpers for the tablet manager.

//
// Utility functions for RPC service
//

// lock is used at the beginning of an RPC call, to acquire the
// action semaphore. It returns ctx.Err() if the context expires.
func (tm *TabletManager) lock(ctx context.Context) error {
	return tm.actionSema.Acquire(ctx, 1)
}

// unlock is the symmetrical action to lock.
func (tm *TabletManager) unlock() {
	tm.actionSema.Release(1)
}

// HandleRPCPanic is part of the RPCTM interface.
func (tm *TabletManager) HandleRPCPanic(ctx context.Context, name string, args, reply any, verbose bool, err *error) {
	// panic handling
	if x := recover(); x != nil {
		log.Errorf("TabletManager.%v(%v) on %v panic: %v\n%s", name, args, topoproto.TabletAliasString(tm.tabletAlias), x, tb.Stack(4))
		*err = fmt.Errorf("caught panic during %v: %v", name, x)
		return
	}

	// quick check for fast path
	if !verbose && *err == nil {
		return
	}

	// we gotta log something, get the source
	from := ""
	ci, ok := callinfo.FromContext(ctx)
	if ok {
		from = ci.Text()
	}

	if *err != nil {
		// error case
		log.Warningf("TabletManager.%v(%v)(on %v from %v) error: %v", name, args, topoproto.TabletAliasString(tm.tabletAlias), from, (*err).Error())
		*err = vterrors.Wrapf(*err, "TabletManager.%v on %v", name, topoproto.TabletAliasString(tm.tabletAlias))
	} else {
		// success case
		log.Infof("TabletManager.%v(%v)(on %v from %v): %#v", name, args, topoproto.TabletAliasString(tm.tabletAlias), from, reply)
	}
}

// RegisterTabletManager is used to delay registration of RPC servers until we have all the objects.
type RegisterTabletManager func(*TabletManager)

// RegisterTabletManagers is a list of functions to call when the delayed registration is triggered.
var RegisterTabletManagers []RegisterTabletManager

// registerTabletManager will register all the instances.
func (tm *TabletManager) registerTabletManager() {
	for _, f := range RegisterTabletManagers {
		f(tm)
	}
}

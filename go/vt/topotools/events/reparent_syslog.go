//go:build !windows

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

package events

import (
	"fmt"
	"log/syslog"

	"github.com/mdibaiee/vitess/go/vt/proto/topodata"

	"github.com/mdibaiee/vitess/go/event/syslogger"
	"github.com/mdibaiee/vitess/go/vt/topo/topoproto"
)

// Syslog writes a Reparent event to syslog.
func (r *Reparent) Syslog() (syslog.Priority, string) {
	var oldAlias *topodata.TabletAlias
	var newAlias *topodata.TabletAlias
	if r.OldPrimary != nil {
		oldAlias = r.OldPrimary.Alias
	}
	if r.NewPrimary != nil {
		newAlias = r.NewPrimary.Alias
	}

	return syslog.LOG_INFO, fmt.Sprintf("%s/%s [reparent %v -> %v] %s (%s)",
		r.ShardInfo.Keyspace(), r.ShardInfo.ShardName(),
		topoproto.TabletAliasString(oldAlias),
		topoproto.TabletAliasString(newAlias),
		r.Status, r.ExternalID)
}

var _ syslogger.Syslogger = (*Reparent)(nil) // compile-time interface check

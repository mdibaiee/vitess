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

package throttler

import (
	"time"

	"github.com/mdibaiee/vitess/go/vt/discovery"
)

// replicationLagRecord stores the tablet health data for a given point in time.
// This data is obtained via the HealthCheck module.
type replicationLagRecord struct {
	// time is the time at which "value" was observed.
	time time.Time

	// LegacyTabletStats holds a copy of the current health data of the tablet.
	discovery.TabletHealth
}

func (r replicationLagRecord) isZero() bool {
	return r.Stats == nil && r.time.IsZero()
}

func (r replicationLagRecord) lag() int64 {
	return int64(r.Stats.ReplicationLagSeconds)
}

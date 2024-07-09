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

package zk2topo

import (
	"fmt"

	"github.com/mdibaiee/vitess/go/vt/topo"
)

// ZKVersion is zookeeper's idea of a version.
// It implements topo.Version.
// We use the native zookeeper.Stat.Version type, int32.
type ZKVersion int32

// String is part of the topo.Version interface.
func (v ZKVersion) String() string {
	return fmt.Sprintf("%v", int32(v))
}

var _ topo.Version = (ZKVersion)(0) // compile-time interface check

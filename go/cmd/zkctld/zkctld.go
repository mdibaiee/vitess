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

// zkctld is a daemon that starts or initializes ZooKeeper with Vitess-specific
// configuration. It will stay running as long as the underlying ZooKeeper
// server, and will pass along SIGTERM.
package main

import (
	"mdibaiee/vitess/oracle/go/cmd/zkctld/cli"
	"mdibaiee/vitess/oracle/go/exit"
	"mdibaiee/vitess/oracle/go/vt/log"
)

func main() {
	defer exit.Recover()
	if err := cli.Main.Execute(); err != nil {
		log.Error(err)
		exit.Return(1)
	}
}

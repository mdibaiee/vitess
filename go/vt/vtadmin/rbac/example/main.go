/*
Copyright 2021 The Vitess Authors.

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

package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"mdibaiee/vitess/go/vt/vtadmin/rbac"
)

func main() {
	flag.Parse()

	cfg, err := rbac.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", data)
}

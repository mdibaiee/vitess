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

package engine

import (
	"context"

	"github.com/mdibaiee/vitess/go/vt/vterrors"
)

type failDBDDL struct{}

// CreateDatabase implements the DropCreateDB interface
func (failDBDDL) CreateDatabase(context.Context, string) error {
	return vterrors.VT12001("create database by failDBDDL")
}

// DropDatabase implements the DropCreateDB interface
func (failDBDDL) DropDatabase(context.Context, string) error {
	return vterrors.VT12001("drop database by failDBDDL")
}

type noOp struct{}

// CreateDatabase implements the DropCreateDB interface
func (noOp) CreateDatabase(context.Context, string) error {
	return nil
}

// DropDatabase implements the DropCreateDB interface
func (noOp) DropDatabase(context.Context, string) error {
	return nil
}

const (
	faildbDDL          = "fail"
	noOpdbDDL          = "noop"
	defaultDBDDLPlugin = faildbDDL
)

func init() {
	DBDDLRegister(faildbDDL, failDBDDL{})
	DBDDLRegister(noOpdbDDL, noOp{})
}

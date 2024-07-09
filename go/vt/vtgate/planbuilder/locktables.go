/*
Copyright 2020 The Vitess Authors.

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

package planbuilder

import (
	"mdibaiee/vitess/oracle/go/sqltypes"
	"mdibaiee/vitess/oracle/go/vt/log"
	querypb "mdibaiee/vitess/oracle/go/vt/proto/query"
	"mdibaiee/vitess/oracle/go/vt/sqlparser"
	"mdibaiee/vitess/oracle/go/vt/vtgate/engine"
	"mdibaiee/vitess/oracle/go/vt/vtgate/planbuilder/plancontext"
)

// buildLockPlan plans lock tables statement.
func buildLockPlan(stmt sqlparser.Statement, _ *sqlparser.ReservedVars, _ plancontext.VSchema) (*planResult, error) {
	log.Warningf("Lock Tables statement is ignored: %v", stmt)
	return newPlanResult(engine.NewRowsPrimitive(make([][]sqltypes.Value, 0), make([]*querypb.Field, 0))), nil
}

// buildUnlockPlan plans lock tables statement.
func buildUnlockPlan(stmt sqlparser.Statement, _ *sqlparser.ReservedVars, _ plancontext.VSchema) (*planResult, error) {
	return newPlanResult(&engine.Unlock{}), nil
}

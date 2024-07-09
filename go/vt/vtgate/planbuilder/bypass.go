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

package planbuilder

import (
	"github.com/mdibaiee/vitess/go/vt/key"
	"github.com/mdibaiee/vitess/go/vt/sqlparser"
	"github.com/mdibaiee/vitess/go/vt/vterrors"
	"github.com/mdibaiee/vitess/go/vt/vtgate/engine"
	"github.com/mdibaiee/vitess/go/vt/vtgate/planbuilder/plancontext"
	"github.com/mdibaiee/vitess/go/vt/vtgate/vindexes"
)

func buildPlanForBypass(stmt sqlparser.Statement, _ *sqlparser.ReservedVars, vschema plancontext.VSchema) (*planResult, error) {
	keyspace, err := vschema.DefaultKeyspace()
	if err != nil {
		return nil, err
	}
	switch dest := vschema.Destination().(type) {
	case key.DestinationExactKeyRange:
		if _, ok := stmt.(*sqlparser.Insert); ok {
			return nil, vterrors.VT03023(vschema.TargetString())
		}
	case key.DestinationShard:
		if !vschema.IsShardRoutingEnabled() {
			break
		}
		shard := string(dest)
		targetKeyspace, err := GetShardRoute(vschema, keyspace.Name, shard)
		if err != nil {
			return nil, err
		}
		if targetKeyspace != nil {
			keyspace = targetKeyspace
		}
	}

	hints := &queryHints{}
	if comments, ok := stmt.(sqlparser.Commented); ok {
		if qh := getHints(comments.GetParsedComments()); qh != nil {
			hints = qh
		}
	}

	send := &engine.Send{
		Keyspace:             keyspace,
		TargetDestination:    vschema.Destination(),
		Query:                sqlparser.String(stmt),
		IsDML:                sqlparser.IsDMLStatement(stmt),
		SingleShardOnly:      false,
		MultishardAutocommit: hints.multiShardAutocommit,
		QueryTimeout:         hints.queryTimeout,
	}
	return newPlanResult(send), nil
}

func GetShardRoute(vschema plancontext.VSchema, keyspace, shard string) (*vindexes.Keyspace, error) {
	targetKeyspaceName, err := vschema.FindRoutedShard(keyspace, shard)
	if err != nil {
		return nil, err
	}
	return vschema.FindKeyspace(targetKeyspaceName)
}

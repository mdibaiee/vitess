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
// Code generated by Sizegen. DO NOT EDIT.

package planbuilder

import hack "mdibaiee/vitess/oracle/go/hack"

type cachedObject interface {
	CachedSize(alloc bool) int64
}

func (cached *Permission) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(24)
	}
	// field TableName string
	size += hack.RuntimeAllocSize(int64(len(cached.TableName)))
	return size
}
func (cached *Plan) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(128)
	}
	// field Table *mdibaiee/vitess/oracle/go/vt/vttablet/tabletserver/schema.Table
	size += cached.Table.CachedSize(true)
	// field AllTables []*mdibaiee/vitess/oracle/go/vt/vttablet/tabletserver/schema.Table
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.AllTables)) * int64(8))
		for _, elem := range cached.AllTables {
			size += elem.CachedSize(true)
		}
	}
	// field Permissions []mdibaiee/vitess/oracle/go/vt/vttablet/tabletserver/planbuilder.Permission
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.Permissions)) * int64(24))
		for _, elem := range cached.Permissions {
			size += elem.CachedSize(false)
		}
	}
	// field FullQuery *mdibaiee/vitess/oracle/go/vt/sqlparser.ParsedQuery
	size += cached.FullQuery.CachedSize(true)
	// field NextCount mdibaiee/vitess/oracle/go/vt/vtgate/evalengine.Expr
	if cc, ok := cached.NextCount.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	// field WhereClause *mdibaiee/vitess/oracle/go/vt/sqlparser.ParsedQuery
	size += cached.WhereClause.CachedSize(true)
	// field FullStmt mdibaiee/vitess/oracle/go/vt/sqlparser.Statement
	if cc, ok := cached.FullStmt.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	return size
}

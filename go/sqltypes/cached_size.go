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

package sqltypes

import hack "github.com/mdibaiee/vitess/go/hack"

func (cached *Result) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(112)
	}
	// field Fields []*github.com/mdibaiee/vitess/go/vt/proto/query.Field
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.Fields)) * int64(8))
		for _, elem := range cached.Fields {
			size += elem.CachedSize(true)
		}
	}
	// field Rows [][]github.com/mdibaiee/vitess/go/sqltypes.Value
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.Rows)) * int64(24))
		for _, elem := range cached.Rows {
			{
				size += hack.RuntimeAllocSize(int64(cap(elem)) * int64(32))
				for _, elem := range elem {
					size += elem.CachedSize(false)
				}
			}
		}
	}
	// field SessionStateChanges string
	size += hack.RuntimeAllocSize(int64(len(cached.SessionStateChanges)))
	// field Info string
	size += hack.RuntimeAllocSize(int64(len(cached.Info)))
	return size
}
func (cached *Value) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(32)
	}
	// field val []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.val)))
	}
	return size
}

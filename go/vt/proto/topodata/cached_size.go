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

package topodata

import hack "github.com/mdibaiee/vitess/go/hack"

func (cached *KeyRange) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(96)
	}
	// field unknownFields []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.unknownFields)))
	}
	// field Start []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.Start)))
	}
	// field End []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.End)))
	}
	return size
}
func (cached *ThrottledAppRule) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(80)
	}
	// field unknownFields []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.unknownFields)))
	}
	// field Name string
	size += hack.RuntimeAllocSize(int64(len(cached.Name)))
	// field ExpiresAt *github.com/mdibaiee/vitess/go/vt/proto/vttime.Time
	size += cached.ExpiresAt.CachedSize(true)
	return size
}

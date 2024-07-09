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

package integration

import (
	"math"
	"reflect"
	"unsafe"

	hack "mdibaiee/vitess/oracle/go/hack"
)

type cachedObject interface {
	CachedSize(alloc bool) int64
}

func (cached *A) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(16)
	}
	return size
}
func (cached *Bimpl) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(8)
	}
	return size
}
func (cached *C) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(16)
	}
	// field field1 mdibaiee/vitess/oracle/go/tools/sizegen/integration.B
	if cc, ok := cached.field1.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	return size
}
func (cached *D) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(8)
	}
	// field field1 *mdibaiee/vitess/oracle/go/tools/sizegen/integration.Bimpl
	if cached.field1 != nil {
		size += hack.RuntimeAllocSize(int64(8))
	}
	return size
}

//go:nocheckptr
func (cached *Map1) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(8)
	}
	// field field1 map[uint8]uint8
	if cached.field1 != nil {
		size += int64(48)
		hmap := reflect.ValueOf(cached.field1)
		numBuckets := int(math.Pow(2, float64((*(*uint8)(unsafe.Pointer(hmap.Pointer() + uintptr(9)))))))
		numOldBuckets := (*(*uint16)(unsafe.Pointer(hmap.Pointer() + uintptr(10))))
		size += hack.RuntimeAllocSize(int64(numOldBuckets * 32))
		if len(cached.field1) > 0 || numBuckets > 1 {
			size += hack.RuntimeAllocSize(int64(numBuckets * 32))
		}
	}
	return size
}

//go:nocheckptr
func (cached *Map2) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(8)
	}
	// field field1 map[uint64]mdibaiee/vitess/oracle/go/tools/sizegen/integration.A
	if cached.field1 != nil {
		size += int64(48)
		hmap := reflect.ValueOf(cached.field1)
		numBuckets := int(math.Pow(2, float64((*(*uint8)(unsafe.Pointer(hmap.Pointer() + uintptr(9)))))))
		numOldBuckets := (*(*uint16)(unsafe.Pointer(hmap.Pointer() + uintptr(10))))
		size += hack.RuntimeAllocSize(int64(numOldBuckets * 208))
		if len(cached.field1) > 0 || numBuckets > 1 {
			size += hack.RuntimeAllocSize(int64(numBuckets * 208))
		}
	}
	return size
}

//go:nocheckptr
func (cached *Map3) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(8)
	}
	// field field1 map[uint64]mdibaiee/vitess/oracle/go/tools/sizegen/integration.B
	if cached.field1 != nil {
		size += int64(48)
		hmap := reflect.ValueOf(cached.field1)
		numBuckets := int(math.Pow(2, float64((*(*uint8)(unsafe.Pointer(hmap.Pointer() + uintptr(9)))))))
		numOldBuckets := (*(*uint16)(unsafe.Pointer(hmap.Pointer() + uintptr(10))))
		size += hack.RuntimeAllocSize(int64(numOldBuckets * 208))
		if len(cached.field1) > 0 || numBuckets > 1 {
			size += hack.RuntimeAllocSize(int64(numBuckets * 208))
		}
		for _, v := range cached.field1 {
			if cc, ok := v.(cachedObject); ok {
				size += cc.CachedSize(true)
			}
		}
	}
	return size
}
func (cached *Padded) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(24)
	}
	return size
}
func (cached *Slice1) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(24)
	}
	// field field1 []mdibaiee/vitess/oracle/go/tools/sizegen/integration.A
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.field1)) * int64(16))
	}
	return size
}
func (cached *Slice2) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(24)
	}
	// field field1 []mdibaiee/vitess/oracle/go/tools/sizegen/integration.B
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.field1)) * int64(16))
		for _, elem := range cached.field1 {
			if cc, ok := elem.(cachedObject); ok {
				size += cc.CachedSize(true)
			}
		}
	}
	return size
}
func (cached *Slice3) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(24)
	}
	// field field1 []*mdibaiee/vitess/oracle/go/tools/sizegen/integration.Bimpl
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.field1)) * int64(8))
		for _, elem := range cached.field1 {
			if elem != nil {
				size += hack.RuntimeAllocSize(int64(8))
			}
		}
	}
	return size
}
func (cached *String1) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(24)
	}
	// field field1 string
	size += hack.RuntimeAllocSize(int64(len(cached.field1)))
	return size
}

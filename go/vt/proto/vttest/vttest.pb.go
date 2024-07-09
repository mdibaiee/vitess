//
//Copyright 2019 The Vitess Authors.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

// This file contains the vttest topology configuration structures.
//
// The protobuf message "VTTestTopology" specifies the Vitess and
// database configuration of the "vttest" test component which can be
// used in end-to-end tests to test an application against an actual
// Vitess and MySQL server.
//
// To start a "vttest" instance, pass the "VTTestTopology" message,
// encoded in the protobuf compact text format, to
// py/vttest/run_local_database.py which in turn will send it to the
// Vitess test binary called "vtcombo".
//
// To encode a "VTTestTopology" message in the protobuf compact text
// format, create the protobuf in your test's native language first
// and then use the protobuf library to encode it as text.
// For an example in Python, see: test/vttest_sample_test.py
// In go, see: go/vt/vttest/local_cluster_test.go
//
// Sample encoded proto configurations would be as follow. Note there are
// multiple encoding options, see the proto documentation for more info
// (first and last quote not included in the encoding):
// - single keyspace named test_keyspace with one shard '0':
//   'keyspaces:<name:"test_keyspace" shards:<name:"0" > > '
// - two keyspaces, one with two shards, the other one with a redirect:
//   'keyspaces { name: "test_keyspace" shards { name: "-80" } shards { name: "80-" } } keyspaces { name: "redirect" served_from: "test_keyspace" }'

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.3
// source: vttest.proto

package vttest

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	vschema "mdibaiee/vitess/oracle/go/vt/proto/vschema"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Shard describes a single shard in a keyspace.
type Shard struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// name has to be unique in a keyspace. For unsharded keyspaces, it
	// should be '0'. For sharded keyspace, it should be derived from
	// the keyrange, like '-80' or '40-80'.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// db_name_override is the mysql db name for this shard. Has to be
	// globally unique. If not specified, we will by default use
	// 'vt_<keyspace>_<shard>'.
	DbNameOverride string `protobuf:"bytes,2,opt,name=db_name_override,json=dbNameOverride,proto3" json:"db_name_override,omitempty"`
}

func (x *Shard) Reset() {
	*x = Shard{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vttest_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Shard) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Shard) ProtoMessage() {}

func (x *Shard) ProtoReflect() protoreflect.Message {
	mi := &file_vttest_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Shard.ProtoReflect.Descriptor instead.
func (*Shard) Descriptor() ([]byte, []int) {
	return file_vttest_proto_rawDescGZIP(), []int{0}
}

func (x *Shard) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Shard) GetDbNameOverride() string {
	if x != nil {
		return x.DbNameOverride
	}
	return ""
}

// Keyspace describes a single keyspace.
type Keyspace struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// name has to be unique in a VTTestTopology.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// shards inside this keyspace. Ignored if redirect is set.
	Shards []*Shard `protobuf:"bytes,2,rep,name=shards,proto3" json:"shards,omitempty"`
	// number of replica tablets to instantiate. This includes the primary tablet.
	ReplicaCount int32 `protobuf:"varint,6,opt,name=replica_count,json=replicaCount,proto3" json:"replica_count,omitempty"`
	// number of rdonly tablets to instantiate.
	RdonlyCount int32 `protobuf:"varint,7,opt,name=rdonly_count,json=rdonlyCount,proto3" json:"rdonly_count,omitempty"`
}

func (x *Keyspace) Reset() {
	*x = Keyspace{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vttest_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Keyspace) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Keyspace) ProtoMessage() {}

func (x *Keyspace) ProtoReflect() protoreflect.Message {
	mi := &file_vttest_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Keyspace.ProtoReflect.Descriptor instead.
func (*Keyspace) Descriptor() ([]byte, []int) {
	return file_vttest_proto_rawDescGZIP(), []int{1}
}

func (x *Keyspace) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Keyspace) GetShards() []*Shard {
	if x != nil {
		return x.Shards
	}
	return nil
}

func (x *Keyspace) GetReplicaCount() int32 {
	if x != nil {
		return x.ReplicaCount
	}
	return 0
}

func (x *Keyspace) GetRdonlyCount() int32 {
	if x != nil {
		return x.RdonlyCount
	}
	return 0
}

// VTTestTopology describes the keyspaces in the topology.
type VTTestTopology struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// all keyspaces in the topology.
	Keyspaces []*Keyspace `protobuf:"bytes,1,rep,name=keyspaces,proto3" json:"keyspaces,omitempty"`
	// list of cells the keyspaces reside in. Vtgate is started in only the first cell.
	Cells []string `protobuf:"bytes,2,rep,name=cells,proto3" json:"cells,omitempty"`
	// routing rules for the topology.
	RoutingRules *vschema.RoutingRules `protobuf:"bytes,3,opt,name=routing_rules,json=routingRules,proto3" json:"routing_rules,omitempty"`
}

func (x *VTTestTopology) Reset() {
	*x = VTTestTopology{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vttest_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VTTestTopology) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VTTestTopology) ProtoMessage() {}

func (x *VTTestTopology) ProtoReflect() protoreflect.Message {
	mi := &file_vttest_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VTTestTopology.ProtoReflect.Descriptor instead.
func (*VTTestTopology) Descriptor() ([]byte, []int) {
	return file_vttest_proto_rawDescGZIP(), []int{2}
}

func (x *VTTestTopology) GetKeyspaces() []*Keyspace {
	if x != nil {
		return x.Keyspaces
	}
	return nil
}

func (x *VTTestTopology) GetCells() []string {
	if x != nil {
		return x.Cells
	}
	return nil
}

func (x *VTTestTopology) GetRoutingRules() *vschema.RoutingRules {
	if x != nil {
		return x.RoutingRules
	}
	return nil
}

var File_vttest_proto protoreflect.FileDescriptor

var file_vttest_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x76, 0x74, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x76, 0x74, 0x74, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x45, 0x0a, 0x05, 0x53, 0x68, 0x61, 0x72, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x64, 0x62, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x6f, 0x76,
	0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x64, 0x62,
	0x4e, 0x61, 0x6d, 0x65, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x22, 0x9f, 0x01, 0x0a,
	0x08, 0x4b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a,
	0x06, 0x73, 0x68, 0x61, 0x72, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x76, 0x74, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x52, 0x06, 0x73, 0x68,
	0x61, 0x72, 0x64, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x5f,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x72, 0x65, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x64, 0x6f,
	0x6e, 0x6c, 0x79, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0b, 0x72, 0x64, 0x6f, 0x6e, 0x6c, 0x79, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x4a, 0x04, 0x08, 0x03,
	0x10, 0x04, 0x4a, 0x04, 0x08, 0x04, 0x10, 0x05, 0x4a, 0x04, 0x08, 0x05, 0x10, 0x06, 0x22, 0x92,
	0x01, 0x0a, 0x0e, 0x56, 0x54, 0x54, 0x65, 0x73, 0x74, 0x54, 0x6f, 0x70, 0x6f, 0x6c, 0x6f, 0x67,
	0x79, 0x12, 0x2e, 0x0a, 0x09, 0x6b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x76, 0x74, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x4b, 0x65,
	0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x09, 0x6b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x65, 0x6c, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x05, 0x63, 0x65, 0x6c, 0x6c, 0x73, 0x12, 0x3a, 0x0a, 0x0d, 0x72, 0x6f, 0x75, 0x74, 0x69,
	0x6e, 0x67, 0x5f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x0c, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x75,
	0x6c, 0x65, 0x73, 0x42, 0x25, 0x5a, 0x23, 0x76, 0x69, 0x74, 0x65, 0x73, 0x73, 0x2e, 0x69, 0x6f,
	0x2f, 0x76, 0x69, 0x74, 0x65, 0x73, 0x73, 0x2f, 0x67, 0x6f, 0x2f, 0x76, 0x74, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x74, 0x74, 0x65, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_vttest_proto_rawDescOnce sync.Once
	file_vttest_proto_rawDescData = file_vttest_proto_rawDesc
)

func file_vttest_proto_rawDescGZIP() []byte {
	file_vttest_proto_rawDescOnce.Do(func() {
		file_vttest_proto_rawDescData = protoimpl.X.CompressGZIP(file_vttest_proto_rawDescData)
	})
	return file_vttest_proto_rawDescData
}

var file_vttest_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_vttest_proto_goTypes = []any{
	(*Shard)(nil),                // 0: vttest.Shard
	(*Keyspace)(nil),             // 1: vttest.Keyspace
	(*VTTestTopology)(nil),       // 2: vttest.VTTestTopology
	(*vschema.RoutingRules)(nil), // 3: vschema.RoutingRules
}
var file_vttest_proto_depIdxs = []int32{
	0, // 0: vttest.Keyspace.shards:type_name -> vttest.Shard
	1, // 1: vttest.VTTestTopology.keyspaces:type_name -> vttest.Keyspace
	3, // 2: vttest.VTTestTopology.routing_rules:type_name -> vschema.RoutingRules
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_vttest_proto_init() }
func file_vttest_proto_init() {
	if File_vttest_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_vttest_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Shard); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vttest_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Keyspace); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vttest_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*VTTestTopology); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_vttest_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_vttest_proto_goTypes,
		DependencyIndexes: file_vttest_proto_depIdxs,
		MessageInfos:      file_vttest_proto_msgTypes,
	}.Build()
	File_vttest_proto = out.File
	file_vttest_proto_rawDesc = nil
	file_vttest_proto_goTypes = nil
	file_vttest_proto_depIdxs = nil
}

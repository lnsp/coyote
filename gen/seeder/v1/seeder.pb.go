// Copyright 2020 Lennart Espe
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: seeder/v1/seeder.proto

package seederv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type HasRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *HasRequest) Reset() {
	*x = HasRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_seeder_v1_seeder_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HasRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HasRequest) ProtoMessage() {}

func (x *HasRequest) ProtoReflect() protoreflect.Message {
	mi := &file_seeder_v1_seeder_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HasRequest.ProtoReflect.Descriptor instead.
func (*HasRequest) Descriptor() ([]byte, []int) {
	return file_seeder_v1_seeder_proto_rawDescGZIP(), []int{0}
}

func (x *HasRequest) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

type HasResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chunks []uint64 `protobuf:"varint,1,rep,packed,name=chunks,proto3" json:"chunks,omitempty"`
}

func (x *HasResponse) Reset() {
	*x = HasResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_seeder_v1_seeder_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HasResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HasResponse) ProtoMessage() {}

func (x *HasResponse) ProtoReflect() protoreflect.Message {
	mi := &file_seeder_v1_seeder_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HasResponse.ProtoReflect.Descriptor instead.
func (*HasResponse) Descriptor() ([]byte, []int) {
	return file_seeder_v1_seeder_proto_rawDescGZIP(), []int{1}
}

func (x *HasResponse) GetChunks() []uint64 {
	if x != nil {
		return x.Chunks
	}
	return nil
}

type FetchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash  []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Chunk int64  `protobuf:"varint,2,opt,name=chunk,proto3" json:"chunk,omitempty"`
}

func (x *FetchRequest) Reset() {
	*x = FetchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_seeder_v1_seeder_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchRequest) ProtoMessage() {}

func (x *FetchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_seeder_v1_seeder_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchRequest.ProtoReflect.Descriptor instead.
func (*FetchRequest) Descriptor() ([]byte, []int) {
	return file_seeder_v1_seeder_proto_rawDescGZIP(), []int{2}
}

func (x *FetchRequest) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *FetchRequest) GetChunk() int64 {
	if x != nil {
		return x.Chunk
	}
	return 0
}

type FetchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *FetchResponse) Reset() {
	*x = FetchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_seeder_v1_seeder_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchResponse) ProtoMessage() {}

func (x *FetchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_seeder_v1_seeder_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchResponse.ProtoReflect.Descriptor instead.
func (*FetchResponse) Descriptor() ([]byte, []int) {
	return file_seeder_v1_seeder_proto_rawDescGZIP(), []int{3}
}

func (x *FetchResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_seeder_v1_seeder_proto protoreflect.FileDescriptor

var file_seeder_v1_seeder_proto_rawDesc = []byte{
	0x0a, 0x16, 0x73, 0x65, 0x65, 0x64, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x65, 0x64,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x65, 0x65, 0x64, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x22, 0x20, 0x0a, 0x0a, 0x48, 0x61, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x68, 0x61, 0x73, 0x68, 0x22, 0x25, 0x0a, 0x0b, 0x48, 0x61, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x04, 0x52, 0x06, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x22, 0x38, 0x0a, 0x0c,
	0x46, 0x65, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x22, 0x23, 0x0a, 0x0d, 0x46, 0x65, 0x74, 0x63, 0x68, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x85, 0x01, 0x0a, 0x0d,
	0x53, 0x65, 0x65, 0x64, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x36, 0x0a,
	0x03, 0x48, 0x61, 0x73, 0x12, 0x15, 0x2e, 0x73, 0x65, 0x65, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x48, 0x61, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x73, 0x65,
	0x65, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x61, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x05, 0x46, 0x65, 0x74, 0x63, 0x68, 0x12, 0x17,
	0x2e, 0x73, 0x65, 0x65, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x65, 0x65, 0x64, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6c, 0x6e, 0x73, 0x70, 0x2f, 0x63, 0x6f, 0x79, 0x6f, 0x74, 0x65, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x73, 0x65, 0x65, 0x64, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x65, 0x65, 0x64,
	0x65, 0x72, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_seeder_v1_seeder_proto_rawDescOnce sync.Once
	file_seeder_v1_seeder_proto_rawDescData = file_seeder_v1_seeder_proto_rawDesc
)

func file_seeder_v1_seeder_proto_rawDescGZIP() []byte {
	file_seeder_v1_seeder_proto_rawDescOnce.Do(func() {
		file_seeder_v1_seeder_proto_rawDescData = protoimpl.X.CompressGZIP(file_seeder_v1_seeder_proto_rawDescData)
	})
	return file_seeder_v1_seeder_proto_rawDescData
}

var file_seeder_v1_seeder_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_seeder_v1_seeder_proto_goTypes = []interface{}{
	(*HasRequest)(nil),    // 0: seeder.v1.HasRequest
	(*HasResponse)(nil),   // 1: seeder.v1.HasResponse
	(*FetchRequest)(nil),  // 2: seeder.v1.FetchRequest
	(*FetchResponse)(nil), // 3: seeder.v1.FetchResponse
}
var file_seeder_v1_seeder_proto_depIdxs = []int32{
	0, // 0: seeder.v1.SeederService.Has:input_type -> seeder.v1.HasRequest
	2, // 1: seeder.v1.SeederService.Fetch:input_type -> seeder.v1.FetchRequest
	1, // 2: seeder.v1.SeederService.Has:output_type -> seeder.v1.HasResponse
	3, // 3: seeder.v1.SeederService.Fetch:output_type -> seeder.v1.FetchResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_seeder_v1_seeder_proto_init() }
func file_seeder_v1_seeder_proto_init() {
	if File_seeder_v1_seeder_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_seeder_v1_seeder_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HasRequest); i {
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
		file_seeder_v1_seeder_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HasResponse); i {
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
		file_seeder_v1_seeder_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FetchRequest); i {
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
		file_seeder_v1_seeder_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FetchResponse); i {
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
			RawDescriptor: file_seeder_v1_seeder_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_seeder_v1_seeder_proto_goTypes,
		DependencyIndexes: file_seeder_v1_seeder_proto_depIdxs,
		MessageInfos:      file_seeder_v1_seeder_proto_msgTypes,
	}.Build()
	File_seeder_v1_seeder_proto = out.File
	file_seeder_v1_seeder_proto_rawDesc = nil
	file_seeder_v1_seeder_proto_goTypes = nil
	file_seeder_v1_seeder_proto_depIdxs = nil
}

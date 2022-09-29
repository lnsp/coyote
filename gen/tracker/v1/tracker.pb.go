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
// source: tracker/v1/tracker.proto

package trackerv1

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

type Tracker struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Taverns     []string `protobuf:"bytes,1,rep,name=taverns,proto3" json:"taverns,omitempty"`
	Hash        []byte   `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	Name        string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Size        int64    `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	ChunkSize   int64    `protobuf:"varint,5,opt,name=chunk_size,json=chunkSize,proto3" json:"chunk_size,omitempty"`
	ChunkHashes [][]byte `protobuf:"bytes,6,rep,name=chunk_hashes,json=chunkHashes,proto3" json:"chunk_hashes,omitempty"`
}

func (x *Tracker) Reset() {
	*x = Tracker{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_v1_tracker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tracker) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tracker) ProtoMessage() {}

func (x *Tracker) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_v1_tracker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tracker.ProtoReflect.Descriptor instead.
func (*Tracker) Descriptor() ([]byte, []int) {
	return file_tracker_v1_tracker_proto_rawDescGZIP(), []int{0}
}

func (x *Tracker) GetTaverns() []string {
	if x != nil {
		return x.Taverns
	}
	return nil
}

func (x *Tracker) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Tracker) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Tracker) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Tracker) GetChunkSize() int64 {
	if x != nil {
		return x.ChunkSize
	}
	return 0
}

func (x *Tracker) GetChunkHashes() [][]byte {
	if x != nil {
		return x.ChunkHashes
	}
	return nil
}

var File_tracker_v1_tracker_proto protoreflect.FileDescriptor

var file_tracker_v1_tracker_proto_rawDesc = []byte{
	0x0a, 0x18, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x72, 0x61,
	0x63, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x74, 0x72, 0x61, 0x63,
	0x6b, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x22, 0xa1, 0x01, 0x0a, 0x07, 0x54, 0x72, 0x61, 0x63, 0x6b,
	0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x61, 0x76, 0x65, 0x72, 0x6e, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x07, 0x74, 0x61, 0x76, 0x65, 0x72, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x75, 0x6e,
	0x6b, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x75, 0x6e, 0x6b,
	0x5f, 0x68, 0x61, 0x73, 0x68, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0b, 0x63,
	0x68, 0x75, 0x6e, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x65, 0x73, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x6e, 0x73, 0x70, 0x2f, 0x66, 0x74,
	0x70, 0x32, 0x70, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2f,
	0x76, 0x31, 0x3b, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tracker_v1_tracker_proto_rawDescOnce sync.Once
	file_tracker_v1_tracker_proto_rawDescData = file_tracker_v1_tracker_proto_rawDesc
)

func file_tracker_v1_tracker_proto_rawDescGZIP() []byte {
	file_tracker_v1_tracker_proto_rawDescOnce.Do(func() {
		file_tracker_v1_tracker_proto_rawDescData = protoimpl.X.CompressGZIP(file_tracker_v1_tracker_proto_rawDescData)
	})
	return file_tracker_v1_tracker_proto_rawDescData
}

var file_tracker_v1_tracker_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_tracker_v1_tracker_proto_goTypes = []interface{}{
	(*Tracker)(nil), // 0: tracker.v1.Tracker
}
var file_tracker_v1_tracker_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_tracker_v1_tracker_proto_init() }
func file_tracker_v1_tracker_proto_init() {
	if File_tracker_v1_tracker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tracker_v1_tracker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tracker); i {
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
			RawDescriptor: file_tracker_v1_tracker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tracker_v1_tracker_proto_goTypes,
		DependencyIndexes: file_tracker_v1_tracker_proto_depIdxs,
		MessageInfos:      file_tracker_v1_tracker_proto_msgTypes,
	}.Build()
	File_tracker_v1_tracker_proto = out.File
	file_tracker_v1_tracker_proto_rawDesc = nil
	file_tracker_v1_tracker_proto_goTypes = nil
	file_tracker_v1_tracker_proto_depIdxs = nil
}

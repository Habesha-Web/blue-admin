// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.21.12
// source: bluerpc/bluerpc.proto

package bluerpc

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

type BlueSalt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SaltA string `protobuf:"bytes,1,opt,name=salt_a,json=saltA,proto3" json:"salt_a,omitempty"`
	SaltB string `protobuf:"bytes,2,opt,name=salt_b,json=saltB,proto3" json:"salt_b,omitempty"`
}

func (x *BlueSalt) Reset() {
	*x = BlueSalt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bluerpc_bluerpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlueSalt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlueSalt) ProtoMessage() {}

func (x *BlueSalt) ProtoReflect() protoreflect.Message {
	mi := &file_bluerpc_bluerpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlueSalt.ProtoReflect.Descriptor instead.
func (*BlueSalt) Descriptor() ([]byte, []int) {
	return file_bluerpc_bluerpc_proto_rawDescGZIP(), []int{0}
}

func (x *BlueSalt) GetSaltA() string {
	if x != nil {
		return x.SaltA
	}
	return ""
}

func (x *BlueSalt) GetSaltB() string {
	if x != nil {
		return x.SaltB
	}
	return ""
}

type BlueAppID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppId string `protobuf:"bytes,1,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"`
}

func (x *BlueAppID) Reset() {
	*x = BlueAppID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bluerpc_bluerpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlueAppID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlueAppID) ProtoMessage() {}

func (x *BlueAppID) ProtoReflect() protoreflect.Message {
	mi := &file_bluerpc_bluerpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlueAppID.ProtoReflect.Descriptor instead.
func (*BlueAppID) Descriptor() ([]byte, []int) {
	return file_bluerpc_bluerpc_proto_rawDescGZIP(), []int{1}
}

func (x *BlueAppID) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

var File_bluerpc_bluerpc_proto protoreflect.FileDescriptor

var file_bluerpc_bluerpc_proto_rawDesc = []byte{
	0x0a, 0x15, 0x62, 0x6c, 0x75, 0x65, 0x72, 0x70, 0x63, 0x2f, 0x62, 0x6c, 0x75, 0x65, 0x72, 0x70,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x38, 0x0a, 0x08, 0x42, 0x6c, 0x75, 0x65, 0x53,
	0x61, 0x6c, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x61, 0x6c, 0x74, 0x5f, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x61, 0x6c, 0x74, 0x41, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x61,
	0x6c, 0x74, 0x5f, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x61, 0x6c, 0x74,
	0x42, 0x22, 0x22, 0x0a, 0x09, 0x42, 0x6c, 0x75, 0x65, 0x41, 0x70, 0x70, 0x49, 0x44, 0x12, 0x15,
	0x0a, 0x06, 0x61, 0x70, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x61, 0x70, 0x70, 0x49, 0x64, 0x32, 0x31, 0x0a, 0x0b, 0x42, 0x6c, 0x75, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x22, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x53, 0x61, 0x6c, 0x74, 0x12,
	0x0a, 0x2e, 0x42, 0x6c, 0x75, 0x65, 0x41, 0x70, 0x70, 0x49, 0x44, 0x1a, 0x09, 0x2e, 0x42, 0x6c,
	0x75, 0x65, 0x53, 0x61, 0x6c, 0x74, 0x22, 0x00, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x2f, 0x62, 0x6c,
	0x75, 0x65, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_bluerpc_bluerpc_proto_rawDescOnce sync.Once
	file_bluerpc_bluerpc_proto_rawDescData = file_bluerpc_bluerpc_proto_rawDesc
)

func file_bluerpc_bluerpc_proto_rawDescGZIP() []byte {
	file_bluerpc_bluerpc_proto_rawDescOnce.Do(func() {
		file_bluerpc_bluerpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_bluerpc_bluerpc_proto_rawDescData)
	})
	return file_bluerpc_bluerpc_proto_rawDescData
}

var file_bluerpc_bluerpc_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_bluerpc_bluerpc_proto_goTypes = []interface{}{
	(*BlueSalt)(nil),  // 0: BlueSalt
	(*BlueAppID)(nil), // 1: BlueAppID
}
var file_bluerpc_bluerpc_proto_depIdxs = []int32{
	1, // 0: BlueService.GetSalt:input_type -> BlueAppID
	0, // 1: BlueService.GetSalt:output_type -> BlueSalt
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_bluerpc_bluerpc_proto_init() }
func file_bluerpc_bluerpc_proto_init() {
	if File_bluerpc_bluerpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_bluerpc_bluerpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlueSalt); i {
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
		file_bluerpc_bluerpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlueAppID); i {
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
			RawDescriptor: file_bluerpc_bluerpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_bluerpc_bluerpc_proto_goTypes,
		DependencyIndexes: file_bluerpc_bluerpc_proto_depIdxs,
		MessageInfos:      file_bluerpc_bluerpc_proto_msgTypes,
	}.Build()
	File_bluerpc_bluerpc_proto = out.File
	file_bluerpc_bluerpc_proto_rawDesc = nil
	file_bluerpc_bluerpc_proto_goTypes = nil
	file_bluerpc_bluerpc_proto_depIdxs = nil
}
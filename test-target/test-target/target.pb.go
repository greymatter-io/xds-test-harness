// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: test-target/target.proto

package test_target

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

type ComplimentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *ComplimentRequest) Reset() {
	*x = ComplimentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_target_target_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ComplimentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComplimentRequest) ProtoMessage() {}

func (x *ComplimentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_test_target_target_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComplimentRequest.ProtoReflect.Descriptor instead.
func (*ComplimentRequest) Descriptor() ([]byte, []int) {
	return file_test_target_target_proto_rawDescGZIP(), []int{0}
}

func (x *ComplimentRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type ComplimentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Compliment string `protobuf:"bytes,1,opt,name=compliment,proto3" json:"compliment,omitempty"`
}

func (x *ComplimentResponse) Reset() {
	*x = ComplimentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_target_target_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ComplimentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComplimentResponse) ProtoMessage() {}

func (x *ComplimentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_test_target_target_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComplimentResponse.ProtoReflect.Descriptor instead.
func (*ComplimentResponse) Descriptor() ([]byte, []int) {
	return file_test_target_target_proto_rawDescGZIP(), []int{1}
}

func (x *ComplimentResponse) GetCompliment() string {
	if x != nil {
		return x.Compliment
	}
	return ""
}

var File_test_target_target_proto protoreflect.FileDescriptor

var file_test_target_target_proto_rawDesc = []byte{
	0x0a, 0x18, 0x74, 0x65, 0x73, 0x74, 0x2d, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x2f, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x74, 0x65, 0x73, 0x74,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x22, 0x27, 0x0a, 0x11, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x34, 0x0a, 0x12, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x6d,
	0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x70, 0x6c,
	0x69, 0x6d, 0x65, 0x6e, 0x74, 0x32, 0x59, 0x0a, 0x04, 0x53, 0x68, 0x69, 0x6d, 0x12, 0x51, 0x0a,
	0x0e, 0x47, 0x69, 0x76, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x12,
	0x1d, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x2e, 0x43, 0x6f, 0x6d,
	0x70, 0x6c, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e,
	0x2e, 0x74, 0x65, 0x73, 0x74, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x2e, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x7a,
	0x61, 0x63, 0x68, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x76, 0x69, 0x6c, 0x6c, 0x65, 0x2f, 0x74, 0x65,
	0x73, 0x74, 0x65, 0x72, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x74,
	0x65, 0x73, 0x74, 0x2d, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2d,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_test_target_target_proto_rawDescOnce sync.Once
	file_test_target_target_proto_rawDescData = file_test_target_target_proto_rawDesc
)

func file_test_target_target_proto_rawDescGZIP() []byte {
	file_test_target_target_proto_rawDescOnce.Do(func() {
		file_test_target_target_proto_rawDescData = protoimpl.X.CompressGZIP(file_test_target_target_proto_rawDescData)
	})
	return file_test_target_target_proto_rawDescData
}

var file_test_target_target_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_test_target_target_proto_goTypes = []interface{}{
	(*ComplimentRequest)(nil),  // 0: testtarget.ComplimentRequest
	(*ComplimentResponse)(nil), // 1: testtarget.ComplimentResponse
}
var file_test_target_target_proto_depIdxs = []int32{
	0, // 0: testtarget.Shim.GiveCompliment:input_type -> testtarget.ComplimentRequest
	1, // 1: testtarget.Shim.GiveCompliment:output_type -> testtarget.ComplimentResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_test_target_target_proto_init() }
func file_test_target_target_proto_init() {
	if File_test_target_target_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_test_target_target_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ComplimentRequest); i {
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
		file_test_target_target_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ComplimentResponse); i {
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
			RawDescriptor: file_test_target_target_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_test_target_target_proto_goTypes,
		DependencyIndexes: file_test_target_target_proto_depIdxs,
		MessageInfos:      file_test_target_target_proto_msgTypes,
	}.Build()
	File_test_target_target_proto = out.File
	file_test_target_target_proto_rawDesc = nil
	file_test_target_target_proto_goTypes = nil
	file_test_target_target_proto_depIdxs = nil
}

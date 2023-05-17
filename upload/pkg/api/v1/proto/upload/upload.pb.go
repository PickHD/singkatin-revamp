// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: api/v1/proto/upload/upload.proto

package uploadpb

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

type UploadAvatarMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName    string `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	Avatars     []byte `protobuf:"bytes,2,opt,name=avatars,proto3" json:"avatars,omitempty"`
	ContentType string `protobuf:"bytes,3,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
}

func (x *UploadAvatarMessage) Reset() {
	*x = UploadAvatarMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_proto_upload_upload_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadAvatarMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadAvatarMessage) ProtoMessage() {}

func (x *UploadAvatarMessage) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_proto_upload_upload_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadAvatarMessage.ProtoReflect.Descriptor instead.
func (*UploadAvatarMessage) Descriptor() ([]byte, []int) {
	return file_api_v1_proto_upload_upload_proto_rawDescGZIP(), []int{0}
}

func (x *UploadAvatarMessage) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *UploadAvatarMessage) GetAvatars() []byte {
	if x != nil {
		return x.Avatars
	}
	return nil
}

func (x *UploadAvatarMessage) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

var File_api_v1_proto_upload_upload_proto protoreflect.FileDescriptor

var file_api_v1_proto_upload_upload_proto_rawDesc = []byte{
	0x0a, 0x20, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x2f, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x13, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x6f, 0x0a, 0x13, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x61, 0x76,
	0x61, 0x74, 0x61, 0x72, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x42, 0x4c, 0x5a, 0x4a, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x50, 0x69, 0x63, 0x6b, 0x48, 0x44, 0x2f, 0x73, 0x69,
	0x6e, 0x67, 0x6b, 0x61, 0x74, 0x69, 0x6e, 0x2d, 0x72, 0x65, 0x76, 0x61, 0x6d, 0x70, 0x2f, 0x75,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x3b, 0x75, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_proto_upload_upload_proto_rawDescOnce sync.Once
	file_api_v1_proto_upload_upload_proto_rawDescData = file_api_v1_proto_upload_upload_proto_rawDesc
)

func file_api_v1_proto_upload_upload_proto_rawDescGZIP() []byte {
	file_api_v1_proto_upload_upload_proto_rawDescOnce.Do(func() {
		file_api_v1_proto_upload_upload_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_proto_upload_upload_proto_rawDescData)
	})
	return file_api_v1_proto_upload_upload_proto_rawDescData
}

var file_api_v1_proto_upload_upload_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_v1_proto_upload_upload_proto_goTypes = []interface{}{
	(*UploadAvatarMessage)(nil), // 0: api.v1.proto.upload.UploadAvatarMessage
}
var file_api_v1_proto_upload_upload_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_v1_proto_upload_upload_proto_init() }
func file_api_v1_proto_upload_upload_proto_init() {
	if File_api_v1_proto_upload_upload_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_proto_upload_upload_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadAvatarMessage); i {
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
			RawDescriptor: file_api_v1_proto_upload_upload_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1_proto_upload_upload_proto_goTypes,
		DependencyIndexes: file_api_v1_proto_upload_upload_proto_depIdxs,
		MessageInfos:      file_api_v1_proto_upload_upload_proto_msgTypes,
	}.Build()
	File_api_v1_proto_upload_upload_proto = out.File
	file_api_v1_proto_upload_upload_proto_rawDesc = nil
	file_api_v1_proto_upload_upload_proto_goTypes = nil
	file_api_v1_proto_upload_upload_proto_depIdxs = nil
}
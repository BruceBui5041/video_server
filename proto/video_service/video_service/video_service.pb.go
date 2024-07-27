// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: video_service/video_service.proto

package proto

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

type VideoInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VideoId     string `protobuf:"bytes,1,opt,name=video_id,json=videoId,proto3" json:"video_id,omitempty"`
	Title       string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	UploadedBy  string `protobuf:"bytes,4,opt,name=uploaded_by,json=uploadedBy,proto3" json:"uploaded_by,omitempty"`
	Timestamp   int64  `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	S3Key       string `protobuf:"bytes,6,opt,name=s3_key,json=s3Key,proto3" json:"s3_key,omitempty"`
	Slug        string `protobuf:"bytes,7,opt,name=slug,proto3" json:"slug,omitempty"`
}

func (x *VideoInfo) Reset() {
	*x = VideoInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_video_service_video_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VideoInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoInfo) ProtoMessage() {}

func (x *VideoInfo) ProtoReflect() protoreflect.Message {
	mi := &file_video_service_video_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoInfo.ProtoReflect.Descriptor instead.
func (*VideoInfo) Descriptor() ([]byte, []int) {
	return file_video_service_video_service_proto_rawDescGZIP(), []int{0}
}

func (x *VideoInfo) GetVideoId() string {
	if x != nil {
		return x.VideoId
	}
	return ""
}

func (x *VideoInfo) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *VideoInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *VideoInfo) GetUploadedBy() string {
	if x != nil {
		return x.UploadedBy
	}
	return ""
}

func (x *VideoInfo) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *VideoInfo) GetS3Key() string {
	if x != nil {
		return x.S3Key
	}
	return ""
}

func (x *VideoInfo) GetSlug() string {
	if x != nil {
		return x.Slug
	}
	return ""
}

type ProcessNewVideoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ProcessNewVideoResponse) Reset() {
	*x = ProcessNewVideoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_video_service_video_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessNewVideoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessNewVideoResponse) ProtoMessage() {}

func (x *ProcessNewVideoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_video_service_video_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessNewVideoResponse.ProtoReflect.Descriptor instead.
func (*ProcessNewVideoResponse) Descriptor() ([]byte, []int) {
	return file_video_service_video_service_proto_rawDescGZIP(), []int{1}
}

func (x *ProcessNewVideoResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_video_service_video_service_proto protoreflect.FileDescriptor

var file_video_service_video_service_proto_rawDesc = []byte{
	0x0a, 0x21, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x22, 0xc8, 0x01, 0x0a, 0x09, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x19, 0x0a, 0x08, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x5f, 0x62,
	0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65,
	0x64, 0x42, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x33, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x73, 0x33, 0x4b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6c, 0x75, 0x67,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x22, 0x31, 0x0a, 0x17,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x4e, 0x65, 0x77, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32,
	0x74, 0x0a, 0x16, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69,
	0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5a, 0x0a, 0x16, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x4e, 0x65, 0x77, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x17, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x25, 0x2e, 0x76,
	0x69, 0x64, 0x65, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x4e, 0x65, 0x77, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x14, 0x5a, 0x12, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_video_service_video_service_proto_rawDescOnce sync.Once
	file_video_service_video_service_proto_rawDescData = file_video_service_video_service_proto_rawDesc
)

func file_video_service_video_service_proto_rawDescGZIP() []byte {
	file_video_service_video_service_proto_rawDescOnce.Do(func() {
		file_video_service_video_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_video_service_video_service_proto_rawDescData)
	})
	return file_video_service_video_service_proto_rawDescData
}

var file_video_service_video_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_video_service_video_service_proto_goTypes = []any{
	(*VideoInfo)(nil),               // 0: videoservice.VideoInfo
	(*ProcessNewVideoResponse)(nil), // 1: videoservice.ProcessNewVideoResponse
}
var file_video_service_video_service_proto_depIdxs = []int32{
	0, // 0: videoservice.VideoProcessingService.ProcessNewVideoRequest:input_type -> videoservice.VideoInfo
	1, // 1: videoservice.VideoProcessingService.ProcessNewVideoRequest:output_type -> videoservice.ProcessNewVideoResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_video_service_video_service_proto_init() }
func file_video_service_video_service_proto_init() {
	if File_video_service_video_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_video_service_video_service_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*VideoInfo); i {
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
		file_video_service_video_service_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ProcessNewVideoResponse); i {
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
			RawDescriptor: file_video_service_video_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_video_service_video_service_proto_goTypes,
		DependencyIndexes: file_video_service_video_service_proto_depIdxs,
		MessageInfos:      file_video_service_video_service_proto_msgTypes,
	}.Build()
	File_video_service_video_service_proto = out.File
	file_video_service_video_service_proto_rawDesc = nil
	file_video_service_video_service_proto_goTypes = nil
	file_video_service_video_service_proto_depIdxs = nil
}
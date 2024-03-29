// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: repositories.proto

package service

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type RepoGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	CreatorId string `protobuf:"bytes,1,opt,name=creator_id,json=creatorId,proto3" json:"creator_id,omitempty"`
}

func (x *RepoGetRequest) Reset() {
	*x = RepoGetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repositories_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepoGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepoGetRequest) ProtoMessage() {}

func (x *RepoGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_repositories_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepoGetRequest.ProtoReflect.Descriptor instead.
func (*RepoGetRequest) Descriptor() ([]byte, []int) {
	return file_repositories_proto_rawDescGZIP(), []int{0}
}

func (x *RepoGetRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RepoGetRequest) GetCreatorId() string {
	if x != nil {
		return x.CreatorId
	}
	return ""
}

type Repository struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Url   string `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Owner *User  `protobuf:"bytes,4,opt,name=owner,proto3" json:"owner,omitempty"`
}

func (x *Repository) Reset() {
	*x = Repository{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repositories_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Repository) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Repository) ProtoMessage() {}

func (x *Repository) ProtoReflect() protoreflect.Message {
	mi := &file_repositories_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Repository.ProtoReflect.Descriptor instead.
func (*Repository) Descriptor() ([]byte, []int) {
	return file_repositories_proto_rawDescGZIP(), []int{1}
}

func (x *Repository) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Repository) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Repository) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Repository) GetOwner() *User {
	if x != nil {
		return x.Owner
	}
	return nil
}

type RepoGetReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repo *Repository `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
}

func (x *RepoGetReply) Reset() {
	*x = RepoGetReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repositories_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepoGetReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepoGetReply) ProtoMessage() {}

func (x *RepoGetReply) ProtoReflect() protoreflect.Message {
	mi := &file_repositories_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepoGetReply.ProtoReflect.Descriptor instead.
func (*RepoGetReply) Descriptor() ([]byte, []int) {
	return file_repositories_proto_rawDescGZIP(), []int{2}
}

func (x *RepoGetReply) GetRepo() *Repository {
	if x != nil {
		return x.Repo
	}
	return nil
}

type RepoBuildLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogLine   string                 `protobuf:"bytes,1,opt,name=log_line,json=logLine,proto3" json:"log_line,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *RepoBuildLog) Reset() {
	*x = RepoBuildLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repositories_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepoBuildLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepoBuildLog) ProtoMessage() {}

func (x *RepoBuildLog) ProtoReflect() protoreflect.Message {
	mi := &file_repositories_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepoBuildLog.ProtoReflect.Descriptor instead.
func (*RepoBuildLog) Descriptor() ([]byte, []int) {
	return file_repositories_proto_rawDescGZIP(), []int{3}
}

func (x *RepoBuildLog) GetLogLine() string {
	if x != nil {
		return x.LogLine
	}
	return ""
}

func (x *RepoBuildLog) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

var File_repositories_proto protoreflect.FileDescriptor

var file_repositories_proto_rawDesc = []byte{
	0x0a, 0x12, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x0b, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x0e, 0x52,
	0x65, 0x70, 0x6f, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x22, 0x67, 0x0a, 0x0a,
	0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c,
	0x12, 0x23, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x05,
	0x6f, 0x77, 0x6e, 0x65, 0x72, 0x22, 0x37, 0x0a, 0x0c, 0x52, 0x65, 0x70, 0x6f, 0x47, 0x65, 0x74,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x27, 0x0a, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x22, 0x63,
	0x0a, 0x0c, 0x52, 0x65, 0x70, 0x6f, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x4c, 0x6f, 0x67, 0x12, 0x19,
	0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x5f, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6c, 0x6f, 0x67, 0x4c, 0x69, 0x6e, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x32, 0x85, 0x01, 0x0a, 0x04, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x3e, 0x0a, 0x08,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x12, 0x17, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x15, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x70, 0x6f,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x30, 0x01, 0x12, 0x3d, 0x0a, 0x0b,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x12, 0x13, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x1a, 0x15, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x42,
	0x75, 0x69, 0x6c, 0x64, 0x4c, 0x6f, 0x67, 0x22, 0x00, 0x30, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_repositories_proto_rawDescOnce sync.Once
	file_repositories_proto_rawDescData = file_repositories_proto_rawDesc
)

func file_repositories_proto_rawDescGZIP() []byte {
	file_repositories_proto_rawDescOnce.Do(func() {
		file_repositories_proto_rawDescData = protoimpl.X.CompressGZIP(file_repositories_proto_rawDescData)
	})
	return file_repositories_proto_rawDescData
}

var file_repositories_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_repositories_proto_goTypes = []interface{}{
	(*RepoGetRequest)(nil),        // 0: service.RepoGetRequest
	(*Repository)(nil),            // 1: service.Repository
	(*RepoGetReply)(nil),          // 2: service.RepoGetReply
	(*RepoBuildLog)(nil),          // 3: service.RepoBuildLog
	(*User)(nil),                  // 4: service.User
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_repositories_proto_depIdxs = []int32{
	4, // 0: service.Repository.owner:type_name -> service.User
	1, // 1: service.RepoGetReply.repo:type_name -> service.Repository
	5, // 2: service.RepoBuildLog.timestamp:type_name -> google.protobuf.Timestamp
	0, // 3: service.Repo.GetRepos:input_type -> service.RepoGetRequest
	1, // 4: service.Repo.CreateBuild:input_type -> service.Repository
	2, // 5: service.Repo.GetRepos:output_type -> service.RepoGetReply
	3, // 6: service.Repo.CreateBuild:output_type -> service.RepoBuildLog
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_repositories_proto_init() }
func file_repositories_proto_init() {
	if File_repositories_proto != nil {
		return
	}
	file_users_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_repositories_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepoGetRequest); i {
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
		file_repositories_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Repository); i {
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
		file_repositories_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepoGetReply); i {
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
		file_repositories_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepoBuildLog); i {
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
			RawDescriptor: file_repositories_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_repositories_proto_goTypes,
		DependencyIndexes: file_repositories_proto_depIdxs,
		MessageInfos:      file_repositories_proto_msgTypes,
	}.Build()
	File_repositories_proto = out.File
	file_repositories_proto_rawDesc = nil
	file_repositories_proto_goTypes = nil
	file_repositories_proto_depIdxs = nil
}

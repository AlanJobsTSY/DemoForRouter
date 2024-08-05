// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.2
// source: minigame_router.proto

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

type Service struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Ip       string `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	Port     string `protobuf:"bytes,3,opt,name=port,proto3" json:"port,omitempty"`
	Protocol string `protobuf:"bytes,4,opt,name=protocol,proto3" json:"protocol,omitempty"`
	Weight   string `protobuf:"bytes,5,opt,name=weight,proto3" json:"weight,omitempty"`
}

func (x *Service) Reset() {
	*x = Service{}
	if protoimpl.UnsafeEnabled {
		mi := &file_minigame_router_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Service) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Service) ProtoMessage() {}

func (x *Service) ProtoReflect() protoreflect.Message {
	mi := &file_minigame_router_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Service.ProtoReflect.Descriptor instead.
func (*Service) Descriptor() ([]byte, []int) {
	return file_minigame_router_proto_rawDescGZIP(), []int{0}
}

func (x *Service) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Service) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *Service) GetPort() string {
	if x != nil {
		return x.Port
	}
	return ""
}

func (x *Service) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

func (x *Service) GetWeight() string {
	if x != nil {
		return x.Weight
	}
	return ""
}

type RegisterServiceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Service *Service `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
}

func (x *RegisterServiceRequest) Reset() {
	*x = RegisterServiceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_minigame_router_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterServiceRequest) ProtoMessage() {}

func (x *RegisterServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_minigame_router_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterServiceRequest.ProtoReflect.Descriptor instead.
func (*RegisterServiceRequest) Descriptor() ([]byte, []int) {
	return file_minigame_router_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterServiceRequest) GetService() *Service {
	if x != nil {
		return x.Service
	}
	return nil
}

type RegisterServiceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *RegisterServiceResponse) Reset() {
	*x = RegisterServiceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_minigame_router_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterServiceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterServiceResponse) ProtoMessage() {}

func (x *RegisterServiceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_minigame_router_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterServiceResponse.ProtoReflect.Descriptor instead.
func (*RegisterServiceResponse) Descriptor() ([]byte, []int) {
	return file_minigame_router_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterServiceResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type DiscoverServiceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromMsg string `protobuf:"bytes,1,opt,name=from_msg,json=fromMsg,proto3" json:"from_msg,omitempty"`
	ToMsg   string `protobuf:"bytes,2,opt,name=to_msg,json=toMsg,proto3" json:"to_msg,omitempty"`
}

func (x *DiscoverServiceRequest) Reset() {
	*x = DiscoverServiceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_minigame_router_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoverServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoverServiceRequest) ProtoMessage() {}

func (x *DiscoverServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_minigame_router_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoverServiceRequest.ProtoReflect.Descriptor instead.
func (*DiscoverServiceRequest) Descriptor() ([]byte, []int) {
	return file_minigame_router_proto_rawDescGZIP(), []int{3}
}

func (x *DiscoverServiceRequest) GetFromMsg() string {
	if x != nil {
		return x.FromMsg
	}
	return ""
}

func (x *DiscoverServiceRequest) GetToMsg() string {
	if x != nil {
		return x.ToMsg
	}
	return ""
}

type DiscoverServiceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *DiscoverServiceResponse) Reset() {
	*x = DiscoverServiceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_minigame_router_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoverServiceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoverServiceResponse) ProtoMessage() {}

func (x *DiscoverServiceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_minigame_router_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoverServiceResponse.ProtoReflect.Descriptor instead.
func (*DiscoverServiceResponse) Descriptor() ([]byte, []int) {
	return file_minigame_router_proto_rawDescGZIP(), []int{4}
}

func (x *DiscoverServiceResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type HelloRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *HelloRequest) Reset() {
	*x = HelloRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_minigame_router_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloRequest) ProtoMessage() {}

func (x *HelloRequest) ProtoReflect() protoreflect.Message {
	mi := &file_minigame_router_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloRequest.ProtoReflect.Descriptor instead.
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return file_minigame_router_proto_rawDescGZIP(), []int{5}
}

func (x *HelloRequest) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type HelloResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *HelloResponse) Reset() {
	*x = HelloResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_minigame_router_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloResponse) ProtoMessage() {}

func (x *HelloResponse) ProtoReflect() protoreflect.Message {
	mi := &file_minigame_router_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloResponse.ProtoReflect.Descriptor instead.
func (*HelloResponse) Descriptor() ([]byte, []int) {
	return file_minigame_router_proto_rawDescGZIP(), []int{6}
}

func (x *HelloResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_minigame_router_proto protoreflect.FileDescriptor

var file_minigame_router_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6d, 0x69, 0x6e, 0x69, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x75,
	0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x6f, 0x72,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x16, 0x0a,
	0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x77,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x42, 0x0a, 0x16, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x28, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x2b, 0x0a, 0x17, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x4a, 0x0a, 0x16, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x19, 0x0a, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x66, 0x72, 0x6f, 0x6d, 0x4d, 0x73, 0x67, 0x12, 0x15, 0x0a, 0x06, 0x74,
	0x6f, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x4d,
	0x73, 0x67, 0x22, 0x2b, 0x0a, 0x17, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22,
	0x20, 0x0a, 0x0c, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73,
	0x67, 0x22, 0x21, 0x0a, 0x0d, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6d, 0x73, 0x67, 0x32, 0xeb, 0x01, 0x0a, 0x0e, 0x4d, 0x69, 0x6e, 0x69, 0x47, 0x61, 0x6d,
	0x65, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x12, 0x50, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a, 0x0f, 0x44, 0x69, 0x73,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x53,
	0x61, 0x79, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x0f, 0x5a, 0x0d, 0x4d, 0x6f, 0x72, 0x65, 0x46, 0x75, 0x6e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_minigame_router_proto_rawDescOnce sync.Once
	file_minigame_router_proto_rawDescData = file_minigame_router_proto_rawDesc
)

func file_minigame_router_proto_rawDescGZIP() []byte {
	file_minigame_router_proto_rawDescOnce.Do(func() {
		file_minigame_router_proto_rawDescData = protoimpl.X.CompressGZIP(file_minigame_router_proto_rawDescData)
	})
	return file_minigame_router_proto_rawDescData
}

var file_minigame_router_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_minigame_router_proto_goTypes = []any{
	(*Service)(nil),                 // 0: proto.Service
	(*RegisterServiceRequest)(nil),  // 1: proto.RegisterServiceRequest
	(*RegisterServiceResponse)(nil), // 2: proto.RegisterServiceResponse
	(*DiscoverServiceRequest)(nil),  // 3: proto.DiscoverServiceRequest
	(*DiscoverServiceResponse)(nil), // 4: proto.DiscoverServiceResponse
	(*HelloRequest)(nil),            // 5: proto.HelloRequest
	(*HelloResponse)(nil),           // 6: proto.HelloResponse
}
var file_minigame_router_proto_depIdxs = []int32{
	0, // 0: proto.RegisterServiceRequest.service:type_name -> proto.Service
	1, // 1: proto.MiniGameRouter.RegisterService:input_type -> proto.RegisterServiceRequest
	3, // 2: proto.MiniGameRouter.DiscoverService:input_type -> proto.DiscoverServiceRequest
	5, // 3: proto.MiniGameRouter.SayHello:input_type -> proto.HelloRequest
	2, // 4: proto.MiniGameRouter.RegisterService:output_type -> proto.RegisterServiceResponse
	4, // 5: proto.MiniGameRouter.DiscoverService:output_type -> proto.DiscoverServiceResponse
	6, // 6: proto.MiniGameRouter.SayHello:output_type -> proto.HelloResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_minigame_router_proto_init() }
func file_minigame_router_proto_init() {
	if File_minigame_router_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_minigame_router_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Service); i {
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
		file_minigame_router_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*RegisterServiceRequest); i {
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
		file_minigame_router_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*RegisterServiceResponse); i {
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
		file_minigame_router_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*DiscoverServiceRequest); i {
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
		file_minigame_router_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*DiscoverServiceResponse); i {
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
		file_minigame_router_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*HelloRequest); i {
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
		file_minigame_router_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*HelloResponse); i {
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
			RawDescriptor: file_minigame_router_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_minigame_router_proto_goTypes,
		DependencyIndexes: file_minigame_router_proto_depIdxs,
		MessageInfos:      file_minigame_router_proto_msgTypes,
	}.Build()
	File_minigame_router_proto = out.File
	file_minigame_router_proto_rawDesc = nil
	file_minigame_router_proto_goTypes = nil
	file_minigame_router_proto_depIdxs = nil
}

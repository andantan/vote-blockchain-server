// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: vote_message.proto

package vote_message

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VoteRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Hash          string                 `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Option        string                 `protobuf:"bytes,2,opt,name=option,proto3" json:"option,omitempty"`
	Topic         string                 `protobuf:"bytes,3,opt,name=topic,proto3" json:"topic,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VoteRequest) Reset() {
	*x = VoteRequest{}
	mi := &file_vote_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequest) ProtoMessage() {}

func (x *VoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vote_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequest.ProtoReflect.Descriptor instead.
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return file_vote_message_proto_rawDescGZIP(), []int{0}
}

func (x *VoteRequest) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *VoteRequest) GetOption() string {
	if x != nil {
		return x.Option
	}
	return ""
}

func (x *VoteRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

type VoteResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Success       bool                   `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VoteResponse) Reset() {
	*x = VoteResponse{}
	mi := &file_vote_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteResponse) ProtoMessage() {}

func (x *VoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vote_message_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteResponse.ProtoReflect.Descriptor instead.
func (*VoteResponse) Descriptor() ([]byte, []int) {
	return file_vote_message_proto_rawDescGZIP(), []int{1}
}

func (x *VoteResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *VoteResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *VoteResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_vote_message_proto protoreflect.FileDescriptor

const file_vote_message_proto_rawDesc = "" +
	"\n" +
	"\x12vote_message.proto\x12\fvote_message\"O\n" +
	"\vVoteRequest\x12\x12\n" +
	"\x04hash\x18\x01 \x01(\tR\x04hash\x12\x16\n" +
	"\x06option\x18\x02 \x01(\tR\x06option\x12\x14\n" +
	"\x05topic\x18\x03 \x01(\tR\x05topic\"Z\n" +
	"\fVoteResponse\x12\x16\n" +
	"\x06status\x18\x01 \x01(\tR\x06status\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\x12\x18\n" +
	"\asuccess\x18\x03 \x01(\bR\asuccess2\\\n" +
	"\x15BlockchainVoteService\x12C\n" +
	"\n" +
	"SubmitVote\x12\x19.vote_message.VoteRequest\x1a\x1a.vote_message.VoteResponseB\x1eZ\x1c../network/gRPC/vote_messageb\x06proto3"

var (
	file_vote_message_proto_rawDescOnce sync.Once
	file_vote_message_proto_rawDescData []byte
)

func file_vote_message_proto_rawDescGZIP() []byte {
	file_vote_message_proto_rawDescOnce.Do(func() {
		file_vote_message_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_vote_message_proto_rawDesc), len(file_vote_message_proto_rawDesc)))
	})
	return file_vote_message_proto_rawDescData
}

var file_vote_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_vote_message_proto_goTypes = []any{
	(*VoteRequest)(nil),  // 0: vote_message.VoteRequest
	(*VoteResponse)(nil), // 1: vote_message.VoteResponse
}
var file_vote_message_proto_depIdxs = []int32{
	0, // 0: vote_message.BlockchainVoteService.SubmitVote:input_type -> vote_message.VoteRequest
	1, // 1: vote_message.BlockchainVoteService.SubmitVote:output_type -> vote_message.VoteResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_vote_message_proto_init() }
func file_vote_message_proto_init() {
	if File_vote_message_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_vote_message_proto_rawDesc), len(file_vote_message_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_vote_message_proto_goTypes,
		DependencyIndexes: file_vote_message_proto_depIdxs,
		MessageInfos:      file_vote_message_proto_msgTypes,
	}.Build()
	File_vote_message_proto = out.File
	file_vote_message_proto_goTypes = nil
	file_vote_message_proto_depIdxs = nil
}

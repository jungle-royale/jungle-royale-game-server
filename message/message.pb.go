// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        v5.29.2
// source: message.proto

package message

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

type Wrapper struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to MessageType:
	//
	//	*Wrapper_Change
	//	*Wrapper_State
	//	*Wrapper_Gameinit
	MessageType   isWrapper_MessageType `protobuf_oneof:"MessageType"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Wrapper) Reset() {
	*x = Wrapper{}
	mi := &file_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Wrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Wrapper) ProtoMessage() {}

func (x *Wrapper) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Wrapper.ProtoReflect.Descriptor instead.
func (*Wrapper) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0}
}

func (x *Wrapper) GetMessageType() isWrapper_MessageType {
	if x != nil {
		return x.MessageType
	}
	return nil
}

func (x *Wrapper) GetChange() *Change {
	if x != nil {
		if x, ok := x.MessageType.(*Wrapper_Change); ok {
			return x.Change
		}
	}
	return nil
}

func (x *Wrapper) GetState() *GameState {
	if x != nil {
		if x, ok := x.MessageType.(*Wrapper_State); ok {
			return x.State
		}
	}
	return nil
}

func (x *Wrapper) GetGameinit() *GameInit {
	if x != nil {
		if x, ok := x.MessageType.(*Wrapper_Gameinit); ok {
			return x.Gameinit
		}
	}
	return nil
}

type isWrapper_MessageType interface {
	isWrapper_MessageType()
}

type Wrapper_Change struct {
	Change *Change `protobuf:"bytes,1,opt,name=change,proto3,oneof"`
}

type Wrapper_State struct {
	State *GameState `protobuf:"bytes,2,opt,name=state,proto3,oneof"`
}

type Wrapper_Gameinit struct {
	Gameinit *GameInit `protobuf:"bytes,3,opt,name=gameinit,proto3,oneof"`
}

func (*Wrapper_Change) isWrapper_MessageType() {}

func (*Wrapper_State) isWrapper_MessageType() {}

func (*Wrapper_Gameinit) isWrapper_MessageType() {}

type GameInit struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GameInit) Reset() {
	*x = GameInit{}
	mi := &file_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GameInit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameInit) ProtoMessage() {}

func (x *GameInit) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameInit.ProtoReflect.Descriptor instead.
func (*GameInit) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{1}
}

func (x *GameInit) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Change struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	X             float32                `protobuf:"fixed32,1,opt,name=x,proto3" json:"x,omitempty"`
	Y             float32                `protobuf:"fixed32,2,opt,name=y,proto3" json:"y,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Change) Reset() {
	*x = Change{}
	mi := &file_message_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Change) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Change) ProtoMessage() {}

func (x *Change) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Change.ProtoReflect.Descriptor instead.
func (*Change) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{2}
}

func (x *Change) GetX() float32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Change) GetY() float32 {
	if x != nil {
		return x.Y
	}
	return 0
}

type GameState struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Players       []*Player              `protobuf:"bytes,1,rep,name=players,proto3" json:"players,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GameState) Reset() {
	*x = GameState{}
	mi := &file_message_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GameState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameState) ProtoMessage() {}

func (x *GameState) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameState.ProtoReflect.Descriptor instead.
func (*GameState) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{3}
}

func (x *GameState) GetPlayers() []*Player {
	if x != nil {
		return x.Players
	}
	return nil
}

type Player struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	X             int32                  `protobuf:"varint,2,opt,name=x,proto3" json:"x,omitempty"`
	Y             int32                  `protobuf:"varint,3,opt,name=y,proto3" json:"y,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Player) Reset() {
	*x = Player{}
	mi := &file_message_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Player) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Player) ProtoMessage() {}

func (x *Player) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Player.ProtoReflect.Descriptor instead.
func (*Player) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{4}
}

func (x *Player) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Player) GetX() int32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Player) GetY() int32 {
	if x != nil {
		return x.Y
	}
	return 0
}

var File_message_proto protoreflect.FileDescriptor

var file_message_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xa0, 0x01, 0x0a, 0x07, 0x57, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x72, 0x12, 0x29, 0x0a, 0x06, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x48, 0x00, 0x52, 0x06, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12,
	0x2a, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x48, 0x00, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x2f, 0x0a, 0x08, 0x67,
	0x61, 0x6d, 0x65, 0x69, 0x6e, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x49, 0x6e, 0x69, 0x74,
	0x48, 0x00, 0x52, 0x08, 0x67, 0x61, 0x6d, 0x65, 0x69, 0x6e, 0x69, 0x74, 0x42, 0x0d, 0x0a, 0x0b,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x22, 0x1a, 0x0a, 0x08, 0x47,
	0x61, 0x6d, 0x65, 0x49, 0x6e, 0x69, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x24, 0x0a, 0x06, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x78, 0x12,
	0x0c, 0x0a, 0x01, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x79, 0x22, 0x36, 0x0a,
	0x09, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x29, 0x0a, 0x07, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x07, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x73, 0x22, 0x34, 0x0a, 0x06, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x0c, 0x0a, 0x01, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a,
	0x01, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x79, 0x42, 0x0a, 0x5a, 0x08, 0x2f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_proto_rawDescOnce sync.Once
	file_message_proto_rawDescData = file_message_proto_rawDesc
)

func file_message_proto_rawDescGZIP() []byte {
	file_message_proto_rawDescOnce.Do(func() {
		file_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_proto_rawDescData)
	})
	return file_message_proto_rawDescData
}

var file_message_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_message_proto_goTypes = []any{
	(*Wrapper)(nil),   // 0: message.Wrapper
	(*GameInit)(nil),  // 1: message.GameInit
	(*Change)(nil),    // 2: message.Change
	(*GameState)(nil), // 3: message.GameState
	(*Player)(nil),    // 4: message.Player
}
var file_message_proto_depIdxs = []int32{
	2, // 0: message.Wrapper.change:type_name -> message.Change
	3, // 1: message.Wrapper.state:type_name -> message.GameState
	1, // 2: message.Wrapper.gameinit:type_name -> message.GameInit
	4, // 3: message.GameState.players:type_name -> message.Player
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_message_proto_init() }
func file_message_proto_init() {
	if File_message_proto != nil {
		return
	}
	file_message_proto_msgTypes[0].OneofWrappers = []any{
		(*Wrapper_Change)(nil),
		(*Wrapper_State)(nil),
		(*Wrapper_Gameinit)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_proto_goTypes,
		DependencyIndexes: file_message_proto_depIdxs,
		MessageInfos:      file_message_proto_msgTypes,
	}.Build()
	File_message_proto = out.File
	file_message_proto_rawDesc = nil
	file_message_proto_goTypes = nil
	file_message_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.29.2
// source: message/message.proto

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
	//	*Wrapper_ChangeDir
	//	*Wrapper_DoDash
	//	*Wrapper_CreateBullet
	//	*Wrapper_GameState
	//	*Wrapper_GameCount
	//	*Wrapper_GameInit
	//	*Wrapper_GameStart
	MessageType   isWrapper_MessageType `protobuf_oneof:"MessageType"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Wrapper) Reset() {
	*x = Wrapper{}
	mi := &file_message_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Wrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Wrapper) ProtoMessage() {}

func (x *Wrapper) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[0]
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
	return file_message_message_proto_rawDescGZIP(), []int{0}
}

func (x *Wrapper) GetMessageType() isWrapper_MessageType {
	if x != nil {
		return x.MessageType
	}
	return nil
}

func (x *Wrapper) GetChangeDir() *ChangeDir {
	if x != nil {
		if x, ok := x.MessageType.(*Wrapper_ChangeDir); ok {
			return x.ChangeDir
		}
	}
	return nil
}

func (x *Wrapper) GetDoDash() *DoDash {
	if x != nil {
		if x, ok := x.MessageType.(*Wrapper_DoDash); ok {
			return x.DoDash
		}
	}
	return nil
}

func (x *Wrapper) GetCreateBullet() *CreateBullet {
	if x != nil {
		if x, ok := x.MessageType.(*Wrapper_CreateBullet); ok {
			return x.CreateBullet
		}
	}
	return nil
}

func (x *Wrapper) GetGameState() *GameState {
	if x != nil {
		if x, ok := x.MessageType.(*Wrapper_GameState); ok {
			return x.GameState
		}
	}
	return nil
}

func (x *Wrapper) GetGameCount() *GameCount {
	if x != nil {
		if x, ok := x.MessageType.(*Wrapper_GameCount); ok {
			return x.GameCount
		}
	}
	return nil
}

func (x *Wrapper) GetGameInit() *GameInit {
	if x != nil {
		if x, ok := x.MessageType.(*Wrapper_GameInit); ok {
			return x.GameInit
		}
	}
	return nil
}

func (x *Wrapper) GetGameStart() *GameStart {
	if x != nil {
		if x, ok := x.MessageType.(*Wrapper_GameStart); ok {
			return x.GameStart
		}
	}
	return nil
}

type isWrapper_MessageType interface {
	isWrapper_MessageType()
}

type Wrapper_ChangeDir struct {
	ChangeDir *ChangeDir `protobuf:"bytes,1,opt,name=changeDir,proto3,oneof"`
}

type Wrapper_DoDash struct {
	DoDash *DoDash `protobuf:"bytes,2,opt,name=doDash,proto3,oneof"`
}

type Wrapper_CreateBullet struct {
	CreateBullet *CreateBullet `protobuf:"bytes,3,opt,name=createBullet,proto3,oneof"`
}

type Wrapper_GameState struct {
	GameState *GameState `protobuf:"bytes,4,opt,name=gameState,proto3,oneof"`
}

type Wrapper_GameCount struct {
	GameCount *GameCount `protobuf:"bytes,5,opt,name=gameCount,proto3,oneof"`
}

type Wrapper_GameInit struct {
	GameInit *GameInit `protobuf:"bytes,6,opt,name=gameInit,proto3,oneof"`
}

type Wrapper_GameStart struct {
	GameStart *GameStart `protobuf:"bytes,7,opt,name=gameStart,proto3,oneof"`
}

func (*Wrapper_ChangeDir) isWrapper_MessageType() {}

func (*Wrapper_DoDash) isWrapper_MessageType() {}

func (*Wrapper_CreateBullet) isWrapper_MessageType() {}

func (*Wrapper_GameState) isWrapper_MessageType() {}

func (*Wrapper_GameCount) isWrapper_MessageType() {}

func (*Wrapper_GameInit) isWrapper_MessageType() {}

func (*Wrapper_GameStart) isWrapper_MessageType() {}

// server → client : init message
type GameInit struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // player id
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GameInit) Reset() {
	*x = GameInit{}
	mi := &file_message_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GameInit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameInit) ProtoMessage() {}

func (x *GameInit) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[1]
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
	return file_message_message_proto_rawDescGZIP(), []int{1}
}

func (x *GameInit) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// client → server : changing state
type ChangeDir struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Angle         float32                `protobuf:"fixed32,1,opt,name=angle,proto3" json:"angle,omitempty"`
	IsMoved       bool                   `protobuf:"varint,2,opt,name=isMoved,proto3" json:"isMoved,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangeDir) Reset() {
	*x = ChangeDir{}
	mi := &file_message_message_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangeDir) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeDir) ProtoMessage() {}

func (x *ChangeDir) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeDir.ProtoReflect.Descriptor instead.
func (*ChangeDir) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{2}
}

func (x *ChangeDir) GetAngle() float32 {
	if x != nil {
		return x.Angle
	}
	return 0
}

func (x *ChangeDir) GetIsMoved() bool {
	if x != nil {
		return x.IsMoved
	}
	return false
}

// server → client
type GameState struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	PlayerState     []*PlayerState         `protobuf:"bytes,1,rep,name=playerState,proto3" json:"playerState,omitempty"`
	BulletState     []*BulletState         `protobuf:"bytes,2,rep,name=bulletState,proto3" json:"bulletState,omitempty"`
	HealPackState   []*HealPackState       `protobuf:"bytes,3,rep,name=healPackState,proto3" json:"healPackState,omitempty"`
	MagicItemState  []*MagicItemState      `protobuf:"bytes,4,rep,name=magicItemState,proto3" json:"magicItemState,omitempty"`
	PlayerDeadState []*PlayerDeadState     `protobuf:"bytes,5,rep,name=playerDeadState,proto3" json:"playerDeadState,omitempty"`
	TileState       []*TileState           `protobuf:"bytes,6,rep,name=tileState,proto3" json:"tileState,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *GameState) Reset() {
	*x = GameState{}
	mi := &file_message_message_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GameState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameState) ProtoMessage() {}

func (x *GameState) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[3]
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
	return file_message_message_proto_rawDescGZIP(), []int{3}
}

func (x *GameState) GetPlayerState() []*PlayerState {
	if x != nil {
		return x.PlayerState
	}
	return nil
}

func (x *GameState) GetBulletState() []*BulletState {
	if x != nil {
		return x.BulletState
	}
	return nil
}

func (x *GameState) GetHealPackState() []*HealPackState {
	if x != nil {
		return x.HealPackState
	}
	return nil
}

func (x *GameState) GetMagicItemState() []*MagicItemState {
	if x != nil {
		return x.MagicItemState
	}
	return nil
}

func (x *GameState) GetPlayerDeadState() []*PlayerDeadState {
	if x != nil {
		return x.PlayerDeadState
	}
	return nil
}

func (x *GameState) GetTileState() []*TileState {
	if x != nil {
		return x.TileState
	}
	return nil
}

type PlayerState struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	X             float32                `protobuf:"fixed32,2,opt,name=x,proto3" json:"x,omitempty"`
	Y             float32                `protobuf:"fixed32,3,opt,name=y,proto3" json:"y,omitempty"`
	Health        int32                  `protobuf:"varint,4,opt,name=health,proto3" json:"health,omitempty"`
	MagicType     int32                  `protobuf:"varint,5,opt,name=magicType,proto3" json:"magicType,omitempty"`
	Angle         float32                `protobuf:"fixed32,6,opt,name=angle,proto3" json:"angle,omitempty"`
	DashCoolTime  int32                  `protobuf:"varint,7,opt,name=dashCoolTime,proto3" json:"dashCoolTime,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PlayerState) Reset() {
	*x = PlayerState{}
	mi := &file_message_message_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PlayerState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerState) ProtoMessage() {}

func (x *PlayerState) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerState.ProtoReflect.Descriptor instead.
func (*PlayerState) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{4}
}

func (x *PlayerState) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PlayerState) GetX() float32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *PlayerState) GetY() float32 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (x *PlayerState) GetHealth() int32 {
	if x != nil {
		return x.Health
	}
	return 0
}

func (x *PlayerState) GetMagicType() int32 {
	if x != nil {
		return x.MagicType
	}
	return 0
}

func (x *PlayerState) GetAngle() float32 {
	if x != nil {
		return x.Angle
	}
	return 0
}

func (x *PlayerState) GetDashCoolTime() int32 {
	if x != nil {
		return x.DashCoolTime
	}
	return 0
}

type CreateBullet struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Angle         float32                `protobuf:"fixed32,1,opt,name=angle,proto3" json:"angle,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateBullet) Reset() {
	*x = CreateBullet{}
	mi := &file_message_message_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateBullet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateBullet) ProtoMessage() {}

func (x *CreateBullet) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateBullet.ProtoReflect.Descriptor instead.
func (*CreateBullet) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{5}
}

func (x *CreateBullet) GetAngle() float32 {
	if x != nil {
		return x.Angle
	}
	return 0
}

type BulletState struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BulletId      string                 `protobuf:"bytes,1,opt,name=bulletId,proto3" json:"bulletId,omitempty"`
	BulletType    int32                  `protobuf:"varint,2,opt,name=bulletType,proto3" json:"bulletType,omitempty"`
	X             float32                `protobuf:"fixed32,3,opt,name=x,proto3" json:"x,omitempty"`
	Y             float32                `protobuf:"fixed32,4,opt,name=y,proto3" json:"y,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BulletState) Reset() {
	*x = BulletState{}
	mi := &file_message_message_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BulletState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BulletState) ProtoMessage() {}

func (x *BulletState) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BulletState.ProtoReflect.Descriptor instead.
func (*BulletState) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{6}
}

func (x *BulletState) GetBulletId() string {
	if x != nil {
		return x.BulletId
	}
	return ""
}

func (x *BulletState) GetBulletType() int32 {
	if x != nil {
		return x.BulletType
	}
	return 0
}

func (x *BulletState) GetX() float32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *BulletState) GetY() float32 {
	if x != nil {
		return x.Y
	}
	return 0
}

type GameCount struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Count         int32                  `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GameCount) Reset() {
	*x = GameCount{}
	mi := &file_message_message_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GameCount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameCount) ProtoMessage() {}

func (x *GameCount) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameCount.ProtoReflect.Descriptor instead.
func (*GameCount) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{7}
}

func (x *GameCount) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type DoDash struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Dash          bool                   `protobuf:"varint,1,opt,name=dash,proto3" json:"dash,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DoDash) Reset() {
	*x = DoDash{}
	mi := &file_message_message_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DoDash) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DoDash) ProtoMessage() {}

func (x *DoDash) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DoDash.ProtoReflect.Descriptor instead.
func (*DoDash) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{8}
}

func (x *DoDash) GetDash() bool {
	if x != nil {
		return x.Dash
	}
	return false
}

type GameStart struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MapLength     int32                  `protobuf:"varint,1,opt,name=mapLength,proto3" json:"mapLength,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GameStart) Reset() {
	*x = GameStart{}
	mi := &file_message_message_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GameStart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameStart) ProtoMessage() {}

func (x *GameStart) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameStart.ProtoReflect.Descriptor instead.
func (*GameStart) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{9}
}

func (x *GameStart) GetMapLength() int32 {
	if x != nil {
		return x.MapLength
	}
	return 0
}

type HealPackState struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ItemId        string                 `protobuf:"bytes,1,opt,name=itemId,proto3" json:"itemId,omitempty"`
	X             float32                `protobuf:"fixed32,2,opt,name=x,proto3" json:"x,omitempty"`
	Y             float32                `protobuf:"fixed32,3,opt,name=y,proto3" json:"y,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HealPackState) Reset() {
	*x = HealPackState{}
	mi := &file_message_message_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealPackState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealPackState) ProtoMessage() {}

func (x *HealPackState) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealPackState.ProtoReflect.Descriptor instead.
func (*HealPackState) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{10}
}

func (x *HealPackState) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

func (x *HealPackState) GetX() float32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *HealPackState) GetY() float32 {
	if x != nil {
		return x.Y
	}
	return 0
}

type MagicItemState struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ItemId        string                 `protobuf:"bytes,1,opt,name=itemId,proto3" json:"itemId,omitempty"`
	MagicType     int32                  `protobuf:"varint,2,opt,name=magicType,proto3" json:"magicType,omitempty"`
	X             float32                `protobuf:"fixed32,3,opt,name=x,proto3" json:"x,omitempty"`
	Y             float32                `protobuf:"fixed32,4,opt,name=y,proto3" json:"y,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MagicItemState) Reset() {
	*x = MagicItemState{}
	mi := &file_message_message_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MagicItemState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MagicItemState) ProtoMessage() {}

func (x *MagicItemState) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MagicItemState.ProtoReflect.Descriptor instead.
func (*MagicItemState) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{11}
}

func (x *MagicItemState) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

func (x *MagicItemState) GetMagicType() int32 {
	if x != nil {
		return x.MagicType
	}
	return 0
}

func (x *MagicItemState) GetX() float32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *MagicItemState) GetY() float32 {
	if x != nil {
		return x.Y
	}
	return 0
}

type PlayerDeadState struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	KillerId      string                 `protobuf:"bytes,1,opt,name=killerId,proto3" json:"killerId,omitempty"`
	DeadId        string                 `protobuf:"bytes,2,opt,name=deadId,proto3" json:"deadId,omitempty"`
	DyingStatus   int32                  `protobuf:"varint,3,opt,name=dyingStatus,proto3" json:"dyingStatus,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PlayerDeadState) Reset() {
	*x = PlayerDeadState{}
	mi := &file_message_message_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PlayerDeadState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerDeadState) ProtoMessage() {}

func (x *PlayerDeadState) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerDeadState.ProtoReflect.Descriptor instead.
func (*PlayerDeadState) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{12}
}

func (x *PlayerDeadState) GetKillerId() string {
	if x != nil {
		return x.KillerId
	}
	return ""
}

func (x *PlayerDeadState) GetDeadId() string {
	if x != nil {
		return x.DeadId
	}
	return ""
}

func (x *PlayerDeadState) GetDyingStatus() int32 {
	if x != nil {
		return x.DyingStatus
	}
	return 0
}

type TileState struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TileId        int32                  `protobuf:"varint,1,opt,name=tileId,proto3" json:"tileId,omitempty"`
	X             float32                `protobuf:"fixed32,2,opt,name=x,proto3" json:"x,omitempty"`
	Y             float32                `protobuf:"fixed32,3,opt,name=y,proto3" json:"y,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TileState) Reset() {
	*x = TileState{}
	mi := &file_message_message_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TileState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TileState) ProtoMessage() {}

func (x *TileState) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TileState.ProtoReflect.Descriptor instead.
func (*TileState) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{13}
}

func (x *TileState) GetTileId() int32 {
	if x != nil {
		return x.TileId
	}
	return 0
}

func (x *TileState) GetX() float32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *TileState) GetY() float32 {
	if x != nil {
		return x.Y
	}
	return 0
}

var File_message_message_proto protoreflect.FileDescriptor

var file_message_message_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0x81, 0x03, 0x0a, 0x07, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x32, 0x0a, 0x09,
	0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x44, 0x69, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x44, 0x69, 0x72, 0x48, 0x00, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x44, 0x69, 0x72,
	0x12, 0x29, 0x0a, 0x06, 0x64, 0x6f, 0x44, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x44, 0x6f, 0x44, 0x61, 0x73,
	0x68, 0x48, 0x00, 0x52, 0x06, 0x64, 0x6f, 0x44, 0x61, 0x73, 0x68, 0x12, 0x3b, 0x0a, 0x0c, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x42, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x48, 0x00, 0x52, 0x0c, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x42, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x12, 0x32, 0x0a, 0x09, 0x67, 0x61, 0x6d, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48,
	0x00, 0x52, 0x09, 0x67, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x32, 0x0a, 0x09,
	0x67, 0x61, 0x6d, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x09, 0x67, 0x61, 0x6d, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x2f, 0x0a, 0x08, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x6e, 0x69, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x47, 0x61, 0x6d,
	0x65, 0x49, 0x6e, 0x69, 0x74, 0x48, 0x00, 0x52, 0x08, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x6e, 0x69,
	0x74, 0x12, 0x32, 0x0a, 0x09, 0x67, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x72, 0x74, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x47,
	0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x72, 0x74, 0x48, 0x00, 0x52, 0x09, 0x67, 0x61, 0x6d, 0x65,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x22, 0x1a, 0x0a, 0x08, 0x47, 0x61, 0x6d, 0x65, 0x49, 0x6e, 0x69, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x3b, 0x0a, 0x09, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x44, 0x69, 0x72, 0x12, 0x14, 0x0a,
	0x05, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x61, 0x6e,
	0x67, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x69, 0x73, 0x4d, 0x6f, 0x76, 0x65, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x4d, 0x6f, 0x76, 0x65, 0x64, 0x22, 0xf0, 0x02,
	0x0a, 0x09, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x36, 0x0a, 0x0b, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x0b, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x12, 0x36, 0x0a, 0x0b, 0x62, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x42, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x0b,
	0x62, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x3c, 0x0a, 0x0d, 0x68,
	0x65, 0x61, 0x6c, 0x50, 0x61, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x48, 0x65, 0x61,
	0x6c, 0x50, 0x61, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x0d, 0x68, 0x65, 0x61, 0x6c,
	0x50, 0x61, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x3f, 0x0a, 0x0e, 0x6d, 0x61, 0x67,
	0x69, 0x63, 0x49, 0x74, 0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x4d, 0x61, 0x67, 0x69,
	0x63, 0x49, 0x74, 0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x0e, 0x6d, 0x61, 0x67, 0x69,
	0x63, 0x49, 0x74, 0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x42, 0x0a, 0x0f, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x44, 0x65, 0x61, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x44, 0x65, 0x61, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x0f, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x44, 0x65, 0x61, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x30,
	0x0a, 0x09, 0x74, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x54, 0x69, 0x6c, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x09, 0x74, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x22, 0xa9, 0x01, 0x0a, 0x0b, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x78, 0x12, 0x0c,
	0x0a, 0x01, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x79, 0x12, 0x16, 0x0a, 0x06,
	0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x68, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x61, 0x67, 0x69, 0x63, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6d, 0x61, 0x67, 0x69, 0x63, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x02, 0x52, 0x05, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x64, 0x61, 0x73, 0x68,
	0x43, 0x6f, 0x6f, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c,
	0x64, 0x61, 0x73, 0x68, 0x43, 0x6f, 0x6f, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x24, 0x0a, 0x0c,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x61, 0x6e, 0x67, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x61, 0x6e, 0x67,
	0x6c, 0x65, 0x22, 0x65, 0x0a, 0x0b, 0x42, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x49, 0x64, 0x12, 0x1e, 0x0a,
	0x0a, 0x62, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0a, 0x62, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0c, 0x0a,
	0x01, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x79, 0x22, 0x21, 0x0a, 0x09, 0x47, 0x61, 0x6d,
	0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x1c, 0x0a, 0x06,
	0x44, 0x6f, 0x44, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x73, 0x68, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x64, 0x61, 0x73, 0x68, 0x22, 0x29, 0x0a, 0x09, 0x47, 0x61,
	0x6d, 0x65, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x61, 0x70, 0x4c, 0x65,
	0x6e, 0x67, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6d, 0x61, 0x70, 0x4c,
	0x65, 0x6e, 0x67, 0x74, 0x68, 0x22, 0x43, 0x0a, 0x0d, 0x48, 0x65, 0x61, 0x6c, 0x50, 0x61, 0x63,
	0x6b, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x0c,
	0x0a, 0x01, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01,
	0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x79, 0x22, 0x62, 0x0a, 0x0e, 0x4d, 0x61,
	0x67, 0x69, 0x63, 0x49, 0x74, 0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x74,
	0x65, 0x6d, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x61, 0x67, 0x69, 0x63, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6d, 0x61, 0x67, 0x69, 0x63, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x78,
	0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x79, 0x22, 0x67,
	0x0a, 0x0f, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x44, 0x65, 0x61, 0x64, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x69, 0x6c, 0x6c, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x69, 0x6c, 0x6c, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x64, 0x65, 0x61, 0x64, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64,
	0x65, 0x61, 0x64, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x79, 0x69, 0x6e, 0x67, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x64, 0x79, 0x69, 0x6e,
	0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x3f, 0x0a, 0x09, 0x54, 0x69, 0x6c, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x74, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x0c, 0x0a, 0x01,
	0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x79, 0x42, 0x0a, 0x5a, 0x08, 0x2f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_message_proto_rawDescOnce sync.Once
	file_message_message_proto_rawDescData = file_message_message_proto_rawDesc
)

func file_message_message_proto_rawDescGZIP() []byte {
	file_message_message_proto_rawDescOnce.Do(func() {
		file_message_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_message_proto_rawDescData)
	})
	return file_message_message_proto_rawDescData
}

var file_message_message_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_message_message_proto_goTypes = []any{
	(*Wrapper)(nil),         // 0: message.Wrapper
	(*GameInit)(nil),        // 1: message.GameInit
	(*ChangeDir)(nil),       // 2: message.ChangeDir
	(*GameState)(nil),       // 3: message.GameState
	(*PlayerState)(nil),     // 4: message.PlayerState
	(*CreateBullet)(nil),    // 5: message.CreateBullet
	(*BulletState)(nil),     // 6: message.BulletState
	(*GameCount)(nil),       // 7: message.GameCount
	(*DoDash)(nil),          // 8: message.DoDash
	(*GameStart)(nil),       // 9: message.GameStart
	(*HealPackState)(nil),   // 10: message.HealPackState
	(*MagicItemState)(nil),  // 11: message.MagicItemState
	(*PlayerDeadState)(nil), // 12: message.PlayerDeadState
	(*TileState)(nil),       // 13: message.TileState
}
var file_message_message_proto_depIdxs = []int32{
	2,  // 0: message.Wrapper.changeDir:type_name -> message.ChangeDir
	8,  // 1: message.Wrapper.doDash:type_name -> message.DoDash
	5,  // 2: message.Wrapper.createBullet:type_name -> message.CreateBullet
	3,  // 3: message.Wrapper.gameState:type_name -> message.GameState
	7,  // 4: message.Wrapper.gameCount:type_name -> message.GameCount
	1,  // 5: message.Wrapper.gameInit:type_name -> message.GameInit
	9,  // 6: message.Wrapper.gameStart:type_name -> message.GameStart
	4,  // 7: message.GameState.playerState:type_name -> message.PlayerState
	6,  // 8: message.GameState.bulletState:type_name -> message.BulletState
	10, // 9: message.GameState.healPackState:type_name -> message.HealPackState
	11, // 10: message.GameState.magicItemState:type_name -> message.MagicItemState
	12, // 11: message.GameState.playerDeadState:type_name -> message.PlayerDeadState
	13, // 12: message.GameState.tileState:type_name -> message.TileState
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_message_message_proto_init() }
func file_message_message_proto_init() {
	if File_message_message_proto != nil {
		return
	}
	file_message_message_proto_msgTypes[0].OneofWrappers = []any{
		(*Wrapper_ChangeDir)(nil),
		(*Wrapper_DoDash)(nil),
		(*Wrapper_CreateBullet)(nil),
		(*Wrapper_GameState)(nil),
		(*Wrapper_GameCount)(nil),
		(*Wrapper_GameInit)(nil),
		(*Wrapper_GameStart)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_message_proto_goTypes,
		DependencyIndexes: file_message_message_proto_depIdxs,
		MessageInfos:      file_message_message_proto_msgTypes,
	}.Build()
	File_message_message_proto = out.File
	file_message_message_proto_rawDesc = nil
	file_message_message_proto_goTypes = nil
	file_message_message_proto_depIdxs = nil
}

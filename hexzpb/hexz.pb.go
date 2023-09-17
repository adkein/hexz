// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: hexzpb/hexz.proto

package hexzpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Board_GameState int32

const (
	Board_UNSPECIFIED Board_GameState = 0
	Board_INITIAL     Board_GameState = 1
	Board_RUNNING     Board_GameState = 2
	Board_FINISHED    Board_GameState = 3
)

// Enum value maps for Board_GameState.
var (
	Board_GameState_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "INITIAL",
		2: "RUNNING",
		3: "FINISHED",
	}
	Board_GameState_value = map[string]int32{
		"UNSPECIFIED": 0,
		"INITIAL":     1,
		"RUNNING":     2,
		"FINISHED":    3,
	}
)

func (x Board_GameState) Enum() *Board_GameState {
	p := new(Board_GameState)
	*p = x
	return p
}

func (x Board_GameState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Board_GameState) Descriptor() protoreflect.EnumDescriptor {
	return file_hexzpb_hexz_proto_enumTypes[0].Descriptor()
}

func (Board_GameState) Type() protoreflect.EnumType {
	return &file_hexzpb_hexz_proto_enumTypes[0]
}

func (x Board_GameState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Board_GameState.Descriptor instead.
func (Board_GameState) EnumDescriptor() ([]byte, []int) {
	return file_hexzpb_hexz_proto_rawDescGZIP(), []int{0, 0}
}

type Field_CellType int32

const (
	Field_UNSPECIFIED Field_CellType = 0
	Field_NORMAL      Field_CellType = 1
	Field_DEAD        Field_CellType = 2
	Field_GRASS       Field_CellType = 3
	Field_ROCK        Field_CellType = 4
	Field_FIRE        Field_CellType = 5
	Field_FLAG        Field_CellType = 6
	Field_PEST        Field_CellType = 7
	Field_DEATH       Field_CellType = 8
)

// Enum value maps for Field_CellType.
var (
	Field_CellType_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "NORMAL",
		2: "DEAD",
		3: "GRASS",
		4: "ROCK",
		5: "FIRE",
		6: "FLAG",
		7: "PEST",
		8: "DEATH",
	}
	Field_CellType_value = map[string]int32{
		"UNSPECIFIED": 0,
		"NORMAL":      1,
		"DEAD":        2,
		"GRASS":       3,
		"ROCK":        4,
		"FIRE":        5,
		"FLAG":        6,
		"PEST":        7,
		"DEATH":       8,
	}
)

func (x Field_CellType) Enum() *Field_CellType {
	p := new(Field_CellType)
	*p = x
	return p
}

func (x Field_CellType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Field_CellType) Descriptor() protoreflect.EnumDescriptor {
	return file_hexzpb_hexz_proto_enumTypes[1].Descriptor()
}

func (Field_CellType) Type() protoreflect.EnumType {
	return &file_hexzpb_hexz_proto_enumTypes[1]
}

func (x Field_CellType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Field_CellType.Descriptor instead.
func (Field_CellType) EnumDescriptor() ([]byte, []int) {
	return file_hexzpb_hexz_proto_rawDescGZIP(), []int{1, 0}
}

type Board struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Turn         int32           `protobuf:"varint,1,opt,name=turn,proto3" json:"turn,omitempty"`
	Move         int32           `protobuf:"varint,2,opt,name=move,proto3" json:"move,omitempty"`
	LastRevealed int32           `protobuf:"varint,3,opt,name=last_revealed,json=lastRevealed,proto3" json:"last_revealed,omitempty"`
	FlatFields   []*Field        `protobuf:"bytes,4,rep,name=flat_fields,json=flatFields,proto3" json:"flat_fields,omitempty"`
	Score        []int32         `protobuf:"varint,5,rep,packed,name=score,proto3" json:"score,omitempty"`
	Resources    []*ResourceInfo `protobuf:"bytes,6,rep,name=resources,proto3" json:"resources,omitempty"`
	State        Board_GameState `protobuf:"varint,7,opt,name=state,proto3,enum=github.com.dnswlt.hexz.Board_GameState" json:"state,omitempty"`
}

func (x *Board) Reset() {
	*x = Board{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hexzpb_hexz_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Board) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Board) ProtoMessage() {}

func (x *Board) ProtoReflect() protoreflect.Message {
	mi := &file_hexzpb_hexz_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Board.ProtoReflect.Descriptor instead.
func (*Board) Descriptor() ([]byte, []int) {
	return file_hexzpb_hexz_proto_rawDescGZIP(), []int{0}
}

func (x *Board) GetTurn() int32 {
	if x != nil {
		return x.Turn
	}
	return 0
}

func (x *Board) GetMove() int32 {
	if x != nil {
		return x.Move
	}
	return 0
}

func (x *Board) GetLastRevealed() int32 {
	if x != nil {
		return x.LastRevealed
	}
	return 0
}

func (x *Board) GetFlatFields() []*Field {
	if x != nil {
		return x.FlatFields
	}
	return nil
}

func (x *Board) GetScore() []int32 {
	if x != nil {
		return x.Score
	}
	return nil
}

func (x *Board) GetResources() []*ResourceInfo {
	if x != nil {
		return x.Resources
	}
	return nil
}

func (x *Board) GetState() Board_GameState {
	if x != nil {
		return x.State
	}
	return Board_UNSPECIFIED
}

type Field struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type     Field_CellType `protobuf:"varint,1,opt,name=type,proto3,enum=github.com.dnswlt.hexz.Field_CellType" json:"type,omitempty"`
	Owner    int32          `protobuf:"varint,2,opt,name=owner,proto3" json:"owner,omitempty"`
	Hidden   bool           `protobuf:"varint,3,opt,name=hidden,proto3" json:"hidden,omitempty"`
	Value    int32          `protobuf:"varint,4,opt,name=value,proto3" json:"value,omitempty"`
	Blocked  int32          `protobuf:"varint,5,opt,name=blocked,proto3" json:"blocked,omitempty"`
	Lifetime int32          `protobuf:"varint,6,opt,name=lifetime,proto3" json:"lifetime,omitempty"`
	NextVal  []int32        `protobuf:"varint,7,rep,packed,name=next_val,json=nextVal,proto3" json:"next_val,omitempty"`
}

func (x *Field) Reset() {
	*x = Field{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hexzpb_hexz_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Field) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Field) ProtoMessage() {}

func (x *Field) ProtoReflect() protoreflect.Message {
	mi := &file_hexzpb_hexz_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Field.ProtoReflect.Descriptor instead.
func (*Field) Descriptor() ([]byte, []int) {
	return file_hexzpb_hexz_proto_rawDescGZIP(), []int{1}
}

func (x *Field) GetType() Field_CellType {
	if x != nil {
		return x.Type
	}
	return Field_UNSPECIFIED
}

func (x *Field) GetOwner() int32 {
	if x != nil {
		return x.Owner
	}
	return 0
}

func (x *Field) GetHidden() bool {
	if x != nil {
		return x.Hidden
	}
	return false
}

func (x *Field) GetValue() int32 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Field) GetBlocked() int32 {
	if x != nil {
		return x.Blocked
	}
	return 0
}

func (x *Field) GetLifetime() int32 {
	if x != nil {
		return x.Lifetime
	}
	return 0
}

func (x *Field) GetNextVal() []int32 {
	if x != nil {
		return x.NextVal
	}
	return nil
}

type ResourceInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Has exactly as many elements as there are cell types.
	NumPieces []int32 `protobuf:"varint,1,rep,packed,name=num_pieces,json=numPieces,proto3" json:"num_pieces,omitempty"`
}

func (x *ResourceInfo) Reset() {
	*x = ResourceInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hexzpb_hexz_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceInfo) ProtoMessage() {}

func (x *ResourceInfo) ProtoReflect() protoreflect.Message {
	mi := &file_hexzpb_hexz_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceInfo.ProtoReflect.Descriptor instead.
func (*ResourceInfo) Descriptor() ([]byte, []int) {
	return file_hexzpb_hexz_proto_rawDescGZIP(), []int{2}
}

func (x *ResourceInfo) GetNumPieces() []int32 {
	if x != nil {
		return x.NumPieces
	}
	return nil
}

// The encoded state of a GameEngineFlagz. Used for saving and loading games.
type GameEngineFlagzState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Board       *Board            `protobuf:"bytes,1,opt,name=board,proto3" json:"board,omitempty"`
	FreeCells   int32             `protobuf:"varint,2,opt,name=free_cells,json=freeCells,proto3" json:"free_cells,omitempty"`
	NormalMoves []int32           `protobuf:"varint,3,rep,packed,name=normal_moves,json=normalMoves,proto3" json:"normal_moves,omitempty"`
	Moves       []*GameEngineMove `protobuf:"bytes,4,rep,name=moves,proto3" json:"moves,omitempty"`
}

func (x *GameEngineFlagzState) Reset() {
	*x = GameEngineFlagzState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hexzpb_hexz_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameEngineFlagzState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameEngineFlagzState) ProtoMessage() {}

func (x *GameEngineFlagzState) ProtoReflect() protoreflect.Message {
	mi := &file_hexzpb_hexz_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameEngineFlagzState.ProtoReflect.Descriptor instead.
func (*GameEngineFlagzState) Descriptor() ([]byte, []int) {
	return file_hexzpb_hexz_proto_rawDescGZIP(), []int{3}
}

func (x *GameEngineFlagzState) GetBoard() *Board {
	if x != nil {
		return x.Board
	}
	return nil
}

func (x *GameEngineFlagzState) GetFreeCells() int32 {
	if x != nil {
		return x.FreeCells
	}
	return 0
}

func (x *GameEngineFlagzState) GetNormalMoves() []int32 {
	if x != nil {
		return x.NormalMoves
	}
	return nil
}

func (x *GameEngineFlagzState) GetMoves() []*GameEngineMove {
	if x != nil {
		return x.Moves
	}
	return nil
}

type GameEngineMove struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerNum int32          `protobuf:"varint,1,opt,name=player_num,json=playerNum,proto3" json:"player_num,omitempty"`
	Move      int32          `protobuf:"varint,2,opt,name=move,proto3" json:"move,omitempty"`
	Row       int32          `protobuf:"varint,3,opt,name=row,proto3" json:"row,omitempty"`
	Col       int32          `protobuf:"varint,4,opt,name=col,proto3" json:"col,omitempty"`
	CellType  Field_CellType `protobuf:"varint,5,opt,name=cell_type,json=cellType,proto3,enum=github.com.dnswlt.hexz.Field_CellType" json:"cell_type,omitempty"`
}

func (x *GameEngineMove) Reset() {
	*x = GameEngineMove{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hexzpb_hexz_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameEngineMove) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameEngineMove) ProtoMessage() {}

func (x *GameEngineMove) ProtoReflect() protoreflect.Message {
	mi := &file_hexzpb_hexz_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameEngineMove.ProtoReflect.Descriptor instead.
func (*GameEngineMove) Descriptor() ([]byte, []int) {
	return file_hexzpb_hexz_proto_rawDescGZIP(), []int{4}
}

func (x *GameEngineMove) GetPlayerNum() int32 {
	if x != nil {
		return x.PlayerNum
	}
	return 0
}

func (x *GameEngineMove) GetMove() int32 {
	if x != nil {
		return x.Move
	}
	return 0
}

func (x *GameEngineMove) GetRow() int32 {
	if x != nil {
		return x.Row
	}
	return 0
}

func (x *GameEngineMove) GetCol() int32 {
	if x != nil {
		return x.Col
	}
	return 0
}

func (x *GameEngineMove) GetCellType() Field_CellType {
	if x != nil {
		return x.CellType
	}
	return Field_UNSPECIFIED
}

var File_hexzpb_hexz_proto protoreflect.FileDescriptor

var file_hexzpb_hexz_proto_rawDesc = []byte{
	0x0a, 0x11, 0x68, 0x65, 0x78, 0x7a, 0x70, 0x62, 0x2f, 0x68, 0x65, 0x78, 0x7a, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x16, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e,
	0x64, 0x6e, 0x73, 0x77, 0x6c, 0x74, 0x2e, 0x68, 0x65, 0x78, 0x7a, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf3, 0x02, 0x0a,
	0x05, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x75, 0x72, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x75, 0x72, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x6f,
	0x76, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x6d, 0x6f, 0x76, 0x65, 0x12, 0x23,
	0x0a, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x76, 0x65, 0x61, 0x6c, 0x65, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x52, 0x65, 0x76, 0x65, 0x61,
	0x6c, 0x65, 0x64, 0x12, 0x3e, 0x0a, 0x0b, 0x66, 0x6c, 0x61, 0x74, 0x5f, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x6e, 0x73, 0x77, 0x6c, 0x74, 0x2e, 0x68, 0x65, 0x78,
	0x7a, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x0a, 0x66, 0x6c, 0x61, 0x74, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x05, 0x20, 0x03,
	0x28, 0x05, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x42, 0x0a, 0x09, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x6e, 0x73, 0x77, 0x6c, 0x74,
	0x2e, 0x68, 0x65, 0x78, 0x7a, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x3d, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x27, 0x2e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x6e, 0x73, 0x77, 0x6c, 0x74,
	0x2e, 0x68, 0x65, 0x78, 0x7a, 0x2e, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x47, 0x61, 0x6d, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x44, 0x0a, 0x09,
	0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x49, 0x4e,
	0x49, 0x54, 0x49, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x55, 0x4e, 0x4e, 0x49,
	0x4e, 0x47, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x46, 0x49, 0x4e, 0x49, 0x53, 0x48, 0x45, 0x44,
	0x10, 0x03, 0x22, 0xc9, 0x02, 0x0a, 0x05, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x3a, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x6e, 0x73, 0x77, 0x6c, 0x74, 0x2e, 0x68,
	0x65, 0x78, 0x7a, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x43, 0x65, 0x6c, 0x6c, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x69, 0x66, 0x65, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6c, 0x69, 0x66, 0x65, 0x74, 0x69,
	0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x18, 0x07,
	0x20, 0x03, 0x28, 0x05, 0x52, 0x07, 0x6e, 0x65, 0x78, 0x74, 0x56, 0x61, 0x6c, 0x22, 0x6f, 0x0a,
	0x08, 0x43, 0x65, 0x6c, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x4f,
	0x52, 0x4d, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x45, 0x41, 0x44, 0x10, 0x02,
	0x12, 0x09, 0x0a, 0x05, 0x47, 0x52, 0x41, 0x53, 0x53, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04, 0x52,
	0x4f, 0x43, 0x4b, 0x10, 0x04, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x49, 0x52, 0x45, 0x10, 0x05, 0x12,
	0x08, 0x0a, 0x04, 0x46, 0x4c, 0x41, 0x47, 0x10, 0x06, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x45, 0x53,
	0x54, 0x10, 0x07, 0x12, 0x09, 0x0a, 0x05, 0x44, 0x45, 0x41, 0x54, 0x48, 0x10, 0x08, 0x22, 0x2d,
	0x0a, 0x0c, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1d,
	0x0a, 0x0a, 0x6e, 0x75, 0x6d, 0x5f, 0x70, 0x69, 0x65, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x05, 0x52, 0x09, 0x6e, 0x75, 0x6d, 0x50, 0x69, 0x65, 0x63, 0x65, 0x73, 0x22, 0xcb, 0x01,
	0x0a, 0x14, 0x47, 0x61, 0x6d, 0x65, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x46, 0x6c, 0x61, 0x67,
	0x7a, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x33, 0x0a, 0x05, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x64, 0x6e, 0x73, 0x77, 0x6c, 0x74, 0x2e, 0x68, 0x65, 0x78, 0x7a, 0x2e, 0x42,
	0x6f, 0x61, 0x72, 0x64, 0x52, 0x05, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x66,
	0x72, 0x65, 0x65, 0x5f, 0x63, 0x65, 0x6c, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x09, 0x66, 0x72, 0x65, 0x65, 0x43, 0x65, 0x6c, 0x6c, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x6f,
	0x72, 0x6d, 0x61, 0x6c, 0x5f, 0x6d, 0x6f, 0x76, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x05,
	0x52, 0x0b, 0x6e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x4d, 0x6f, 0x76, 0x65, 0x73, 0x12, 0x3c, 0x0a,
	0x05, 0x6d, 0x6f, 0x76, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x6e, 0x73, 0x77, 0x6c, 0x74,
	0x2e, 0x68, 0x65, 0x78, 0x7a, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65,
	0x4d, 0x6f, 0x76, 0x65, 0x52, 0x05, 0x6d, 0x6f, 0x76, 0x65, 0x73, 0x22, 0xac, 0x01, 0x0a, 0x0e,
	0x47, 0x61, 0x6d, 0x65, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x4d, 0x6f, 0x76, 0x65, 0x12, 0x1d,
	0x0a, 0x0a, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4e, 0x75, 0x6d, 0x12, 0x12, 0x0a,
	0x04, 0x6d, 0x6f, 0x76, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x6d, 0x6f, 0x76,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x6f, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x72, 0x6f, 0x77, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x6f, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x03, 0x63, 0x6f, 0x6c, 0x12, 0x43, 0x0a, 0x09, 0x63, 0x65, 0x6c, 0x6c, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x6e, 0x73, 0x77, 0x6c, 0x74, 0x2e, 0x68, 0x65, 0x78,
	0x7a, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x43, 0x65, 0x6c, 0x6c, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x08, 0x63, 0x65, 0x6c, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x42, 0x1f, 0x5a, 0x1d, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6e, 0x73, 0x77, 0x6c, 0x74, 0x2f,
	0x68, 0x65, 0x78, 0x7a, 0x2f, 0x68, 0x65, 0x78, 0x7a, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_hexzpb_hexz_proto_rawDescOnce sync.Once
	file_hexzpb_hexz_proto_rawDescData = file_hexzpb_hexz_proto_rawDesc
)

func file_hexzpb_hexz_proto_rawDescGZIP() []byte {
	file_hexzpb_hexz_proto_rawDescOnce.Do(func() {
		file_hexzpb_hexz_proto_rawDescData = protoimpl.X.CompressGZIP(file_hexzpb_hexz_proto_rawDescData)
	})
	return file_hexzpb_hexz_proto_rawDescData
}

var file_hexzpb_hexz_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_hexzpb_hexz_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_hexzpb_hexz_proto_goTypes = []interface{}{
	(Board_GameState)(0),         // 0: github.com.dnswlt.hexz.Board.GameState
	(Field_CellType)(0),          // 1: github.com.dnswlt.hexz.Field.CellType
	(*Board)(nil),                // 2: github.com.dnswlt.hexz.Board
	(*Field)(nil),                // 3: github.com.dnswlt.hexz.Field
	(*ResourceInfo)(nil),         // 4: github.com.dnswlt.hexz.ResourceInfo
	(*GameEngineFlagzState)(nil), // 5: github.com.dnswlt.hexz.GameEngineFlagzState
	(*GameEngineMove)(nil),       // 6: github.com.dnswlt.hexz.GameEngineMove
}
var file_hexzpb_hexz_proto_depIdxs = []int32{
	3, // 0: github.com.dnswlt.hexz.Board.flat_fields:type_name -> github.com.dnswlt.hexz.Field
	4, // 1: github.com.dnswlt.hexz.Board.resources:type_name -> github.com.dnswlt.hexz.ResourceInfo
	0, // 2: github.com.dnswlt.hexz.Board.state:type_name -> github.com.dnswlt.hexz.Board.GameState
	1, // 3: github.com.dnswlt.hexz.Field.type:type_name -> github.com.dnswlt.hexz.Field.CellType
	2, // 4: github.com.dnswlt.hexz.GameEngineFlagzState.board:type_name -> github.com.dnswlt.hexz.Board
	6, // 5: github.com.dnswlt.hexz.GameEngineFlagzState.moves:type_name -> github.com.dnswlt.hexz.GameEngineMove
	1, // 6: github.com.dnswlt.hexz.GameEngineMove.cell_type:type_name -> github.com.dnswlt.hexz.Field.CellType
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_hexzpb_hexz_proto_init() }
func file_hexzpb_hexz_proto_init() {
	if File_hexzpb_hexz_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_hexzpb_hexz_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Board); i {
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
		file_hexzpb_hexz_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Field); i {
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
		file_hexzpb_hexz_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceInfo); i {
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
		file_hexzpb_hexz_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GameEngineFlagzState); i {
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
		file_hexzpb_hexz_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GameEngineMove); i {
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
			RawDescriptor: file_hexzpb_hexz_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_hexzpb_hexz_proto_goTypes,
		DependencyIndexes: file_hexzpb_hexz_proto_depIdxs,
		EnumInfos:         file_hexzpb_hexz_proto_enumTypes,
		MessageInfos:      file_hexzpb_hexz_proto_msgTypes,
	}.Build()
	File_hexzpb_hexz_proto = out.File
	file_hexzpb_hexz_proto_rawDesc = nil
	file_hexzpb_hexz_proto_goTypes = nil
	file_hexzpb_hexz_proto_depIdxs = nil
}

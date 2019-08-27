// Code generated by protoc-gen-go. DO NOT EDIT.
// source: plugin.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// TextFormat ...
type Message_TextFormat int32

const (
	Message_PLAIN_TEXT Message_TextFormat = 0
)

var Message_TextFormat_name = map[int32]string{
	0: "PLAIN_TEXT",
}
var Message_TextFormat_value = map[string]int32{
	"PLAIN_TEXT": 0,
}

func (x Message_TextFormat) String() string {
	return proto.EnumName(Message_TextFormat_name, int32(x))
}
func (Message_TextFormat) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{4, 0}
}

// Errors to occur in the communication
type Error_Code int32

const (
	Error_UNKNOWN  Error_Code = 0
	Error_REGISTER Error_Code = 1
)

var Error_Code_name = map[int32]string{
	0: "UNKNOWN",
	1: "REGISTER",
}
var Error_Code_value = map[string]int32{
	"UNKNOWN":  0,
	"REGISTER": 1,
}

func (x Error_Code) String() string {
	return proto.EnumName(Error_Code_name, int32(x))
}
func (Error_Code) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{5, 0}
}

// Event ...
//
// These are events that may occur.
type Event struct {
	// Event ...
	//
	// Types that are valid to be assigned to Event:
	//	*Event_Message
	//	*Event_Reply
	//	*Event_Private
	//	*Event_Config
	//	*Event_Register
	//	*Event_Error
	Event                isEvent_Event `protobuf_oneof:"event"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{0}
}
func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (dst *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(dst, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

type isEvent_Event interface {
	isEvent_Event()
}

type Event_Message struct {
	Message *Message `protobuf:"bytes,10,opt,name=message,proto3,oneof"`
}

type Event_Reply struct {
	Reply *Message `protobuf:"bytes,11,opt,name=reply,proto3,oneof"`
}

type Event_Private struct {
	Private *Message `protobuf:"bytes,12,opt,name=private,proto3,oneof"`
}

type Event_Config struct {
	Config *Config `protobuf:"bytes,15,opt,name=config,proto3,oneof"`
}

type Event_Register struct {
	Register *Register `protobuf:"bytes,16,opt,name=register,proto3,oneof"`
}

type Event_Error struct {
	Error *Error `protobuf:"bytes,20,opt,name=error,proto3,oneof"`
}

func (*Event_Message) isEvent_Event() {}

func (*Event_Reply) isEvent_Event() {}

func (*Event_Private) isEvent_Event() {}

func (*Event_Config) isEvent_Event() {}

func (*Event_Register) isEvent_Event() {}

func (*Event_Error) isEvent_Event() {}

func (m *Event) GetEvent() isEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (m *Event) GetMessage() *Message {
	if x, ok := m.GetEvent().(*Event_Message); ok {
		return x.Message
	}
	return nil
}

func (m *Event) GetReply() *Message {
	if x, ok := m.GetEvent().(*Event_Reply); ok {
		return x.Reply
	}
	return nil
}

func (m *Event) GetPrivate() *Message {
	if x, ok := m.GetEvent().(*Event_Private); ok {
		return x.Private
	}
	return nil
}

func (m *Event) GetConfig() *Config {
	if x, ok := m.GetEvent().(*Event_Config); ok {
		return x.Config
	}
	return nil
}

func (m *Event) GetRegister() *Register {
	if x, ok := m.GetEvent().(*Event_Register); ok {
		return x.Register
	}
	return nil
}

func (m *Event) GetError() *Error {
	if x, ok := m.GetEvent().(*Event_Error); ok {
		return x.Error
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Event) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Event_OneofMarshaler, _Event_OneofUnmarshaler, _Event_OneofSizer, []interface{}{
		(*Event_Message)(nil),
		(*Event_Reply)(nil),
		(*Event_Private)(nil),
		(*Event_Config)(nil),
		(*Event_Register)(nil),
		(*Event_Error)(nil),
	}
}

func _Event_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Event)
	// event
	switch x := m.Event.(type) {
	case *Event_Message:
		b.EncodeVarint(10<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Message); err != nil {
			return err
		}
	case *Event_Reply:
		b.EncodeVarint(11<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Reply); err != nil {
			return err
		}
	case *Event_Private:
		b.EncodeVarint(12<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Private); err != nil {
			return err
		}
	case *Event_Config:
		b.EncodeVarint(15<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Config); err != nil {
			return err
		}
	case *Event_Register:
		b.EncodeVarint(16<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Register); err != nil {
			return err
		}
	case *Event_Error:
		b.EncodeVarint(20<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Error); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Event.Event has unexpected type %T", x)
	}
	return nil
}

func _Event_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Event)
	switch tag {
	case 10: // event.message
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Message)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Message{msg}
		return true, err
	case 11: // event.reply
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Message)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Reply{msg}
		return true, err
	case 12: // event.private
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Message)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Private{msg}
		return true, err
	case 15: // event.config
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Config)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Config{msg}
		return true, err
	case 16: // event.register
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Register)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Register{msg}
		return true, err
	case 20: // event.error
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Error)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Error{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Event_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Event)
	// event
	switch x := m.Event.(type) {
	case *Event_Message:
		s := proto.Size(x.Message)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Reply:
		s := proto.Size(x.Reply)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Private:
		s := proto.Size(x.Private)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Config:
		s := proto.Size(x.Config)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Register:
		s := proto.Size(x.Register)
		n += 2 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Error:
		s := proto.Size(x.Error)
		n += 2 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Register ...
//
// Register action for the server.
type Register struct {
	Plugin               *Plugin  `protobuf:"bytes,1,opt,name=plugin,proto3" json:"plugin,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Register) Reset()         { *m = Register{} }
func (m *Register) String() string { return proto.CompactTextString(m) }
func (*Register) ProtoMessage()    {}
func (*Register) Descriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{1}
}
func (m *Register) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Register.Unmarshal(m, b)
}
func (m *Register) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Register.Marshal(b, m, deterministic)
}
func (dst *Register) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Register.Merge(dst, src)
}
func (m *Register) XXX_Size() int {
	return xxx_messageInfo_Register.Size(m)
}
func (m *Register) XXX_DiscardUnknown() {
	xxx_messageInfo_Register.DiscardUnknown(m)
}

var xxx_messageInfo_Register proto.InternalMessageInfo

func (m *Register) GetPlugin() *Plugin {
	if m != nil {
		return m.Plugin
	}
	return nil
}

// Config ...
//
// Config is a plugin config from the server
type Config struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{2}
}
func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (dst *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(dst, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

// Plugin ...
//
// These are all info for a plugin.
type Plugin struct {
	// name ...
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// identifier ...
	Identifier string `protobuf:"bytes,2,opt,name=identifier,proto3" json:"identifier,omitempty"`
	// version ...
	Version string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	// path ...
	Path                 string   `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Plugin) Reset()         { *m = Plugin{} }
func (m *Plugin) String() string { return proto.CompactTextString(m) }
func (*Plugin) ProtoMessage()    {}
func (*Plugin) Descriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{3}
}
func (m *Plugin) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Plugin.Unmarshal(m, b)
}
func (m *Plugin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Plugin.Marshal(b, m, deterministic)
}
func (dst *Plugin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Plugin.Merge(dst, src)
}
func (m *Plugin) XXX_Size() int {
	return xxx_messageInfo_Plugin.Size(m)
}
func (m *Plugin) XXX_DiscardUnknown() {
	xxx_messageInfo_Plugin.DiscardUnknown(m)
}

var xxx_messageInfo_Plugin proto.InternalMessageInfo

func (m *Plugin) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Plugin) GetIdentifier() string {
	if m != nil {
		return m.Identifier
	}
	return ""
}

func (m *Plugin) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Plugin) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

// Message
//
// Information from a chat.
type Message struct {
	// UUID ...
	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	// ID ...
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	// Type ...
	Type string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	// Channel ...
	Channel *Message_Channel `protobuf:"bytes,4,opt,name=channel,proto3" json:"channel,omitempty"`
	// From ...
	From *Message_User `protobuf:"bytes,5,opt,name=from,proto3" json:"from,omitempty"`
	// Recipient ...
	Recipient *Message_Recipient `protobuf:"bytes,6,opt,name=recipient,proto3" json:"recipient,omitempty"`
	// isBot ...
	IsBot bool `protobuf:"varint,7,opt,name=is_bot,json=isBot,proto3" json:"is_bot,omitempty"`
	// isDirectMessage ...
	IsDirectMessage bool `protobuf:"varint,8,opt,name=is_direct_message,json=isDirectMessage,proto3" json:"is_direct_message,omitempty"`
	// Timestamp ...
	Timestamp *timestamp.Timestamp `protobuf:"bytes,10,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// TextFormat ...
	TextFormat Message_TextFormat `protobuf:"varint,20,opt,name=text_format,json=textFormat,proto3,enum=proto.Message_TextFormat" json:"text_format,omitempty"`
	// Text ...
	Text                 string   `protobuf:"bytes,21,opt,name=text,proto3" json:"text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{4}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Message) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Message) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Message) GetChannel() *Message_Channel {
	if m != nil {
		return m.Channel
	}
	return nil
}

func (m *Message) GetFrom() *Message_User {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Message) GetRecipient() *Message_Recipient {
	if m != nil {
		return m.Recipient
	}
	return nil
}

func (m *Message) GetIsBot() bool {
	if m != nil {
		return m.IsBot
	}
	return false
}

func (m *Message) GetIsDirectMessage() bool {
	if m != nil {
		return m.IsDirectMessage
	}
	return false
}

func (m *Message) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Message) GetTextFormat() Message_TextFormat {
	if m != nil {
		return m.TextFormat
	}
	return Message_PLAIN_TEXT
}

func (m *Message) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

// User ...
//
// This is a user interaction ...
type Message_User struct {
	// ID ...
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name ...
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Team ...
	Team                 *Message_Team `protobuf:"bytes,3,opt,name=team,proto3" json:"team,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Message_User) Reset()         { *m = Message_User{} }
func (m *Message_User) String() string { return proto.CompactTextString(m) }
func (*Message_User) ProtoMessage()    {}
func (*Message_User) Descriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{4, 0}
}
func (m *Message_User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message_User.Unmarshal(m, b)
}
func (m *Message_User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message_User.Marshal(b, m, deterministic)
}
func (dst *Message_User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message_User.Merge(dst, src)
}
func (m *Message_User) XXX_Size() int {
	return xxx_messageInfo_Message_User.Size(m)
}
func (m *Message_User) XXX_DiscardUnknown() {
	xxx_messageInfo_Message_User.DiscardUnknown(m)
}

var xxx_messageInfo_Message_User proto.InternalMessageInfo

func (m *Message_User) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Message_User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Message_User) GetTeam() *Message_Team {
	if m != nil {
		return m.Team
	}
	return nil
}

// Team ...
//
// This is a team interaction ...
type Message_Team struct {
	// ID ...
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name ...
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message_Team) Reset()         { *m = Message_Team{} }
func (m *Message_Team) String() string { return proto.CompactTextString(m) }
func (*Message_Team) ProtoMessage()    {}
func (*Message_Team) Descriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{4, 1}
}
func (m *Message_Team) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message_Team.Unmarshal(m, b)
}
func (m *Message_Team) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message_Team.Marshal(b, m, deterministic)
}
func (dst *Message_Team) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message_Team.Merge(dst, src)
}
func (m *Message_Team) XXX_Size() int {
	return xxx_messageInfo_Message_Team.Size(m)
}
func (m *Message_Team) XXX_DiscardUnknown() {
	xxx_messageInfo_Message_Team.DiscardUnknown(m)
}

var xxx_messageInfo_Message_Team proto.InternalMessageInfo

func (m *Message_Team) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Message_Team) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Channel ...
//
// This is the channel of the interaction ...
type Message_Channel struct {
	// ID ...
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name ...
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message_Channel) Reset()         { *m = Message_Channel{} }
func (m *Message_Channel) String() string { return proto.CompactTextString(m) }
func (*Message_Channel) ProtoMessage()    {}
func (*Message_Channel) Descriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{4, 2}
}
func (m *Message_Channel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message_Channel.Unmarshal(m, b)
}
func (m *Message_Channel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message_Channel.Marshal(b, m, deterministic)
}
func (dst *Message_Channel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message_Channel.Merge(dst, src)
}
func (m *Message_Channel) XXX_Size() int {
	return xxx_messageInfo_Message_Channel.Size(m)
}
func (m *Message_Channel) XXX_DiscardUnknown() {
	xxx_messageInfo_Message_Channel.DiscardUnknown(m)
}

var xxx_messageInfo_Message_Channel proto.InternalMessageInfo

func (m *Message_Channel) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Message_Channel) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Recipient ...
//
// This is the recipient of an interaction ...
type Message_Recipient struct {
	// ID ...
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name ...
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Team ...
	Team                 *Message_Team `protobuf:"bytes,3,opt,name=team,proto3" json:"team,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Message_Recipient) Reset()         { *m = Message_Recipient{} }
func (m *Message_Recipient) String() string { return proto.CompactTextString(m) }
func (*Message_Recipient) ProtoMessage()    {}
func (*Message_Recipient) Descriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{4, 3}
}
func (m *Message_Recipient) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message_Recipient.Unmarshal(m, b)
}
func (m *Message_Recipient) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message_Recipient.Marshal(b, m, deterministic)
}
func (dst *Message_Recipient) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message_Recipient.Merge(dst, src)
}
func (m *Message_Recipient) XXX_Size() int {
	return xxx_messageInfo_Message_Recipient.Size(m)
}
func (m *Message_Recipient) XXX_DiscardUnknown() {
	xxx_messageInfo_Message_Recipient.DiscardUnknown(m)
}

var xxx_messageInfo_Message_Recipient proto.InternalMessageInfo

func (m *Message_Recipient) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Message_Recipient) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Message_Recipient) GetTeam() *Message_Team {
	if m != nil {
		return m.Team
	}
	return nil
}

// Error
//
// These are standard error codes.
type Error struct {
	Code                 Error_Code        `protobuf:"varint,1,opt,name=code,proto3,enum=proto.Error_Code" json:"code,omitempty"`
	Message              string            `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Details              map[string]string `protobuf:"bytes,3,rep,name=details,proto3" json:"details,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{5}
}
func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (dst *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(dst, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetCode() Error_Code {
	if m != nil {
		return m.Code
	}
	return Error_UNKNOWN
}

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Error) GetDetails() map[string]string {
	if m != nil {
		return m.Details
	}
	return nil
}

// Empty ...
//
// Empty reply.
type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_plugin_ac867b57e049421c, []int{6}
}
func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (dst *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(dst, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Event)(nil), "proto.Event")
	proto.RegisterType((*Register)(nil), "proto.Register")
	proto.RegisterType((*Config)(nil), "proto.Config")
	proto.RegisterType((*Plugin)(nil), "proto.Plugin")
	proto.RegisterType((*Message)(nil), "proto.Message")
	proto.RegisterType((*Message_User)(nil), "proto.Message.User")
	proto.RegisterType((*Message_Team)(nil), "proto.Message.Team")
	proto.RegisterType((*Message_Channel)(nil), "proto.Message.Channel")
	proto.RegisterType((*Message_Recipient)(nil), "proto.Message.Recipient")
	proto.RegisterType((*Error)(nil), "proto.Error")
	proto.RegisterMapType((map[string]string)(nil), "proto.Error.DetailsEntry")
	proto.RegisterType((*Empty)(nil), "proto.Empty")
	proto.RegisterEnum("proto.Message_TextFormat", Message_TextFormat_name, Message_TextFormat_value)
	proto.RegisterEnum("proto.Error_Code", Error_Code_name, Error_Code_value)
}

func init() { proto.RegisterFile("plugin.proto", fileDescriptor_plugin_ac867b57e049421c) }

var fileDescriptor_plugin_ac867b57e049421c = []byte{
	// 699 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xdb, 0x4e, 0x1b, 0x49,
	0x10, 0x86, 0x7d, 0x1a, 0x8f, 0x5d, 0xf6, 0x1a, 0xd3, 0x0b, 0xab, 0x5e, 0x6b, 0xb5, 0xcb, 0x8e,
	0x96, 0x05, 0x21, 0x61, 0x12, 0x23, 0x45, 0x88, 0xbb, 0x00, 0x4e, 0x40, 0x49, 0x1c, 0xd4, 0x18,
	0x85, 0x3b, 0x6b, 0xf0, 0x94, 0x4d, 0x27, 0x9e, 0x83, 0x7a, 0xda, 0x96, 0xfd, 0x12, 0x79, 0xaf,
	0xbc, 0x43, 0x1e, 0x26, 0xea, 0xc3, 0x18, 0x43, 0x84, 0xc4, 0x45, 0xae, 0x5c, 0x5d, 0xff, 0xd7,
	0xd5, 0xe5, 0x3a, 0x0c, 0xd4, 0x93, 0xc9, 0x74, 0xcc, 0xa3, 0x76, 0x22, 0x62, 0x19, 0x13, 0x47,
	0xff, 0xb4, 0xfe, 0x19, 0xc7, 0xf1, 0x78, 0x82, 0x07, 0xfa, 0x74, 0x3b, 0x1d, 0x1d, 0x48, 0x1e,
	0x62, 0x2a, 0xfd, 0x30, 0x31, 0x9c, 0xf7, 0xb5, 0x00, 0x4e, 0x77, 0x86, 0x91, 0x24, 0x7b, 0xe0,
	0x86, 0x98, 0xa6, 0xfe, 0x18, 0x29, 0x6c, 0xe5, 0x77, 0x6b, 0x9d, 0x86, 0x41, 0xda, 0x1f, 0x8c,
	0xf7, 0x3c, 0xc7, 0x32, 0x80, 0xfc, 0x0f, 0x8e, 0xc0, 0x64, 0xb2, 0xa0, 0xb5, 0x27, 0x48, 0x23,
	0xab, 0x98, 0x89, 0xe0, 0x33, 0x5f, 0x22, 0xad, 0x3f, 0x15, 0xd3, 0x02, 0x64, 0x07, 0xca, 0xc3,
	0x38, 0x1a, 0xf1, 0x31, 0x5d, 0xd3, 0xe8, 0x6f, 0x16, 0x3d, 0xd5, 0xce, 0xf3, 0x1c, 0xb3, 0x32,
	0xd9, 0x87, 0x8a, 0xc0, 0x31, 0x4f, 0x25, 0x0a, 0xda, 0xd4, 0xe8, 0x9a, 0x45, 0x99, 0x75, 0x9f,
	0xe7, 0xd8, 0x12, 0x21, 0xff, 0x81, 0x83, 0x42, 0xc4, 0x82, 0x6e, 0x68, 0xb6, 0x6e, 0xd9, 0xae,
	0xf2, 0xa9, 0x4c, 0xb5, 0x78, 0xe2, 0x82, 0x83, 0xaa, 0x0c, 0xde, 0x4b, 0xa8, 0x64, 0x61, 0xc8,
	0x36, 0x94, 0x4d, 0x51, 0x69, 0xfe, 0x41, 0x4a, 0x97, 0xda, 0xc9, 0xac, 0xe8, 0x55, 0xa0, 0x6c,
	0x92, 0xf4, 0x3e, 0x43, 0xd9, 0x68, 0x84, 0x40, 0x29, 0xf2, 0x43, 0xd4, 0x17, 0xab, 0x4c, 0xdb,
	0xe4, 0x6f, 0x00, 0x1e, 0x60, 0x24, 0xf9, 0x88, 0xa3, 0xa0, 0x05, 0xad, 0xac, 0x78, 0x08, 0x05,
	0x77, 0x86, 0x22, 0xe5, 0x71, 0x44, 0x8b, 0x5a, 0xcc, 0x8e, 0x2a, 0x5a, 0xe2, 0xcb, 0x3b, 0x5a,
	0x32, 0xd1, 0x94, 0xed, 0x7d, 0x73, 0xc0, 0xb5, 0x65, 0x54, 0xfa, 0x74, 0xca, 0x83, 0xec, 0x35,
	0x65, 0x93, 0x06, 0x14, 0x78, 0x60, 0x5f, 0x29, 0xf0, 0x40, 0x31, 0x72, 0x91, 0xa0, 0x0d, 0xad,
	0x6d, 0xf2, 0x02, 0xdc, 0xe1, 0x9d, 0x1f, 0x45, 0x38, 0xd1, 0xa1, 0x6b, 0x9d, 0x3f, 0x1e, 0xf6,
	0xa7, 0x7d, 0x6a, 0x54, 0x96, 0x61, 0x64, 0x07, 0x4a, 0x23, 0x11, 0x87, 0xd4, 0xd1, 0xf8, 0xef,
	0x8f, 0xf0, 0xeb, 0x14, 0x05, 0xd3, 0x00, 0x79, 0x05, 0x55, 0x81, 0x43, 0x9e, 0x70, 0x8c, 0x24,
	0x2d, 0x6b, 0x9a, 0x3e, 0xa2, 0x59, 0xa6, 0xb3, 0x7b, 0x94, 0x6c, 0x42, 0x99, 0xa7, 0x83, 0xdb,
	0x58, 0x52, 0x77, 0x2b, 0xbf, 0x5b, 0x61, 0x0e, 0x4f, 0x4f, 0x62, 0x35, 0x9d, 0xeb, 0x3c, 0x1d,
	0x04, 0x5c, 0xe0, 0x50, 0x0e, 0xb2, 0x39, 0xad, 0x68, 0x62, 0x8d, 0xa7, 0x67, 0xda, 0x9f, 0x55,
	0xe3, 0x08, 0xaa, 0xcb, 0x31, 0xb7, 0xb3, 0xdc, 0x6a, 0x9b, 0x45, 0x68, 0x67, 0x8b, 0xd0, 0xee,
	0x67, 0x04, 0xbb, 0x87, 0xc9, 0x31, 0xd4, 0x24, 0xce, 0xe5, 0x60, 0x14, 0x8b, 0xd0, 0x97, 0x7a,
	0x62, 0x1a, 0x9d, 0x3f, 0x1f, 0xa5, 0xdd, 0xc7, 0xb9, 0x7c, 0xa3, 0x01, 0x06, 0x72, 0x69, 0xeb,
	0xfa, 0xe2, 0x5c, 0xd2, 0x4d, 0x5b, 0x5f, 0x9c, 0xcb, 0xd6, 0x15, 0x94, 0x54, 0x49, 0x6c, 0x2f,
	0xf2, 0xab, 0xbd, 0xd0, 0xd3, 0x51, 0x58, 0x99, 0x8e, 0x1d, 0x75, 0xdf, 0x0f, 0x75, 0x7f, 0x7e,
	0xae, 0x6c, 0x1f, 0xfd, 0x90, 0x69, 0xa0, 0xb5, 0x07, 0x25, 0x75, 0x7a, 0x4e, 0xd0, 0xd6, 0x3e,
	0xb8, 0xb6, 0x85, 0xcf, 0xc2, 0x6f, 0xa0, 0xba, 0x6c, 0xca, 0x2f, 0x4d, 0xda, 0xfb, 0x0b, 0xe0,
	0xbe, 0x6e, 0xa4, 0x01, 0x70, 0xf9, 0xfe, 0xf5, 0x45, 0x6f, 0xd0, 0xef, 0xde, 0xf4, 0x9b, 0x39,
	0xef, 0x7b, 0x1e, 0x1c, 0xbd, 0x90, 0x64, 0x1b, 0x4a, 0xc3, 0x38, 0x30, 0x7b, 0xd3, 0xe8, 0xac,
	0xaf, 0x2e, 0x6b, 0xfb, 0x34, 0x0e, 0x90, 0x69, 0x59, 0xad, 0x4a, 0x36, 0x04, 0x26, 0x9d, 0xe5,
	0xa7, 0xe9, 0x10, 0xdc, 0x00, 0xa5, 0xcf, 0x27, 0x29, 0x2d, 0x6e, 0x15, 0x77, 0x6b, 0xcb, 0xf6,
	0x99, 0x18, 0x67, 0x46, 0xeb, 0x46, 0x52, 0x2c, 0x58, 0x46, 0xb6, 0x8e, 0xa1, 0xbe, 0x2a, 0x90,
	0x26, 0x14, 0xbf, 0xe0, 0xc2, 0xfe, 0x77, 0x65, 0x92, 0x0d, 0x70, 0x66, 0xfe, 0x64, 0x9a, 0x3d,
	0x67, 0x0e, 0xc7, 0x85, 0xa3, 0xbc, 0xf7, 0x2f, 0x94, 0x54, 0x62, 0xa4, 0x06, 0xee, 0x75, 0xef,
	0x5d, 0xef, 0xe3, 0xa7, 0x5e, 0x33, 0x47, 0xea, 0x50, 0x61, 0xdd, 0xb7, 0x17, 0x57, 0xfd, 0x2e,
	0x6b, 0xe6, 0x3d, 0x17, 0x9c, 0x6e, 0x98, 0xc8, 0xc5, 0x6d, 0x59, 0xa7, 0x72, 0xf8, 0x23, 0x00,
	0x00, 0xff, 0xff, 0xcf, 0x67, 0xff, 0xf9, 0xac, 0x05, 0x00, 0x00,
}

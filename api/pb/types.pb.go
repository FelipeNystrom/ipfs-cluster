// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.3
// source: types.proto

package pb

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Pin_PinType int32

const (
	Pin_BadType        Pin_PinType = 0 // 1 << iota
	Pin_DataType       Pin_PinType = 1 // 2 << iota
	Pin_MetaType       Pin_PinType = 2
	Pin_ClusterDAGType Pin_PinType = 3
	Pin_ShardType      Pin_PinType = 4
)

// Enum value maps for Pin_PinType.
var (
	Pin_PinType_name = map[int32]string{
		0: "BadType",
		1: "DataType",
		2: "MetaType",
		3: "ClusterDAGType",
		4: "ShardType",
	}
	Pin_PinType_value = map[string]int32{
		"BadType":        0,
		"DataType":       1,
		"MetaType":       2,
		"ClusterDAGType": 3,
		"ShardType":      4,
	}
)

func (x Pin_PinType) Enum() *Pin_PinType {
	p := new(Pin_PinType)
	*p = x
	return p
}

func (x Pin_PinType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Pin_PinType) Descriptor() protoreflect.EnumDescriptor {
	return file_types_proto_enumTypes[0].Descriptor()
}

func (Pin_PinType) Type() protoreflect.EnumType {
	return &file_types_proto_enumTypes[0]
}

func (x Pin_PinType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Pin_PinType.Descriptor instead.
func (Pin_PinType) EnumDescriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{0, 0}
}

type Pin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cid         []byte      `protobuf:"bytes,1,opt,name=Cid,proto3" json:"Cid,omitempty"`
	Type        Pin_PinType `protobuf:"varint,2,opt,name=Type,proto3,enum=api.pb.Pin_PinType" json:"Type,omitempty"`
	Allocations [][]byte    `protobuf:"bytes,3,rep,name=Allocations,proto3" json:"Allocations,omitempty"`
	MaxDepth    int32       `protobuf:"zigzag32,4,opt,name=MaxDepth,proto3" json:"MaxDepth,omitempty"`
	Reference   []byte      `protobuf:"bytes,5,opt,name=Reference,proto3" json:"Reference,omitempty"`
	Options     *PinOptions `protobuf:"bytes,6,opt,name=Options,proto3" json:"Options,omitempty"`
	Timestamp   uint64      `protobuf:"varint,7,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
}

func (x *Pin) Reset() {
	*x = Pin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pin) ProtoMessage() {}

func (x *Pin) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pin.ProtoReflect.Descriptor instead.
func (*Pin) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{0}
}

func (x *Pin) GetCid() []byte {
	if x != nil {
		return x.Cid
	}
	return nil
}

func (x *Pin) GetType() Pin_PinType {
	if x != nil {
		return x.Type
	}
	return Pin_BadType
}

func (x *Pin) GetAllocations() [][]byte {
	if x != nil {
		return x.Allocations
	}
	return nil
}

func (x *Pin) GetMaxDepth() int32 {
	if x != nil {
		return x.MaxDepth
	}
	return 0
}

func (x *Pin) GetReference() []byte {
	if x != nil {
		return x.Reference
	}
	return nil
}

func (x *Pin) GetOptions() *PinOptions {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *Pin) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type PinOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReplicationFactorMin int32             `protobuf:"zigzag32,1,opt,name=ReplicationFactorMin,proto3" json:"ReplicationFactorMin,omitempty"`
	ReplicationFactorMax int32             `protobuf:"zigzag32,2,opt,name=ReplicationFactorMax,proto3" json:"ReplicationFactorMax,omitempty"`
	Name                 string            `protobuf:"bytes,3,opt,name=Name,proto3" json:"Name,omitempty"`
	ShardSize            uint64            `protobuf:"varint,4,opt,name=ShardSize,proto3" json:"ShardSize,omitempty"`
	Metadata             map[string]string `protobuf:"bytes,6,rep,name=Metadata,proto3" json:"Metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	PinUpdate            []byte            `protobuf:"bytes,7,opt,name=PinUpdate,proto3" json:"PinUpdate,omitempty"`
	ExpireAt             uint64            `protobuf:"varint,8,opt,name=ExpireAt,proto3" json:"ExpireAt,omitempty"`
	Origins              [][]byte          `protobuf:"bytes,9,rep,name=Origins,proto3" json:"Origins,omitempty"`
}

func (x *PinOptions) Reset() {
	*x = PinOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PinOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PinOptions) ProtoMessage() {}

func (x *PinOptions) ProtoReflect() protoreflect.Message {
	mi := &file_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PinOptions.ProtoReflect.Descriptor instead.
func (*PinOptions) Descriptor() ([]byte, []int) {
	return file_types_proto_rawDescGZIP(), []int{1}
}

func (x *PinOptions) GetReplicationFactorMin() int32 {
	if x != nil {
		return x.ReplicationFactorMin
	}
	return 0
}

func (x *PinOptions) GetReplicationFactorMax() int32 {
	if x != nil {
		return x.ReplicationFactorMax
	}
	return 0
}

func (x *PinOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PinOptions) GetShardSize() uint64 {
	if x != nil {
		return x.ShardSize
	}
	return 0
}

func (x *PinOptions) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *PinOptions) GetPinUpdate() []byte {
	if x != nil {
		return x.PinUpdate
	}
	return nil
}

func (x *PinOptions) GetExpireAt() uint64 {
	if x != nil {
		return x.ExpireAt
	}
	return 0
}

func (x *PinOptions) GetOrigins() [][]byte {
	if x != nil {
		return x.Origins
	}
	return nil
}

var File_types_proto protoreflect.FileDescriptor

var file_types_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61,
	0x70, 0x69, 0x2e, 0x70, 0x62, 0x22, 0xbf, 0x02, 0x0a, 0x03, 0x50, 0x69, 0x6e, 0x12, 0x10, 0x0a,
	0x03, 0x43, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x43, 0x69, 0x64, 0x12,
	0x27, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x69, 0x6e, 0x2e, 0x50, 0x69, 0x6e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x41, 0x6c, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0b, 0x41,
	0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x4d, 0x61,
	0x78, 0x44, 0x65, 0x70, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x11, 0x52, 0x08, 0x4d, 0x61,
	0x78, 0x44, 0x65, 0x70, 0x74, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x52, 0x65, 0x66, 0x65, 0x72,
	0x65, 0x6e, 0x63, 0x65, 0x12, 0x2c, 0x0a, 0x07, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x62, 0x2e, 0x50,
	0x69, 0x6e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x22, 0x55, 0x0a, 0x07, 0x50, 0x69, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x42,
	0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x61, 0x74, 0x61,
	0x54, 0x79, 0x70, 0x65, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x54, 0x79,
	0x70, 0x65, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x44,
	0x41, 0x47, 0x54, 0x79, 0x70, 0x65, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x68, 0x61, 0x72,
	0x64, 0x54, 0x79, 0x70, 0x65, 0x10, 0x04, 0x22, 0xfb, 0x02, 0x0a, 0x0a, 0x50, 0x69, 0x6e, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x32, 0x0a, 0x14, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x4d, 0x69, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x11, 0x52, 0x14, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x4d, 0x69, 0x6e, 0x12, 0x32, 0x0a, 0x14, 0x52, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x4d,
	0x61, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x11, 0x52, 0x14, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x4d, 0x61, 0x78, 0x12, 0x12,
	0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x68, 0x61, 0x72, 0x64, 0x53, 0x69, 0x7a, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x53, 0x68, 0x61, 0x72, 0x64, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x3c, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x20, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x69, 0x6e, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1c,
	0x0a, 0x09, 0x50, 0x69, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x50, 0x69, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x41, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08,
	0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x41, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x72, 0x69, 0x67,
	0x69, 0x6e, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x07, 0x4f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x73, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x4a,
	0x04, 0x08, 0x05, 0x10, 0x06, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_types_proto_rawDescOnce sync.Once
	file_types_proto_rawDescData = file_types_proto_rawDesc
)

func file_types_proto_rawDescGZIP() []byte {
	file_types_proto_rawDescOnce.Do(func() {
		file_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_types_proto_rawDescData)
	})
	return file_types_proto_rawDescData
}

var file_types_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_types_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_types_proto_goTypes = []interface{}{
	(Pin_PinType)(0),   // 0: api.pb.Pin.PinType
	(*Pin)(nil),        // 1: api.pb.Pin
	(*PinOptions)(nil), // 2: api.pb.PinOptions
	nil,                // 3: api.pb.PinOptions.MetadataEntry
}
var file_types_proto_depIdxs = []int32{
	0, // 0: api.pb.Pin.Type:type_name -> api.pb.Pin.PinType
	2, // 1: api.pb.Pin.Options:type_name -> api.pb.PinOptions
	3, // 2: api.pb.PinOptions.Metadata:type_name -> api.pb.PinOptions.MetadataEntry
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_types_proto_init() }
func file_types_proto_init() {
	if File_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pin); i {
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
		file_types_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PinOptions); i {
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
			RawDescriptor: file_types_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_types_proto_goTypes,
		DependencyIndexes: file_types_proto_depIdxs,
		EnumInfos:         file_types_proto_enumTypes,
		MessageInfos:      file_types_proto_msgTypes,
	}.Build()
	File_types_proto = out.File
	file_types_proto_rawDesc = nil
	file_types_proto_goTypes = nil
	file_types_proto_depIdxs = nil
}

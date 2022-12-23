// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: common/money.proto

package common

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

type CurrencyType int32

const (
	CurrencyType_CURRENCY_TYPE_UNSPECIFIED CurrencyType = 0
	CurrencyType_CURRENCY_TYPE_USD         CurrencyType = 1
)

// Enum value maps for CurrencyType.
var (
	CurrencyType_name = map[int32]string{
		0: "CURRENCY_TYPE_UNSPECIFIED",
		1: "CURRENCY_TYPE_USD",
	}
	CurrencyType_value = map[string]int32{
		"CURRENCY_TYPE_UNSPECIFIED": 0,
		"CURRENCY_TYPE_USD":         1,
	}
)

func (x CurrencyType) Enum() *CurrencyType {
	p := new(CurrencyType)
	*p = x
	return p
}

func (x CurrencyType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CurrencyType) Descriptor() protoreflect.EnumDescriptor {
	return file_common_money_proto_enumTypes[0].Descriptor()
}

func (CurrencyType) Type() protoreflect.EnumType {
	return &file_common_money_proto_enumTypes[0]
}

func (x CurrencyType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CurrencyType.Descriptor instead.
func (CurrencyType) EnumDescriptor() ([]byte, []int) {
	return file_common_money_proto_rawDescGZIP(), []int{0}
}

type Money struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount       string       `protobuf:"bytes,1,opt,name=Amount,proto3" json:"Amount,omitempty"`
	CurrencyType CurrencyType `protobuf:"varint,2,opt,name=currency_type,json=currencyType,proto3,enum=money.CurrencyType" json:"currency_type,omitempty"`
}

func (x *Money) Reset() {
	*x = Money{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_money_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Money) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Money) ProtoMessage() {}

func (x *Money) ProtoReflect() protoreflect.Message {
	mi := &file_common_money_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Money.ProtoReflect.Descriptor instead.
func (*Money) Descriptor() ([]byte, []int) {
	return file_common_money_proto_rawDescGZIP(), []int{0}
}

func (x *Money) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *Money) GetCurrencyType() CurrencyType {
	if x != nil {
		return x.CurrencyType
	}
	return CurrencyType_CURRENCY_TYPE_UNSPECIFIED
}

var File_common_money_proto protoreflect.FileDescriptor

var file_common_money_proto_rawDesc = []byte{
	0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x22, 0x59, 0x0a, 0x05, 0x4d,
	0x6f, 0x6e, 0x65, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x38, 0x0a, 0x0d,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x2e, 0x43, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x79, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x63, 0x79, 0x54, 0x79, 0x70, 0x65, 0x2a, 0x44, 0x0a, 0x0c, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x63, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x19, 0x43, 0x55, 0x52, 0x52, 0x45, 0x4e,
	0x43, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x55, 0x52, 0x52, 0x45, 0x4e, 0x43,
	0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x53, 0x44, 0x10, 0x01, 0x42, 0x72, 0x0a, 0x09,
	0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x42, 0x0a, 0x4d, 0x6f, 0x6e, 0x65, 0x79,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x25, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2d, 0x73,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0xa2, 0x02,
	0x03, 0x4d, 0x58, 0x58, 0xaa, 0x02, 0x05, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0xca, 0x02, 0x05, 0x4d,
	0x6f, 0x6e, 0x65, 0x79, 0xe2, 0x02, 0x11, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x05, 0x4d, 0x6f, 0x6e, 0x65, 0x79,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_money_proto_rawDescOnce sync.Once
	file_common_money_proto_rawDescData = file_common_money_proto_rawDesc
)

func file_common_money_proto_rawDescGZIP() []byte {
	file_common_money_proto_rawDescOnce.Do(func() {
		file_common_money_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_money_proto_rawDescData)
	})
	return file_common_money_proto_rawDescData
}

var file_common_money_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_common_money_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_money_proto_goTypes = []interface{}{
	(CurrencyType)(0), // 0: money.CurrencyType
	(*Money)(nil),     // 1: money.Money
}
var file_common_money_proto_depIdxs = []int32{
	0, // 0: money.Money.currency_type:type_name -> money.CurrencyType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_common_money_proto_init() }
func file_common_money_proto_init() {
	if File_common_money_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_money_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Money); i {
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
			RawDescriptor: file_common_money_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_money_proto_goTypes,
		DependencyIndexes: file_common_money_proto_depIdxs,
		EnumInfos:         file_common_money_proto_enumTypes,
		MessageInfos:      file_common_money_proto_msgTypes,
	}.Build()
	File_common_money_proto = out.File
	file_common_money_proto_rawDesc = nil
	file_common_money_proto_goTypes = nil
	file_common_money_proto_depIdxs = nil
}
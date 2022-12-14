// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: orders/orders.proto

package orders

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	common "order-sample/internal/protobuf/common"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OrderState int32

const (
	OrderState_ORDER_STATE_UNSPECIFIED OrderState = 0
	OrderState_ORDER_STATE_SUCCEEDED   OrderState = 1
	OrderState_ORDER_STATE_FAILED      OrderState = 2
)

// Enum value maps for OrderState.
var (
	OrderState_name = map[int32]string{
		0: "ORDER_STATE_UNSPECIFIED",
		1: "ORDER_STATE_SUCCEEDED",
		2: "ORDER_STATE_FAILED",
	}
	OrderState_value = map[string]int32{
		"ORDER_STATE_UNSPECIFIED": 0,
		"ORDER_STATE_SUCCEEDED":   1,
		"ORDER_STATE_FAILED":      2,
	}
)

func (x OrderState) Enum() *OrderState {
	p := new(OrderState)
	*p = x
	return p
}

func (x OrderState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrderState) Descriptor() protoreflect.EnumDescriptor {
	return file_orders_orders_proto_enumTypes[0].Descriptor()
}

func (OrderState) Type() protoreflect.EnumType {
	return &file_orders_orders_proto_enumTypes[0]
}

func (x OrderState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrderState.Descriptor instead.
func (OrderState) EnumDescriptor() ([]byte, []int) {
	return file_orders_orders_proto_rawDescGZIP(), []int{0}
}

type AssetType int32

const (
	AssetType_ASSET_TYPE_UNSPECIFIED   AssetType = 0
	AssetType_ASSET_TYPE_DAPPER_CREDIT AssetType = 1
)

// Enum value maps for AssetType.
var (
	AssetType_name = map[int32]string{
		0: "ASSET_TYPE_UNSPECIFIED",
		1: "ASSET_TYPE_DAPPER_CREDIT",
	}
	AssetType_value = map[string]int32{
		"ASSET_TYPE_UNSPECIFIED":   0,
		"ASSET_TYPE_DAPPER_CREDIT": 1,
	}
)

func (x AssetType) Enum() *AssetType {
	p := new(AssetType)
	*p = x
	return p
}

func (x AssetType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AssetType) Descriptor() protoreflect.EnumDescriptor {
	return file_orders_orders_proto_enumTypes[1].Descriptor()
}

func (AssetType) Type() protoreflect.EnumType {
	return &file_orders_orders_proto_enumTypes[1]
}

func (x AssetType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AssetType.Descriptor instead.
func (AssetType) EnumDescriptor() ([]byte, []int) {
	return file_orders_orders_proto_rawDescGZIP(), []int{1}
}

type Order struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId string        `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	State  OrderState    `protobuf:"varint,3,opt,name=state,proto3,enum=orders.OrderState" json:"state,omitempty"`
	Asset  *Asset        `protobuf:"bytes,4,opt,name=asset,proto3" json:"asset,omitempty"`
	Price  *common.Money `protobuf:"bytes,5,opt,name=price,proto3" json:"price,omitempty"`
}

func (x *Order) Reset() {
	*x = Order{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orders_orders_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_orders_orders_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Order.ProtoReflect.Descriptor instead.
func (*Order) Descriptor() ([]byte, []int) {
	return file_orders_orders_proto_rawDescGZIP(), []int{0}
}

func (x *Order) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Order) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Order) GetState() OrderState {
	if x != nil {
		return x.State
	}
	return OrderState_ORDER_STATE_UNSPECIFIED
}

func (x *Order) GetAsset() *Asset {
	if x != nil {
		return x.Asset
	}
	return nil
}

func (x *Order) GetPrice() *common.Money {
	if x != nil {
		return x.Price
	}
	return nil
}

type Asset struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AssetType AssetType `protobuf:"varint,2,opt,name=asset_type,json=assetType,proto3,enum=orders.AssetType" json:"asset_type,omitempty"`
	Name      string    `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Asset) Reset() {
	*x = Asset{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orders_orders_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Asset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Asset) ProtoMessage() {}

func (x *Asset) ProtoReflect() protoreflect.Message {
	mi := &file_orders_orders_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Asset.ProtoReflect.Descriptor instead.
func (*Asset) Descriptor() ([]byte, []int) {
	return file_orders_orders_proto_rawDescGZIP(), []int{1}
}

func (x *Asset) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Asset) GetAssetType() AssetType {
	if x != nil {
		return x.AssetType
	}
	return AssetType_ASSET_TYPE_UNSPECIFIED
}

func (x *Asset) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_orders_orders_proto protoreflect.FileDescriptor

var file_orders_orders_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x1a, 0x12, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xa3, 0x01, 0x0a, 0x05, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x2e, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x23,
	0x0a, 0x05, 0x61, 0x73, 0x73, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x05, 0x61, 0x73,
	0x73, 0x65, 0x74, 0x12, 0x22, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x2e, 0x4d, 0x6f, 0x6e, 0x65, 0x79,
	0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x22, 0x5d, 0x0a, 0x05, 0x41, 0x73, 0x73, 0x65, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x30, 0x0a, 0x0a, 0x61, 0x73, 0x73, 0x65, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x2e, 0x41, 0x73,
	0x73, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x61, 0x73, 0x73, 0x65, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x2a, 0x5c, 0x0a, 0x0a, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x0a, 0x17, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x19, 0x0a, 0x15, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45,
	0x5f, 0x53, 0x55, 0x43, 0x43, 0x45, 0x45, 0x44, 0x45, 0x44, 0x10, 0x01, 0x12, 0x16, 0x0a, 0x12,
	0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c,
	0x45, 0x44, 0x10, 0x02, 0x2a, 0x45, 0x0a, 0x09, 0x41, 0x73, 0x73, 0x65, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x1a, 0x0a, 0x16, 0x41, 0x53, 0x53, 0x45, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1c, 0x0a,
	0x18, 0x41, 0x53, 0x53, 0x45, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x41, 0x50, 0x50,
	0x45, 0x52, 0x5f, 0x43, 0x52, 0x45, 0x44, 0x49, 0x54, 0x10, 0x01, 0x42, 0x78, 0x0a, 0x0a, 0x63,
	0x6f, 0x6d, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x42, 0x0b, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x25, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2d,
	0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0xa2,
	0x02, 0x03, 0x4f, 0x58, 0x58, 0xaa, 0x02, 0x06, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0xca, 0x02,
	0x06, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0xe2, 0x02, 0x12, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x06, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_orders_orders_proto_rawDescOnce sync.Once
	file_orders_orders_proto_rawDescData = file_orders_orders_proto_rawDesc
)

func file_orders_orders_proto_rawDescGZIP() []byte {
	file_orders_orders_proto_rawDescOnce.Do(func() {
		file_orders_orders_proto_rawDescData = protoimpl.X.CompressGZIP(file_orders_orders_proto_rawDescData)
	})
	return file_orders_orders_proto_rawDescData
}

var file_orders_orders_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_orders_orders_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_orders_orders_proto_goTypes = []interface{}{
	(OrderState)(0),      // 0: orders.OrderState
	(AssetType)(0),       // 1: orders.AssetType
	(*Order)(nil),        // 2: orders.Order
	(*Asset)(nil),        // 3: orders.Asset
	(*common.Money)(nil), // 4: money.Money
}
var file_orders_orders_proto_depIdxs = []int32{
	0, // 0: orders.Order.state:type_name -> orders.OrderState
	3, // 1: orders.Order.asset:type_name -> orders.Asset
	4, // 2: orders.Order.price:type_name -> money.Money
	1, // 3: orders.Asset.asset_type:type_name -> orders.AssetType
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_orders_orders_proto_init() }
func file_orders_orders_proto_init() {
	if File_orders_orders_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_orders_orders_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Order); i {
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
		file_orders_orders_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Asset); i {
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
			RawDescriptor: file_orders_orders_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_orders_orders_proto_goTypes,
		DependencyIndexes: file_orders_orders_proto_depIdxs,
		EnumInfos:         file_orders_orders_proto_enumTypes,
		MessageInfos:      file_orders_orders_proto_msgTypes,
	}.Build()
	File_orders_orders_proto = out.File
	file_orders_orders_proto_rawDesc = nil
	file_orders_orders_proto_goTypes = nil
	file_orders_orders_proto_depIdxs = nil
}

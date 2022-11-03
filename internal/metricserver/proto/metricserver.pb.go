// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: proto/metricserver.proto

package proto

import (
	empty "github.com/golang/protobuf/ptypes/empty"
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

type Type int32

const (
	Type_UNKNOWN Type = 0
	Type_COUNTER Type = 1
	Type_GAUGE   Type = 2
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "UNKNOWN",
		1: "COUNTER",
		2: "GAUGE",
	}
	Type_value = map[string]int32{
		"UNKNOWN": 0,
		"COUNTER": 1,
		"GAUGE":   2,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_metricserver_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_proto_metricserver_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_proto_metricserver_proto_rawDescGZIP(), []int{0}
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type Type   `protobuf:"varint,2,opt,name=type,proto3,enum=metricserver.Type" json:"type,omitempty"`
	Hash string `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	// Types that are assignable to Value:
	//
	//	*Metric_Counter
	//	*Metric_Gauge
	Value isMetric_Value `protobuf_oneof:"value"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metricserver_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metricserver_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_proto_metricserver_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metric) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_UNKNOWN
}

func (x *Metric) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (m *Metric) GetValue() isMetric_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *Metric) GetCounter() int64 {
	if x, ok := x.GetValue().(*Metric_Counter); ok {
		return x.Counter
	}
	return 0
}

func (x *Metric) GetGauge() float64 {
	if x, ok := x.GetValue().(*Metric_Gauge); ok {
		return x.Gauge
	}
	return 0
}

type isMetric_Value interface {
	isMetric_Value()
}

type Metric_Counter struct {
	Counter int64 `protobuf:"varint,4,opt,name=counter,proto3,oneof"`
}

type Metric_Gauge struct {
	Gauge float64 `protobuf:"fixed64,5,opt,name=gauge,proto3,oneof"`
}

func (*Metric_Counter) isMetric_Value() {}

func (*Metric_Gauge) isMetric_Value() {}

type UpdateMetricBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *UpdateMetricBatchRequest) Reset() {
	*x = UpdateMetricBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metricserver_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricBatchRequest) ProtoMessage() {}

func (x *UpdateMetricBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metricserver_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricBatchRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricBatchRequest) Descriptor() ([]byte, []int) {
	return file_proto_metricserver_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateMetricBatchRequest) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type UpdateMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *UpdateMetricRequest) Reset() {
	*x = UpdateMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metricserver_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricRequest) ProtoMessage() {}

func (x *UpdateMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metricserver_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricRequest) Descriptor() ([]byte, []int) {
	return file_proto_metricserver_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateMetricRequest) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

type GetMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *GetMetricRequest) Reset() {
	*x = GetMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metricserver_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricRequest) ProtoMessage() {}

func (x *GetMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metricserver_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetricRequest.ProtoReflect.Descriptor instead.
func (*GetMetricRequest) Descriptor() ([]byte, []int) {
	return file_proto_metricserver_proto_rawDescGZIP(), []int{3}
}

func (x *GetMetricRequest) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

var File_proto_metricserver_proto protoreflect.FileDescriptor

var file_proto_metricserver_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x91, 0x01, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x26, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12,
	0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x1a, 0x0a, 0x07,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52,
	0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x05, 0x67, 0x61, 0x75, 0x67,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x05, 0x67, 0x61, 0x75, 0x67, 0x65,
	0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x4a, 0x0a, 0x18, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x07, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x43, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x06,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x22, 0x40, 0x0a, 0x10, 0x47, 0x65,
	0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c,
	0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2a, 0x2b, 0x0a, 0x04,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10,
	0x00, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x45, 0x52, 0x10, 0x01, 0x12, 0x09,
	0x0a, 0x05, 0x47, 0x41, 0x55, 0x47, 0x45, 0x10, 0x02, 0x32, 0xf9, 0x01, 0x0a, 0x0e, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x55, 0x0a, 0x11,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x12, 0x26, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x42, 0x61, 0x74,
	0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x12, 0x4b, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x12, 0x21, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00,
	0x12, 0x43, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x1e, 0x2e,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x22, 0x00, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x6d, 0x78, 0x30, 0x2f, 0x64, 0x65, 0x76, 0x6f, 0x70, 0x73, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_proto_metricserver_proto_rawDescOnce sync.Once
	file_proto_metricserver_proto_rawDescData = file_proto_metricserver_proto_rawDesc
)

func file_proto_metricserver_proto_rawDescGZIP() []byte {
	file_proto_metricserver_proto_rawDescOnce.Do(func() {
		file_proto_metricserver_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_metricserver_proto_rawDescData)
	})
	return file_proto_metricserver_proto_rawDescData
}

var file_proto_metricserver_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_metricserver_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_metricserver_proto_goTypes = []interface{}{
	(Type)(0),                        // 0: metricserver.Type
	(*Metric)(nil),                   // 1: metricserver.Metric
	(*UpdateMetricBatchRequest)(nil), // 2: metricserver.UpdateMetricBatchRequest
	(*UpdateMetricRequest)(nil),      // 3: metricserver.UpdateMetricRequest
	(*GetMetricRequest)(nil),         // 4: metricserver.GetMetricRequest
	(*empty.Empty)(nil),              // 5: google.protobuf.Empty
}
var file_proto_metricserver_proto_depIdxs = []int32{
	0, // 0: metricserver.Metric.type:type_name -> metricserver.Type
	1, // 1: metricserver.UpdateMetricBatchRequest.metrics:type_name -> metricserver.Metric
	1, // 2: metricserver.UpdateMetricRequest.metric:type_name -> metricserver.Metric
	1, // 3: metricserver.GetMetricRequest.metric:type_name -> metricserver.Metric
	2, // 4: metricserver.MetricsService.UpdateMetricBatch:input_type -> metricserver.UpdateMetricBatchRequest
	3, // 5: metricserver.MetricsService.UpdateMetric:input_type -> metricserver.UpdateMetricRequest
	4, // 6: metricserver.MetricsService.GetMetric:input_type -> metricserver.GetMetricRequest
	5, // 7: metricserver.MetricsService.UpdateMetricBatch:output_type -> google.protobuf.Empty
	5, // 8: metricserver.MetricsService.UpdateMetric:output_type -> google.protobuf.Empty
	1, // 9: metricserver.MetricsService.GetMetric:output_type -> metricserver.Metric
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_metricserver_proto_init() }
func file_proto_metricserver_proto_init() {
	if File_proto_metricserver_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_metricserver_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
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
		file_proto_metricserver_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricBatchRequest); i {
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
		file_proto_metricserver_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricRequest); i {
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
		file_proto_metricserver_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMetricRequest); i {
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
	file_proto_metricserver_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Metric_Counter)(nil),
		(*Metric_Gauge)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_metricserver_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_metricserver_proto_goTypes,
		DependencyIndexes: file_proto_metricserver_proto_depIdxs,
		EnumInfos:         file_proto_metricserver_proto_enumTypes,
		MessageInfos:      file_proto_metricserver_proto_msgTypes,
	}.Build()
	File_proto_metricserver_proto = out.File
	file_proto_metricserver_proto_rawDesc = nil
	file_proto_metricserver_proto_goTypes = nil
	file_proto_metricserver_proto_depIdxs = nil
}

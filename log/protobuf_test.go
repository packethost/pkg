package log

import (
	"reflect"
	"sync"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ProtoTest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Foo string `protobuf:"bytes,1,opt,name=foo,proto3" json:"foo,omitempty"`
}

func (x *ProtoTest) Reset() {
	*x = ProtoTest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtoTest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoTest) ProtoMessage() {}

func (x *ProtoTest) ProtoReflect() protoreflect.Message {
	mi := &file_p_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoTest.ProtoReflect.Descriptor instead.
func (*ProtoTest) Descriptor() ([]byte, []int) {
	return file_p_proto_rawDescGZIP(), []int{0}
}

func (x *ProtoTest) GetFoo() string {
	if x != nil {
		return x.Foo
	}
	return ""
}

var File_p_proto protoreflect.FileDescriptor

var file_p_proto_rawDesc = []byte{
	0x0a, 0x07, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1d, 0x0a, 0x09, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x54, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x6f, 0x6f, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x66, 0x6f, 0x6f, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_p_proto_rawDescOnce sync.Once
	file_p_proto_rawDescData = file_p_proto_rawDesc
)

func file_p_proto_rawDescGZIP() []byte {
	file_p_proto_rawDescOnce.Do(func() {
		file_p_proto_rawDescData = protoimpl.X.CompressGZIP(file_p_proto_rawDescData)
	})
	return file_p_proto_rawDescData
}

var file_p_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_p_proto_goTypes = []interface{}{
	(*ProtoTest)(nil), // 0: ProtoTest
}
var file_p_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_p_proto_init() }
func file_p_proto_init() {
	if File_p_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_p_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtoTest); i {
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
			RawDescriptor: file_p_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_p_proto_goTypes,
		DependencyIndexes: file_p_proto_depIdxs,
		MessageInfos:      file_p_proto_msgTypes,
	}.Build()
	File_p_proto = out.File
	file_p_proto_rawDesc = nil
	file_p_proto_goTypes = nil
	file_p_proto_depIdxs = nil
}

func TestProtoAsJSON(t *testing.T) {
	enabler := zap.NewAtomicLevelAt(zap.InfoLevel)
	core, logs := observer.New(enabler)

	logger, err := configureLogger(zap.New(core), "test")
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	var val ProtoTest
	val.Foo = "bar"

	logger.With("proto_field", ProtoAsJSON(&val)).Info("ignore this message")

	if cnt := logs.Len(); cnt != 1 {
		t.Errorf("unexpected entry count: %d", cnt)
		t.FailNow()
	}

	entries := logs.TakeAll()

	if cnt := len(entries[0].Context); cnt != 2 {
		t.Errorf("unexpected field count: %d", cnt)
		t.FailNow()
	}

	if name := entries[0].Context[1].Key; name != "proto_field" {
		t.Errorf("unexpected field name: %s", name)
	}

	if typ := entries[0].Context[1].Type; typ != zapcore.ReflectType {
		t.Errorf("unexpected field type: %v", typ)
	}
}

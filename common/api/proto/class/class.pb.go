// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: api/proto/class/class.proto

package class

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

type Class struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	ImageUrl      string                 `protobuf:"bytes,4,opt,name=imageUrl,proto3" json:"imageUrl,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,5,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt     string                 `protobuf:"bytes,6,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Class) Reset() {
	*x = Class{}
	mi := &file_api_proto_class_class_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Class) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Class) ProtoMessage() {}

func (x *Class) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_class_class_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Class.ProtoReflect.Descriptor instead.
func (*Class) Descriptor() ([]byte, []int) {
	return file_api_proto_class_class_proto_rawDescGZIP(), []int{0}
}

func (x *Class) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Class) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Class) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Class) GetImageUrl() string {
	if x != nil {
		return x.ImageUrl
	}
	return ""
}

func (x *Class) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *Class) GetUpdatedAt() string {
	if x != nil {
		return x.UpdatedAt
	}
	return ""
}

type Ascendancy struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ClassId       string                 `protobuf:"bytes,2,opt,name=classId,proto3" json:"classId,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	ImageUrl      string                 `protobuf:"bytes,5,opt,name=imageUrl,proto3" json:"imageUrl,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,6,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt     string                 `protobuf:"bytes,7,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Ascendancy) Reset() {
	*x = Ascendancy{}
	mi := &file_api_proto_class_class_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Ascendancy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ascendancy) ProtoMessage() {}

func (x *Ascendancy) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_class_class_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ascendancy.ProtoReflect.Descriptor instead.
func (*Ascendancy) Descriptor() ([]byte, []int) {
	return file_api_proto_class_class_proto_rawDescGZIP(), []int{1}
}

func (x *Ascendancy) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Ascendancy) GetClassId() string {
	if x != nil {
		return x.ClassId
	}
	return ""
}

func (x *Ascendancy) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Ascendancy) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Ascendancy) GetImageUrl() string {
	if x != nil {
		return x.ImageUrl
	}
	return ""
}

func (x *Ascendancy) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *Ascendancy) GetUpdatedAt() string {
	if x != nil {
		return x.UpdatedAt
	}
	return ""
}

type GetClassesAndAscendanciesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MemberId      string                 `protobuf:"bytes,1,opt,name=MemberId,proto3" json:"MemberId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetClassesAndAscendanciesRequest) Reset() {
	*x = GetClassesAndAscendanciesRequest{}
	mi := &file_api_proto_class_class_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetClassesAndAscendanciesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClassesAndAscendanciesRequest) ProtoMessage() {}

func (x *GetClassesAndAscendanciesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_class_class_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClassesAndAscendanciesRequest.ProtoReflect.Descriptor instead.
func (*GetClassesAndAscendanciesRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_class_class_proto_rawDescGZIP(), []int{2}
}

func (x *GetClassesAndAscendanciesRequest) GetMemberId() string {
	if x != nil {
		return x.MemberId
	}
	return ""
}

type GetClassesAndAscendanciesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Classes       []*Class               `protobuf:"bytes,1,rep,name=classes,proto3" json:"classes,omitempty"`
	Ascendancies  []*Ascendancy          `protobuf:"bytes,2,rep,name=ascendancies,proto3" json:"ascendancies,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetClassesAndAscendanciesResponse) Reset() {
	*x = GetClassesAndAscendanciesResponse{}
	mi := &file_api_proto_class_class_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetClassesAndAscendanciesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClassesAndAscendanciesResponse) ProtoMessage() {}

func (x *GetClassesAndAscendanciesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_class_class_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClassesAndAscendanciesResponse.ProtoReflect.Descriptor instead.
func (*GetClassesAndAscendanciesResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_class_class_proto_rawDescGZIP(), []int{3}
}

func (x *GetClassesAndAscendanciesResponse) GetClasses() []*Class {
	if x != nil {
		return x.Classes
	}
	return nil
}

func (x *GetClassesAndAscendanciesResponse) GetAscendancies() []*Ascendancy {
	if x != nil {
		return x.Ascendancies
	}
	return nil
}

var File_api_proto_class_class_proto protoreflect.FileDescriptor

const file_api_proto_class_class_proto_rawDesc = "" +
	"\n" +
	"\x1bapi/proto/class/class.proto\x12\x05class\"\xa5\x01\n" +
	"\x05Class\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x12\x1a\n" +
	"\bimageUrl\x18\x04 \x01(\tR\bimageUrl\x12\x1c\n" +
	"\tcreatedAt\x18\x05 \x01(\tR\tcreatedAt\x12\x1c\n" +
	"\tupdatedAt\x18\x06 \x01(\tR\tupdatedAt\"\xc4\x01\n" +
	"\n" +
	"Ascendancy\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x18\n" +
	"\aclassId\x18\x02 \x01(\tR\aclassId\x12\x12\n" +
	"\x04name\x18\x03 \x01(\tR\x04name\x12 \n" +
	"\vdescription\x18\x04 \x01(\tR\vdescription\x12\x1a\n" +
	"\bimageUrl\x18\x05 \x01(\tR\bimageUrl\x12\x1c\n" +
	"\tcreatedAt\x18\x06 \x01(\tR\tcreatedAt\x12\x1c\n" +
	"\tupdatedAt\x18\a \x01(\tR\tupdatedAt\">\n" +
	" GetClassesAndAscendanciesRequest\x12\x1a\n" +
	"\bMemberId\x18\x01 \x01(\tR\bMemberId\"\x82\x01\n" +
	"!GetClassesAndAscendanciesResponse\x12&\n" +
	"\aclasses\x18\x01 \x03(\v2\f.class.ClassR\aclasses\x125\n" +
	"\fascendancies\x18\x02 \x03(\v2\x11.class.AscendancyR\fascendancies2\x80\x01\n" +
	"\fClassService\x12p\n" +
	"\x19GetClassesAndAscendancies\x12'.class.GetClassesAndAscendanciesRequest\x1a(.class.GetClassesAndAscendanciesResponse\"\x00BOZMgithub.com/darkphotonKN/community-builds-microservices/common/api/proto/classb\x06proto3"

var (
	file_api_proto_class_class_proto_rawDescOnce sync.Once
	file_api_proto_class_class_proto_rawDescData []byte
)

func file_api_proto_class_class_proto_rawDescGZIP() []byte {
	file_api_proto_class_class_proto_rawDescOnce.Do(func() {
		file_api_proto_class_class_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_proto_class_class_proto_rawDesc), len(file_api_proto_class_class_proto_rawDesc)))
	})
	return file_api_proto_class_class_proto_rawDescData
}

var file_api_proto_class_class_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_proto_class_class_proto_goTypes = []any{
	(*Class)(nil),                             // 0: class.Class
	(*Ascendancy)(nil),                        // 1: class.Ascendancy
	(*GetClassesAndAscendanciesRequest)(nil),  // 2: class.GetClassesAndAscendanciesRequest
	(*GetClassesAndAscendanciesResponse)(nil), // 3: class.GetClassesAndAscendanciesResponse
}
var file_api_proto_class_class_proto_depIdxs = []int32{
	0, // 0: class.GetClassesAndAscendanciesResponse.classes:type_name -> class.Class
	1, // 1: class.GetClassesAndAscendanciesResponse.ascendancies:type_name -> class.Ascendancy
	2, // 2: class.ClassService.GetClassesAndAscendancies:input_type -> class.GetClassesAndAscendanciesRequest
	3, // 3: class.ClassService.GetClassesAndAscendancies:output_type -> class.GetClassesAndAscendanciesResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_proto_class_class_proto_init() }
func file_api_proto_class_class_proto_init() {
	if File_api_proto_class_class_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_proto_class_class_proto_rawDesc), len(file_api_proto_class_class_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_class_class_proto_goTypes,
		DependencyIndexes: file_api_proto_class_class_proto_depIdxs,
		MessageInfos:      file_api_proto_class_class_proto_msgTypes,
	}.Build()
	File_api_proto_class_class_proto = out.File
	file_api_proto_class_class_proto_goTypes = nil
	file_api_proto_class_class_proto_depIdxs = nil
}

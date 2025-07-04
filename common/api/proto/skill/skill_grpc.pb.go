// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: api/proto/skill/skill.proto

package skill

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SkillService_CreateSkill_FullMethodName = "/skill.SkillService/CreateSkill"
	SkillService_GetSkills_FullMethodName   = "/skill.SkillService/GetSkills"
)

// SkillServiceClient is the client API for SkillService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SkillServiceClient interface {
	CreateSkill(ctx context.Context, in *CreateSkillRequest, opts ...grpc.CallOption) (*CreateSkillResponse, error)
	GetSkills(ctx context.Context, in *GetSkillsRequest, opts ...grpc.CallOption) (*GetSkillsResponse, error)
}

type skillServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSkillServiceClient(cc grpc.ClientConnInterface) SkillServiceClient {
	return &skillServiceClient{cc}
}

func (c *skillServiceClient) CreateSkill(ctx context.Context, in *CreateSkillRequest, opts ...grpc.CallOption) (*CreateSkillResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateSkillResponse)
	err := c.cc.Invoke(ctx, SkillService_CreateSkill_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *skillServiceClient) GetSkills(ctx context.Context, in *GetSkillsRequest, opts ...grpc.CallOption) (*GetSkillsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetSkillsResponse)
	err := c.cc.Invoke(ctx, SkillService_GetSkills_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SkillServiceServer is the server API for SkillService service.
// All implementations must embed UnimplementedSkillServiceServer
// for forward compatibility.
type SkillServiceServer interface {
	CreateSkill(context.Context, *CreateSkillRequest) (*CreateSkillResponse, error)
	GetSkills(context.Context, *GetSkillsRequest) (*GetSkillsResponse, error)
	mustEmbedUnimplementedSkillServiceServer()
}

// UnimplementedSkillServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSkillServiceServer struct{}

func (UnimplementedSkillServiceServer) CreateSkill(context.Context, *CreateSkillRequest) (*CreateSkillResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSkill not implemented")
}
func (UnimplementedSkillServiceServer) GetSkills(context.Context, *GetSkillsRequest) (*GetSkillsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSkills not implemented")
}
func (UnimplementedSkillServiceServer) mustEmbedUnimplementedSkillServiceServer() {}
func (UnimplementedSkillServiceServer) testEmbeddedByValue()                      {}

// UnsafeSkillServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SkillServiceServer will
// result in compilation errors.
type UnsafeSkillServiceServer interface {
	mustEmbedUnimplementedSkillServiceServer()
}

func RegisterSkillServiceServer(s grpc.ServiceRegistrar, srv SkillServiceServer) {
	// If the following call pancis, it indicates UnimplementedSkillServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SkillService_ServiceDesc, srv)
}

func _SkillService_CreateSkill_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSkillRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SkillServiceServer).CreateSkill(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SkillService_CreateSkill_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SkillServiceServer).CreateSkill(ctx, req.(*CreateSkillRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SkillService_GetSkills_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSkillsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SkillServiceServer).GetSkills(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SkillService_GetSkills_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SkillServiceServer).GetSkills(ctx, req.(*GetSkillsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SkillService_ServiceDesc is the grpc.ServiceDesc for SkillService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SkillService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "skill.SkillService",
	HandlerType: (*SkillServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSkill",
			Handler:    _SkillService_CreateSkill_Handler,
		},
		{
			MethodName: "GetSkills",
			Handler:    _SkillService_GetSkills_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/skill/skill.proto",
}

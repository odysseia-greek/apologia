// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.29.3
// source: proto/antisthenes.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AntisthenesClient is the client API for Antisthenes service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AntisthenesClient interface {
	Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error)
	Options(ctx context.Context, in *OptionsRequest, opts ...grpc.CallOption) (*AggregatedOptions, error)
	Question(ctx context.Context, in *CreationRequest, opts ...grpc.CallOption) (*QuizResponse, error)
	Answer(ctx context.Context, in *AnswerRequest, opts ...grpc.CallOption) (*ComprehensiveResponse, error)
}

type antisthenesClient struct {
	cc grpc.ClientConnInterface
}

func NewAntisthenesClient(cc grpc.ClientConnInterface) AntisthenesClient {
	return &antisthenesClient{cc}
}

func (c *antisthenesClient) Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error) {
	out := new(HealthResponse)
	err := c.cc.Invoke(ctx, "/apologia_antisthenes.Antisthenes/Health", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *antisthenesClient) Options(ctx context.Context, in *OptionsRequest, opts ...grpc.CallOption) (*AggregatedOptions, error) {
	out := new(AggregatedOptions)
	err := c.cc.Invoke(ctx, "/apologia_antisthenes.Antisthenes/Options", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *antisthenesClient) Question(ctx context.Context, in *CreationRequest, opts ...grpc.CallOption) (*QuizResponse, error) {
	out := new(QuizResponse)
	err := c.cc.Invoke(ctx, "/apologia_antisthenes.Antisthenes/Question", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *antisthenesClient) Answer(ctx context.Context, in *AnswerRequest, opts ...grpc.CallOption) (*ComprehensiveResponse, error) {
	out := new(ComprehensiveResponse)
	err := c.cc.Invoke(ctx, "/apologia_antisthenes.Antisthenes/Answer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AntisthenesServer is the server API for Antisthenes service.
// All implementations must embed UnimplementedAntisthenesServer
// for forward compatibility
type AntisthenesServer interface {
	Health(context.Context, *HealthRequest) (*HealthResponse, error)
	Options(context.Context, *OptionsRequest) (*AggregatedOptions, error)
	Question(context.Context, *CreationRequest) (*QuizResponse, error)
	Answer(context.Context, *AnswerRequest) (*ComprehensiveResponse, error)
	mustEmbedUnimplementedAntisthenesServer()
}

// UnimplementedAntisthenesServer must be embedded to have forward compatible implementations.
type UnimplementedAntisthenesServer struct {
}

func (UnimplementedAntisthenesServer) Health(context.Context, *HealthRequest) (*HealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedAntisthenesServer) Options(context.Context, *OptionsRequest) (*AggregatedOptions, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Options not implemented")
}
func (UnimplementedAntisthenesServer) Question(context.Context, *CreationRequest) (*QuizResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Question not implemented")
}
func (UnimplementedAntisthenesServer) Answer(context.Context, *AnswerRequest) (*ComprehensiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Answer not implemented")
}
func (UnimplementedAntisthenesServer) mustEmbedUnimplementedAntisthenesServer() {}

// UnsafeAntisthenesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AntisthenesServer will
// result in compilation errors.
type UnsafeAntisthenesServer interface {
	mustEmbedUnimplementedAntisthenesServer()
}

func RegisterAntisthenesServer(s grpc.ServiceRegistrar, srv AntisthenesServer) {
	s.RegisterService(&Antisthenes_ServiceDesc, srv)
}

func _Antisthenes_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AntisthenesServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apologia_antisthenes.Antisthenes/Health",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AntisthenesServer).Health(ctx, req.(*HealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Antisthenes_Options_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OptionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AntisthenesServer).Options(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apologia_antisthenes.Antisthenes/Options",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AntisthenesServer).Options(ctx, req.(*OptionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Antisthenes_Question_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AntisthenesServer).Question(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apologia_antisthenes.Antisthenes/Question",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AntisthenesServer).Question(ctx, req.(*CreationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Antisthenes_Answer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnswerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AntisthenesServer).Answer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apologia_antisthenes.Antisthenes/Answer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AntisthenesServer).Answer(ctx, req.(*AnswerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Antisthenes_ServiceDesc is the grpc.ServiceDesc for Antisthenes service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Antisthenes_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "apologia_antisthenes.Antisthenes",
	HandlerType: (*AntisthenesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Health",
			Handler:    _Antisthenes_Health_Handler,
		},
		{
			MethodName: "Options",
			Handler:    _Antisthenes_Options_Handler,
		},
		{
			MethodName: "Question",
			Handler:    _Antisthenes_Question_Handler,
		},
		{
			MethodName: "Answer",
			Handler:    _Antisthenes_Answer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/antisthenes.proto",
}

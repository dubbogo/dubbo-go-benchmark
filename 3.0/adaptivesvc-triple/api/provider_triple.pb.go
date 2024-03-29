// Code generated by protoc-gen-go-triple. DO NOT EDIT.
// versions:
// - protoc-gen-go-triple v1.0.5
// - protoc             v3.19.4
// source: provider.proto

package api

import (
	context "context"
	protocol "dubbo.apache.org/dubbo-go/v3/protocol"
	dubbo3 "dubbo.apache.org/dubbo-go/v3/protocol/dubbo3"
	invocation "dubbo.apache.org/dubbo-go/v3/protocol/invocation"
	grpc_go "github.com/dubbogo/grpc-go"
	codes "github.com/dubbogo/grpc-go/codes"
	metadata "github.com/dubbogo/grpc-go/metadata"
	status "github.com/dubbogo/grpc-go/status"
	common "github.com/dubbogo/triple/pkg/common"
	constant "github.com/dubbogo/triple/pkg/common/constant"
	triple "github.com/dubbogo/triple/pkg/triple"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc_go.SupportPackageIsVersion7

// ProviderClient is the client API for Provider service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProviderClient interface {
	Fibonacci(ctx context.Context, in *FibonacciRequest, opts ...grpc_go.CallOption) (*FibonacciResult, common.ErrorWithAttachment)
	Sleep(ctx context.Context, in *SleepRequest, opts ...grpc_go.CallOption) (*SleepResult, common.ErrorWithAttachment)
}

type providerClient struct {
	cc *triple.TripleConn
}

type ProviderClientImpl struct {
	Fibonacci func(ctx context.Context, in *FibonacciRequest) (*FibonacciResult, error)
	Sleep     func(ctx context.Context, in *SleepRequest) (*SleepResult, error)
}

func (c *ProviderClientImpl) GetDubboStub(cc *triple.TripleConn) ProviderClient {
	return NewProviderClient(cc)
}

func (c *ProviderClientImpl) XXX_InterfaceName() string {
	return "api.Provider"
}

func NewProviderClient(cc *triple.TripleConn) ProviderClient {
	return &providerClient{cc}
}

func (c *providerClient) Fibonacci(ctx context.Context, in *FibonacciRequest, opts ...grpc_go.CallOption) (*FibonacciResult, common.ErrorWithAttachment) {
	out := new(FibonacciResult)
	interfaceKey := ctx.Value(constant.InterfaceKey).(string)
	return out, c.cc.Invoke(ctx, "/"+interfaceKey+"/Fibonacci", in, out)
}

func (c *providerClient) Sleep(ctx context.Context, in *SleepRequest, opts ...grpc_go.CallOption) (*SleepResult, common.ErrorWithAttachment) {
	out := new(SleepResult)
	interfaceKey := ctx.Value(constant.InterfaceKey).(string)
	return out, c.cc.Invoke(ctx, "/"+interfaceKey+"/Sleep", in, out)
}

// ProviderServer is the server API for Provider service.
// All implementations must embed UnimplementedProviderServer
// for forward compatibility
type ProviderServer interface {
	Fibonacci(context.Context, *FibonacciRequest) (*FibonacciResult, error)
	Sleep(context.Context, *SleepRequest) (*SleepResult, error)
	mustEmbedUnimplementedProviderServer()
}

// UnimplementedProviderServer must be embedded to have forward compatible implementations.
type UnimplementedProviderServer struct {
	proxyImpl protocol.Invoker
}

func (UnimplementedProviderServer) Fibonacci(context.Context, *FibonacciRequest) (*FibonacciResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fibonacci not implemented")
}
func (UnimplementedProviderServer) Sleep(context.Context, *SleepRequest) (*SleepResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sleep not implemented")
}
func (s *UnimplementedProviderServer) XXX_SetProxyImpl(impl protocol.Invoker) {
	s.proxyImpl = impl
}

func (s *UnimplementedProviderServer) XXX_GetProxyImpl() protocol.Invoker {
	return s.proxyImpl
}

func (s *UnimplementedProviderServer) XXX_ServiceDesc() *grpc_go.ServiceDesc {
	return &Provider_ServiceDesc
}
func (s *UnimplementedProviderServer) XXX_InterfaceName() string {
	return "api.Provider"
}

func (UnimplementedProviderServer) mustEmbedUnimplementedProviderServer() {}

// UnsafeProviderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProviderServer will
// result in compilation errors.
type UnsafeProviderServer interface {
	mustEmbedUnimplementedProviderServer()
}

func RegisterProviderServer(s grpc_go.ServiceRegistrar, srv ProviderServer) {
	s.RegisterService(&Provider_ServiceDesc, srv)
}

func _Provider_Fibonacci_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc_go.UnaryServerInterceptor) (interface{}, error) {
	in := new(FibonacciRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	base := srv.(dubbo3.Dubbo3GrpcService)
	args := []interface{}{}
	args = append(args, in)
	md, _ := metadata.FromIncomingContext(ctx)
	invAttachment := make(map[string]interface{}, len(md))
	for k, v := range md {
		invAttachment[k] = v
	}
	invo := invocation.NewRPCInvocation("Fibonacci", args, invAttachment)
	if interceptor == nil {
		result := base.XXX_GetProxyImpl().Invoke(ctx, invo)
		return result, result.Error()
	}
	info := &grpc_go.UnaryServerInfo{
		Server:     srv,
		FullMethod: ctx.Value("XXX_TRIPLE_GO_INTERFACE_NAME").(string),
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		result := base.XXX_GetProxyImpl().Invoke(ctx, invo)
		return result, result.Error()
	}
	return interceptor(ctx, in, info, handler)
}

func _Provider_Sleep_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc_go.UnaryServerInterceptor) (interface{}, error) {
	in := new(SleepRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	base := srv.(dubbo3.Dubbo3GrpcService)
	args := []interface{}{}
	args = append(args, in)
	md, _ := metadata.FromIncomingContext(ctx)
	invAttachment := make(map[string]interface{}, len(md))
	for k, v := range md {
		invAttachment[k] = v
	}
	invo := invocation.NewRPCInvocation("Sleep", args, invAttachment)
	if interceptor == nil {
		result := base.XXX_GetProxyImpl().Invoke(ctx, invo)
		return result, result.Error()
	}
	info := &grpc_go.UnaryServerInfo{
		Server:     srv,
		FullMethod: ctx.Value("XXX_TRIPLE_GO_INTERFACE_NAME").(string),
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		result := base.XXX_GetProxyImpl().Invoke(ctx, invo)
		return result, result.Error()
	}
	return interceptor(ctx, in, info, handler)
}

// Provider_ServiceDesc is the grpc_go.ServiceDesc for Provider service.
// It's only intended for direct use with grpc_go.RegisterService,
// and not to be introspected or modified (even as a copy)
var Provider_ServiceDesc = grpc_go.ServiceDesc{
	ServiceName: "api.Provider",
	HandlerType: (*ProviderServer)(nil),
	Methods: []grpc_go.MethodDesc{
		{
			MethodName: "Fibonacci",
			Handler:    _Provider_Fibonacci_Handler,
		},
		{
			MethodName: "Sleep",
			Handler:    _Provider_Sleep_Handler,
		},
	},
	Streams:  []grpc_go.StreamDesc{},
	Metadata: "provider.proto",
}

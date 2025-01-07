// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: game.proto

package game_backend

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
	GameService_CreateGame_FullMethodName        = "/game.GameService/CreateGame"
	GameService_JoinGame_FullMethodName          = "/game.GameService/JoinGame"
	GameService_UpdateGameState_FullMethodName   = "/game.GameService/UpdateGameState"
	GameService_StreamGameUpdates_FullMethodName = "/game.GameService/StreamGameUpdates"
)

// GameServiceClient is the client API for GameService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GameServiceClient interface {
	CreateGame(ctx context.Context, in *CreateGameRequest, opts ...grpc.CallOption) (*CreateGameResponse, error)
	JoinGame(ctx context.Context, in *JoinGameRequest, opts ...grpc.CallOption) (*JoinGameResponse, error)
	UpdateGameState(ctx context.Context, in *UpdateGameStateRequest, opts ...grpc.CallOption) (*UpdateGameStateResponse, error)
	StreamGameUpdates(ctx context.Context, in *StreamGameUpdatesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[StreamGameUpdatesResponse], error)
}

type gameServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGameServiceClient(cc grpc.ClientConnInterface) GameServiceClient {
	return &gameServiceClient{cc}
}

func (c *gameServiceClient) CreateGame(ctx context.Context, in *CreateGameRequest, opts ...grpc.CallOption) (*CreateGameResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateGameResponse)
	err := c.cc.Invoke(ctx, GameService_CreateGame_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameServiceClient) JoinGame(ctx context.Context, in *JoinGameRequest, opts ...grpc.CallOption) (*JoinGameResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(JoinGameResponse)
	err := c.cc.Invoke(ctx, GameService_JoinGame_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameServiceClient) UpdateGameState(ctx context.Context, in *UpdateGameStateRequest, opts ...grpc.CallOption) (*UpdateGameStateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateGameStateResponse)
	err := c.cc.Invoke(ctx, GameService_UpdateGameState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameServiceClient) StreamGameUpdates(ctx context.Context, in *StreamGameUpdatesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[StreamGameUpdatesResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &GameService_ServiceDesc.Streams[0], GameService_StreamGameUpdates_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[StreamGameUpdatesRequest, StreamGameUpdatesResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type GameService_StreamGameUpdatesClient = grpc.ServerStreamingClient[StreamGameUpdatesResponse]

// GameServiceServer is the server API for GameService service.
// All implementations must embed UnimplementedGameServiceServer
// for forward compatibility.
type GameServiceServer interface {
	CreateGame(context.Context, *CreateGameRequest) (*CreateGameResponse, error)
	JoinGame(context.Context, *JoinGameRequest) (*JoinGameResponse, error)
	UpdateGameState(context.Context, *UpdateGameStateRequest) (*UpdateGameStateResponse, error)
	StreamGameUpdates(*StreamGameUpdatesRequest, grpc.ServerStreamingServer[StreamGameUpdatesResponse]) error
	mustEmbedUnimplementedGameServiceServer()
}

// UnimplementedGameServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGameServiceServer struct{}

func (UnimplementedGameServiceServer) CreateGame(context.Context, *CreateGameRequest) (*CreateGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGame not implemented")
}
func (UnimplementedGameServiceServer) JoinGame(context.Context, *JoinGameRequest) (*JoinGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinGame not implemented")
}
func (UnimplementedGameServiceServer) UpdateGameState(context.Context, *UpdateGameStateRequest) (*UpdateGameStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateGameState not implemented")
}
func (UnimplementedGameServiceServer) StreamGameUpdates(*StreamGameUpdatesRequest, grpc.ServerStreamingServer[StreamGameUpdatesResponse]) error {
	return status.Errorf(codes.Unimplemented, "method StreamGameUpdates not implemented")
}
func (UnimplementedGameServiceServer) mustEmbedUnimplementedGameServiceServer() {}
func (UnimplementedGameServiceServer) testEmbeddedByValue()                     {}

// UnsafeGameServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GameServiceServer will
// result in compilation errors.
type UnsafeGameServiceServer interface {
	mustEmbedUnimplementedGameServiceServer()
}

func RegisterGameServiceServer(s grpc.ServiceRegistrar, srv GameServiceServer) {
	// If the following call pancis, it indicates UnimplementedGameServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GameService_ServiceDesc, srv)
}

func _GameService_CreateGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServiceServer).CreateGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GameService_CreateGame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServiceServer).CreateGame(ctx, req.(*CreateGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameService_JoinGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServiceServer).JoinGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GameService_JoinGame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServiceServer).JoinGame(ctx, req.(*JoinGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameService_UpdateGameState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGameStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServiceServer).UpdateGameState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GameService_UpdateGameState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServiceServer).UpdateGameState(ctx, req.(*UpdateGameStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameService_StreamGameUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamGameUpdatesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GameServiceServer).StreamGameUpdates(m, &grpc.GenericServerStream[StreamGameUpdatesRequest, StreamGameUpdatesResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type GameService_StreamGameUpdatesServer = grpc.ServerStreamingServer[StreamGameUpdatesResponse]

// GameService_ServiceDesc is the grpc.ServiceDesc for GameService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GameService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "game.GameService",
	HandlerType: (*GameServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGame",
			Handler:    _GameService_CreateGame_Handler,
		},
		{
			MethodName: "JoinGame",
			Handler:    _GameService_JoinGame_Handler,
		},
		{
			MethodName: "UpdateGameState",
			Handler:    _GameService_UpdateGameState_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamGameUpdates",
			Handler:       _GameService_StreamGameUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "game.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// AuthenticationClient is the client API for Authentication service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthenticationClient interface {
	RegisterUser(ctx context.Context, in *UserRegisterRequest, opts ...grpc.CallOption) (*UserRegisterResponse, error)
	LoginUser(ctx context.Context, in *UserLoginRequest, opts ...grpc.CallOption) (*UserLoginResponse, error)
}

type authenticationClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthenticationClient(cc grpc.ClientConnInterface) AuthenticationClient {
	return &authenticationClient{cc}
}

func (c *authenticationClient) RegisterUser(ctx context.Context, in *UserRegisterRequest, opts ...grpc.CallOption) (*UserRegisterResponse, error) {
	out := new(UserRegisterResponse)
	err := c.cc.Invoke(ctx, "/main.Authentication/RegisterUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationClient) LoginUser(ctx context.Context, in *UserLoginRequest, opts ...grpc.CallOption) (*UserLoginResponse, error) {
	out := new(UserLoginResponse)
	err := c.cc.Invoke(ctx, "/main.Authentication/LoginUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthenticationServer is the server API for Authentication service.
// All implementations must embed UnimplementedAuthenticationServer
// for forward compatibility
type AuthenticationServer interface {
	RegisterUser(context.Context, *UserRegisterRequest) (*UserRegisterResponse, error)
	LoginUser(context.Context, *UserLoginRequest) (*UserLoginResponse, error)
	mustEmbedUnimplementedAuthenticationServer()
}

// UnimplementedAuthenticationServer must be embedded to have forward compatible implementations.
type UnimplementedAuthenticationServer struct {
}

func (UnimplementedAuthenticationServer) RegisterUser(context.Context, *UserRegisterRequest) (*UserRegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterUser not implemented")
}
func (UnimplementedAuthenticationServer) LoginUser(context.Context, *UserLoginRequest) (*UserLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedAuthenticationServer) mustEmbedUnimplementedAuthenticationServer() {}

// UnsafeAuthenticationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthenticationServer will
// result in compilation errors.
type UnsafeAuthenticationServer interface {
	mustEmbedUnimplementedAuthenticationServer()
}

func RegisterAuthenticationServer(s grpc.ServiceRegistrar, srv AuthenticationServer) {
	s.RegisterService(&_Authentication_serviceDesc, srv)
}

func _Authentication_RegisterUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServer).RegisterUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.Authentication/RegisterUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServer).RegisterUser(ctx, req.(*UserRegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authentication_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.Authentication/LoginUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServer).LoginUser(ctx, req.(*UserLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Authentication_serviceDesc = grpc.ServiceDesc{
	ServiceName: "main.Authentication",
	HandlerType: (*AuthenticationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterUser",
			Handler:    _Authentication_RegisterUser_Handler,
		},
		{
			MethodName: "LoginUser",
			Handler:    _Authentication_LoginUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "logic/proto/api.proto",
}

// NewPhotosClient is the client API for NewPhotos service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NewPhotosClient interface {
	// Returns all the photos of the passed user ...
	AllPhotos(ctx context.Context, in *AllPhotosRequest, opts ...grpc.CallOption) (NewPhotos_AllPhotosClient, error)
	// Uploads equal photo of the passed user ...
	UploadEqualPhoto(ctx context.Context, in *UploadEqualPhotoRequest, opts ...grpc.CallOption) (*UploadEqualPhotoResponse, error)
	// Returns current available storage for storing of the passed user ...
	GetUserinfo(ctx context.Context, in *GetUserinfoRequest, opts ...grpc.CallOption) (*GetUserinfoResponse, error)
	GetFullPhotoByThumbnail(ctx context.Context, in *GetFullPhotoByThumbnailRequest, opts ...grpc.CallOption) (*GetFullPhotoByThumbnailResponse, error)
	AllPhotosAlbum(ctx context.Context, in *AllPhotosAlbumRequest, opts ...grpc.CallOption) (NewPhotos_AllPhotosAlbumClient, error)
	GetAllAlbums(ctx context.Context, in *GetAllAlbumsRequest, opts ...grpc.CallOption) (NewPhotos_GetAllAlbumsClient, error)
	CreateAlbum(ctx context.Context, in *CreateAlbumRequest, opts ...grpc.CallOption) (*CreateAlbumResponse, error)
	DeleteAlbum(ctx context.Context, in *DeleteAlbumRequest, opts ...grpc.CallOption) (*DeleteAlbumResponse, error)
	UploadPhotoToAlbum(ctx context.Context, opts ...grpc.CallOption) (NewPhotos_UploadPhotoToAlbumClient, error)
}

type newPhotosClient struct {
	cc grpc.ClientConnInterface
}

func NewNewPhotosClient(cc grpc.ClientConnInterface) NewPhotosClient {
	return &newPhotosClient{cc}
}

func (c *newPhotosClient) AllPhotos(ctx context.Context, in *AllPhotosRequest, opts ...grpc.CallOption) (NewPhotos_AllPhotosClient, error) {
	stream, err := c.cc.NewStream(ctx, &_NewPhotos_serviceDesc.Streams[0], "/main.NewPhotos/AllPhotos", opts...)
	if err != nil {
		return nil, err
	}
	x := &newPhotosAllPhotosClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NewPhotos_AllPhotosClient interface {
	Recv() (*AllPhotosResponse, error)
	grpc.ClientStream
}

type newPhotosAllPhotosClient struct {
	grpc.ClientStream
}

func (x *newPhotosAllPhotosClient) Recv() (*AllPhotosResponse, error) {
	m := new(AllPhotosResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *newPhotosClient) UploadEqualPhoto(ctx context.Context, in *UploadEqualPhotoRequest, opts ...grpc.CallOption) (*UploadEqualPhotoResponse, error) {
	out := new(UploadEqualPhotoResponse)
	err := c.cc.Invoke(ctx, "/main.NewPhotos/UploadEqualPhoto", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newPhotosClient) GetUserinfo(ctx context.Context, in *GetUserinfoRequest, opts ...grpc.CallOption) (*GetUserinfoResponse, error) {
	out := new(GetUserinfoResponse)
	err := c.cc.Invoke(ctx, "/main.NewPhotos/GetUserinfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newPhotosClient) GetFullPhotoByThumbnail(ctx context.Context, in *GetFullPhotoByThumbnailRequest, opts ...grpc.CallOption) (*GetFullPhotoByThumbnailResponse, error) {
	out := new(GetFullPhotoByThumbnailResponse)
	err := c.cc.Invoke(ctx, "/main.NewPhotos/GetFullPhotoByThumbnail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newPhotosClient) AllPhotosAlbum(ctx context.Context, in *AllPhotosAlbumRequest, opts ...grpc.CallOption) (NewPhotos_AllPhotosAlbumClient, error) {
	stream, err := c.cc.NewStream(ctx, &_NewPhotos_serviceDesc.Streams[1], "/main.NewPhotos/AllPhotosAlbum", opts...)
	if err != nil {
		return nil, err
	}
	x := &newPhotosAllPhotosAlbumClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NewPhotos_AllPhotosAlbumClient interface {
	Recv() (*AllPhotosAlbumResponse, error)
	grpc.ClientStream
}

type newPhotosAllPhotosAlbumClient struct {
	grpc.ClientStream
}

func (x *newPhotosAllPhotosAlbumClient) Recv() (*AllPhotosAlbumResponse, error) {
	m := new(AllPhotosAlbumResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *newPhotosClient) GetAllAlbums(ctx context.Context, in *GetAllAlbumsRequest, opts ...grpc.CallOption) (NewPhotos_GetAllAlbumsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_NewPhotos_serviceDesc.Streams[2], "/main.NewPhotos/GetAllAlbums", opts...)
	if err != nil {
		return nil, err
	}
	x := &newPhotosGetAllAlbumsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NewPhotos_GetAllAlbumsClient interface {
	Recv() (*GetAllAlbumsResponse, error)
	grpc.ClientStream
}

type newPhotosGetAllAlbumsClient struct {
	grpc.ClientStream
}

func (x *newPhotosGetAllAlbumsClient) Recv() (*GetAllAlbumsResponse, error) {
	m := new(GetAllAlbumsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *newPhotosClient) CreateAlbum(ctx context.Context, in *CreateAlbumRequest, opts ...grpc.CallOption) (*CreateAlbumResponse, error) {
	out := new(CreateAlbumResponse)
	err := c.cc.Invoke(ctx, "/main.NewPhotos/CreateAlbum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newPhotosClient) DeleteAlbum(ctx context.Context, in *DeleteAlbumRequest, opts ...grpc.CallOption) (*DeleteAlbumResponse, error) {
	out := new(DeleteAlbumResponse)
	err := c.cc.Invoke(ctx, "/main.NewPhotos/DeleteAlbum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newPhotosClient) UploadPhotoToAlbum(ctx context.Context, opts ...grpc.CallOption) (NewPhotos_UploadPhotoToAlbumClient, error) {
	stream, err := c.cc.NewStream(ctx, &_NewPhotos_serviceDesc.Streams[3], "/main.NewPhotos/UploadPhotoToAlbum", opts...)
	if err != nil {
		return nil, err
	}
	x := &newPhotosUploadPhotoToAlbumClient{stream}
	return x, nil
}

type NewPhotos_UploadPhotoToAlbumClient interface {
	Send(*UploadPhotoToAlbumRequest) error
	CloseAndRecv() (*UploadPhotoToAlbumResponse, error)
	grpc.ClientStream
}

type newPhotosUploadPhotoToAlbumClient struct {
	grpc.ClientStream
}

func (x *newPhotosUploadPhotoToAlbumClient) Send(m *UploadPhotoToAlbumRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *newPhotosUploadPhotoToAlbumClient) CloseAndRecv() (*UploadPhotoToAlbumResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadPhotoToAlbumResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NewPhotosServer is the server API for NewPhotos service.
// All implementations must embed UnimplementedNewPhotosServer
// for forward compatibility
type NewPhotosServer interface {
	// Returns all the photos of the passed user ...
	AllPhotos(*AllPhotosRequest, NewPhotos_AllPhotosServer) error
	// Uploads equal photo of the passed user ...
	UploadEqualPhoto(context.Context, *UploadEqualPhotoRequest) (*UploadEqualPhotoResponse, error)
	// Returns current available storage for storing of the passed user ...
	GetUserinfo(context.Context, *GetUserinfoRequest) (*GetUserinfoResponse, error)
	GetFullPhotoByThumbnail(context.Context, *GetFullPhotoByThumbnailRequest) (*GetFullPhotoByThumbnailResponse, error)
	AllPhotosAlbum(*AllPhotosAlbumRequest, NewPhotos_AllPhotosAlbumServer) error
	GetAllAlbums(*GetAllAlbumsRequest, NewPhotos_GetAllAlbumsServer) error
	CreateAlbum(context.Context, *CreateAlbumRequest) (*CreateAlbumResponse, error)
	DeleteAlbum(context.Context, *DeleteAlbumRequest) (*DeleteAlbumResponse, error)
	UploadPhotoToAlbum(NewPhotos_UploadPhotoToAlbumServer) error
	mustEmbedUnimplementedNewPhotosServer()
}

// UnimplementedNewPhotosServer must be embedded to have forward compatible implementations.
type UnimplementedNewPhotosServer struct {
}

func (UnimplementedNewPhotosServer) AllPhotos(*AllPhotosRequest, NewPhotos_AllPhotosServer) error {
	return status.Errorf(codes.Unimplemented, "method AllPhotos not implemented")
}
func (UnimplementedNewPhotosServer) UploadEqualPhoto(context.Context, *UploadEqualPhotoRequest) (*UploadEqualPhotoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadEqualPhoto not implemented")
}
func (UnimplementedNewPhotosServer) GetUserinfo(context.Context, *GetUserinfoRequest) (*GetUserinfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserinfo not implemented")
}
func (UnimplementedNewPhotosServer) GetFullPhotoByThumbnail(context.Context, *GetFullPhotoByThumbnailRequest) (*GetFullPhotoByThumbnailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFullPhotoByThumbnail not implemented")
}
func (UnimplementedNewPhotosServer) AllPhotosAlbum(*AllPhotosAlbumRequest, NewPhotos_AllPhotosAlbumServer) error {
	return status.Errorf(codes.Unimplemented, "method AllPhotosAlbum not implemented")
}
func (UnimplementedNewPhotosServer) GetAllAlbums(*GetAllAlbumsRequest, NewPhotos_GetAllAlbumsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAllAlbums not implemented")
}
func (UnimplementedNewPhotosServer) CreateAlbum(context.Context, *CreateAlbumRequest) (*CreateAlbumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAlbum not implemented")
}
func (UnimplementedNewPhotosServer) DeleteAlbum(context.Context, *DeleteAlbumRequest) (*DeleteAlbumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAlbum not implemented")
}
func (UnimplementedNewPhotosServer) UploadPhotoToAlbum(NewPhotos_UploadPhotoToAlbumServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadPhotoToAlbum not implemented")
}
func (UnimplementedNewPhotosServer) mustEmbedUnimplementedNewPhotosServer() {}

// UnsafeNewPhotosServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NewPhotosServer will
// result in compilation errors.
type UnsafeNewPhotosServer interface {
	mustEmbedUnimplementedNewPhotosServer()
}

func RegisterNewPhotosServer(s grpc.ServiceRegistrar, srv NewPhotosServer) {
	s.RegisterService(&_NewPhotos_serviceDesc, srv)
}

func _NewPhotos_AllPhotos_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AllPhotosRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NewPhotosServer).AllPhotos(m, &newPhotosAllPhotosServer{stream})
}

type NewPhotos_AllPhotosServer interface {
	Send(*AllPhotosResponse) error
	grpc.ServerStream
}

type newPhotosAllPhotosServer struct {
	grpc.ServerStream
}

func (x *newPhotosAllPhotosServer) Send(m *AllPhotosResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _NewPhotos_UploadEqualPhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadEqualPhotoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewPhotosServer).UploadEqualPhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.NewPhotos/UploadEqualPhoto",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewPhotosServer).UploadEqualPhoto(ctx, req.(*UploadEqualPhotoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewPhotos_GetUserinfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserinfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewPhotosServer).GetUserinfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.NewPhotos/GetUserinfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewPhotosServer).GetUserinfo(ctx, req.(*GetUserinfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewPhotos_GetFullPhotoByThumbnail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFullPhotoByThumbnailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewPhotosServer).GetFullPhotoByThumbnail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.NewPhotos/GetFullPhotoByThumbnail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewPhotosServer).GetFullPhotoByThumbnail(ctx, req.(*GetFullPhotoByThumbnailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewPhotos_AllPhotosAlbum_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AllPhotosAlbumRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NewPhotosServer).AllPhotosAlbum(m, &newPhotosAllPhotosAlbumServer{stream})
}

type NewPhotos_AllPhotosAlbumServer interface {
	Send(*AllPhotosAlbumResponse) error
	grpc.ServerStream
}

type newPhotosAllPhotosAlbumServer struct {
	grpc.ServerStream
}

func (x *newPhotosAllPhotosAlbumServer) Send(m *AllPhotosAlbumResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _NewPhotos_GetAllAlbums_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetAllAlbumsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NewPhotosServer).GetAllAlbums(m, &newPhotosGetAllAlbumsServer{stream})
}

type NewPhotos_GetAllAlbumsServer interface {
	Send(*GetAllAlbumsResponse) error
	grpc.ServerStream
}

type newPhotosGetAllAlbumsServer struct {
	grpc.ServerStream
}

func (x *newPhotosGetAllAlbumsServer) Send(m *GetAllAlbumsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _NewPhotos_CreateAlbum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAlbumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewPhotosServer).CreateAlbum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.NewPhotos/CreateAlbum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewPhotosServer).CreateAlbum(ctx, req.(*CreateAlbumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewPhotos_DeleteAlbum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAlbumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewPhotosServer).DeleteAlbum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.NewPhotos/DeleteAlbum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewPhotosServer).DeleteAlbum(ctx, req.(*DeleteAlbumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewPhotos_UploadPhotoToAlbum_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(NewPhotosServer).UploadPhotoToAlbum(&newPhotosUploadPhotoToAlbumServer{stream})
}

type NewPhotos_UploadPhotoToAlbumServer interface {
	SendAndClose(*UploadPhotoToAlbumResponse) error
	Recv() (*UploadPhotoToAlbumRequest, error)
	grpc.ServerStream
}

type newPhotosUploadPhotoToAlbumServer struct {
	grpc.ServerStream
}

func (x *newPhotosUploadPhotoToAlbumServer) SendAndClose(m *UploadPhotoToAlbumResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *newPhotosUploadPhotoToAlbumServer) Recv() (*UploadPhotoToAlbumRequest, error) {
	m := new(UploadPhotoToAlbumRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _NewPhotos_serviceDesc = grpc.ServiceDesc{
	ServiceName: "main.NewPhotos",
	HandlerType: (*NewPhotosServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UploadEqualPhoto",
			Handler:    _NewPhotos_UploadEqualPhoto_Handler,
		},
		{
			MethodName: "GetUserinfo",
			Handler:    _NewPhotos_GetUserinfo_Handler,
		},
		{
			MethodName: "GetFullPhotoByThumbnail",
			Handler:    _NewPhotos_GetFullPhotoByThumbnail_Handler,
		},
		{
			MethodName: "CreateAlbum",
			Handler:    _NewPhotos_CreateAlbum_Handler,
		},
		{
			MethodName: "DeleteAlbum",
			Handler:    _NewPhotos_DeleteAlbum_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AllPhotos",
			Handler:       _NewPhotos_AllPhotos_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "AllPhotosAlbum",
			Handler:       _NewPhotos_AllPhotosAlbum_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetAllAlbums",
			Handler:       _NewPhotos_GetAllAlbums_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "UploadPhotoToAlbum",
			Handler:       _NewPhotos_UploadPhotoToAlbum_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "logic/proto/api.proto",
}

// TagClient is the client API for Tag service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TagClient interface {
	IsHuman(ctx context.Context, in *IsHumanRequest, opts ...grpc.CallOption) (*IsHumanResponse, error)
	IsDog(ctx context.Context, in *IsDogRequest, opts ...grpc.CallOption) (*IsDogResponse, error)
}

type tagClient struct {
	cc grpc.ClientConnInterface
}

func NewTagClient(cc grpc.ClientConnInterface) TagClient {
	return &tagClient{cc}
}

func (c *tagClient) IsHuman(ctx context.Context, in *IsHumanRequest, opts ...grpc.CallOption) (*IsHumanResponse, error) {
	out := new(IsHumanResponse)
	err := c.cc.Invoke(ctx, "/main.Tag/IsHuman", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tagClient) IsDog(ctx context.Context, in *IsDogRequest, opts ...grpc.CallOption) (*IsDogResponse, error) {
	out := new(IsDogResponse)
	err := c.cc.Invoke(ctx, "/main.Tag/IsDog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TagServer is the server API for Tag service.
// All implementations must embed UnimplementedTagServer
// for forward compatibility
type TagServer interface {
	IsHuman(context.Context, *IsHumanRequest) (*IsHumanResponse, error)
	IsDog(context.Context, *IsDogRequest) (*IsDogResponse, error)
	mustEmbedUnimplementedTagServer()
}

// UnimplementedTagServer must be embedded to have forward compatible implementations.
type UnimplementedTagServer struct {
}

func (UnimplementedTagServer) IsHuman(context.Context, *IsHumanRequest) (*IsHumanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsHuman not implemented")
}
func (UnimplementedTagServer) IsDog(context.Context, *IsDogRequest) (*IsDogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsDog not implemented")
}
func (UnimplementedTagServer) mustEmbedUnimplementedTagServer() {}

// UnsafeTagServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TagServer will
// result in compilation errors.
type UnsafeTagServer interface {
	mustEmbedUnimplementedTagServer()
}

func RegisterTagServer(s grpc.ServiceRegistrar, srv TagServer) {
	s.RegisterService(&_Tag_serviceDesc, srv)
}

func _Tag_IsHuman_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsHumanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TagServer).IsHuman(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.Tag/IsHuman",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TagServer).IsHuman(ctx, req.(*IsHumanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tag_IsDog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsDogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TagServer).IsDog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.Tag/IsDog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TagServer).IsDog(ctx, req.(*IsDogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Tag_serviceDesc = grpc.ServiceDesc{
	ServiceName: "main.Tag",
	HandlerType: (*TagServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsHuman",
			Handler:    _Tag_IsHuman_Handler,
		},
		{
			MethodName: "IsDog",
			Handler:    _Tag_IsDog_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "logic/proto/api.proto",
}

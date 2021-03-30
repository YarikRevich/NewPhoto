package server

import (
	"NewPhoto/logic/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func Run(){
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(50 * 10e6),
		grpc.MaxSendMsgSize(50 * 10e6),
	}
	s := grpc.NewServer(opts...)
	newphoto := proto.NewNewPhoto()
	proto.RegisterNewPhotosServer(s, newphoto)
	newphoto.InitTagClient()
	
	auth := proto.NewAuthentication()
	proto.RegisterAuthenticationServer(s, auth)

	l, err := net.Listen("tcp", ":8082")
	if err != nil{
		log.Fatalln(err)
	}
	if err := s.Serve(l); err != nil{
		log.Fatalln(err)
	}
}
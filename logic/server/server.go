package server

import (
	"NewPhoto/logic/proto"
	
	"net"
	"os"

	"NewPhoto/log"
	"google.golang.org/grpc"
)

func Run() {
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(50 * 10e6),
		grpc.MaxSendMsgSize(50 * 10e6),
	}
	s := grpc.NewServer(opts...)
	newphoto := proto.NewNewPhoto()
	proto.RegisterNewPhotosServer(s, newphoto)

	auth := proto.NewAuthentication()
	proto.RegisterAuthenticationServer(s, auth)

	runAddr, ok := os.LookupEnv("runAddr")
	if !ok {
		log.Logger.Fatalln("runAddr is not written in credentials.sh file")
	}

	l, err := net.Listen("tcp", runAddr)
	if err != nil {
		log.Logger.Fatalln(err)
	}
	if err := s.Serve(l); err != nil {
		log.Logger.Fatalln(err)
	}
}

package client

import (
	"github.com/YarikRevich/NewPhoto/logic/proto"
	"log"
	"os"

	"google.golang.org/grpc"
)

func NewTagClient() proto.TagClient {
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(32*10e6),
			grpc.MaxCallSendMsgSize(32*10e6),
		),
		grpc.WithInsecure(),
	}

	runAddr, ok := os.LookupEnv("runAddr")
	if !ok {
		log.Fatalln(ok)
	}

	c, err := grpc.Dial(runAddr, opts...)
	if err != nil {
		log.Fatalln(err)
	}
	return proto.NewTagClient(c)
}

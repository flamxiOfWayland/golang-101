package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/flamxiOfWayland/golang-101/grpc/greeter"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	greeter.UnimplementedGreeterServer
}

func (s *server) SayHello(_ context.Context, in *greeter.HelloRequest) (*greeter.HelloReply, error) {
	var msg = "Hello " + in.GetName()
	return &greeter.HelloReply{Message: &msg}, nil
}

// func (s *server) SayHelloAgain(_ context.Context, in *greeter.HelloRequest) (*greeter.HelloReply, error) {
// 	var msg = "Hello Again foo" + in.GetName()
// 	return &greeter.HelloReply{Message: &msg}, nil
// }

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greeter.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

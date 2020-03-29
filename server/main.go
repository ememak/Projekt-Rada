package main

import (
	"context"
	"fmt"
	"net"

	pb "github.com/ememak/Projekt-Rada/hello"
	"google.golang.org/grpc"
)

var (
	port = ":12345"
)

type server struct {
	pb.UnimplementedHelloServer
}

func (s *server) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Mess: "Hello World!\n"}, nil
}

func main() {
	rec, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &server{})
	s.Serve(rec)
}

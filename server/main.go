package main

import (
	"context"
	"fmt"
	"net"
  "io"

	pb "github.com/ememak/Projekt-Rada/hello"
	"google.golang.org/grpc"
)

var (
	port = ":12345"
)

type server struct {
	pb.UnimplementedQueryServer

  queries []pb.PollQuestion
  tokens []int
}

func (s *server) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Mess: "Hello World!\n"}, nil
}

func (s *server) QueryInit(stream pb.Query_QueryInitServer) error {
  var q = pb.PollQuestion{
    Id: int32(len(s.queries)),
    Fields: []*pb.PollQuestion_QueryField{
      },
  }
  fmt.Printf("QueryInitReceived, id = %v\n", len(s.queries))
  s.tokens = append(s.tokens, 0)
  for {
    field, err := stream.Recv()
    if err == io.EOF {
      s.queries = append(s.queries, q)
      fmt.Printf("Sending back: %v\n", q)
      fmt.Printf("In Memory: %v\n", s.queries)
      return stream.SendAndClose(&q)
    }
    if err != nil {
			return err
		}
    if field.Which == -1 {
      q.Fields = append(q.Fields,
        &pb.PollQuestion_QueryField{
          Name: field.Name,
          Votes: 0,
        })
    } else if field.Which < int32(len(q.Fields)){
      q.Fields[field.Which].Name = field.Name
    } else {
      fmt.Printf("Wrong Field Number")
    }
  }
}

func (s *server) QueryGetToken(ctx context.Context, in *pb.TokenRequest) (*pb.VoteToken, error){
  if in.Nr < 0 || in.Nr >= int32(len(s.queries)) {
    fmt.Printf("")
  }
  s.tokens[in.Nr]++
  return &pb.VoteToken{Token: int32(s.tokens[in.Nr]-1)}, nil
}

func main() {
	rec, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterQueryServer(s, &server{})
	s.Serve(rec)
}

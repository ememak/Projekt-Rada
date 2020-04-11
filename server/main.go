package main

import (
	"context"
	"fmt"
	"net"
  "io"

	pb "github.com/ememak/Projekt-Rada/query"
	"google.golang.org/grpc"
)

var (
	port = ":12345"
)

type server struct {
	pb.UnimplementedQueryServer

  queries []pb.PollQuestion
}

func (s *server) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Mess: "Hello World!\n"}, nil
}

func (s *server) QueryInit(stream pb.Query_QueryInitServer) error {
  var q = pb.PollQuestion{
    Id: int32(len(s.queries)),
    Fields: []*pb.PollQuestion_QueryField{},
  }
  fmt.Printf("QueryInitReceived, id = %v\n", len(s.queries))
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
      fmt.Printf("Wrong Query field number\n")
    }
  }
}

//for now tokens for votes dont have much informations, no crypto yet
func (s *server) QueryGetToken(ctx context.Context, in *pb.TokenRequest) (*pb.VoteToken, error){
  var t = pb.VoteToken{
    Token: int32(len(s.queries[in.Nr].Tokens)+1),}
  if in.Nr < 0 || in.Nr >= int32(len(s.queries)) {
    fmt.Printf("Wrong Query number in Get Token\n")
    t.Token = -1
    return &t, nil
  }
  s.queries[in.Nr].Tokens = append(s.queries[in.Nr].Tokens, &t)
  fmt.Printf("GetToken, in Memory: %v\n", s.queries)
  return &t, nil
}

//for now tokens for votes dont have much informations, no crypto yet
func (s *server) QueryVote(ctx context.Context, in *pb.Vote) (*pb.VoteReply, error) {
  if in.Nr >= int32(len(s.queries)) || in.Nr < 0 { //security leak, path that out later!
    fmt.Printf("No such Query: %v\n", in.Nr)
	  return &pb.VoteReply{Mess: "No such Query!\n"}, nil
  }
  var nT = len(s.queries[in.Nr].Tokens)
  var nF = len(s.queries[in.Nr].Fields)
  for i:=0; i<nT; i++ {
    if s.queries[in.Nr].Tokens[i].Token == in.Token.Token {
      s.queries[in.Nr].Tokens[i] = s.queries[in.Nr].Tokens[nT-1]
      s.queries[in.Nr].Tokens[nT-1] = &pb.VoteToken{Token:0}
      s.queries[in.Nr].Tokens = s.queries[in.Nr].Tokens[:nT-1]
      break
    }
    if i == nT-1 {
      fmt.Printf("Token not valid\n")
	    return &pb.VoteReply{Mess: "Token not valid!\n"}, nil
    }
  }
  if nF != len(in.Answer) {
    fmt.Printf("Vote have different number of fields than query\n")
	  return &pb.VoteReply{Mess: "Vote have different number of fields than query!\n"}, nil
  }
  for i:=0; i<nF; i++ {
    if in.Answer[i] >= 1 {
      s.queries[in.Nr].Fields[i].Votes++
    }
  }
  fmt.Printf("In Memory: %v\n", s.queries)
	return &pb.VoteReply{Mess: "Thank you for your vote!\n"}, nil
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

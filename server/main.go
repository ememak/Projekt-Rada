package main

import (
  "bytes"
	"context"
	"fmt"
	"net"
  "io"
  "crypto/rand"
  "crypto/rsa"
  "crypto/sha256"
  "math/big"

	pb "github.com/ememak/Projekt-Rada/query"
	"google.golang.org/grpc"
)

var (
	port = ":12345"
)

type server struct {
	pb.UnimplementedQueryServer

  key     *rsa.PrivateKey
  queries []pb.PollQuestion
}

func (s *server) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
    Mess: "Hello World!\n",
    N: s.key.PublicKey.N.Bytes(),
    E: int32(s.key.PublicKey.E),
  }, nil
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
    Token: -1,}
  if in.Nr < 0 || in.Nr >= int32(len(s.queries)) {
    fmt.Printf("Wrong Query number in Get Token\n")
    return &t, nil
  }
  t.Token = int32(len(s.queries[in.Nr].Tokens)+1)
  s.queries[in.Nr].Tokens = append(s.queries[in.Nr].Tokens, &t)
  fmt.Printf("GetToken, in Memory: %v\n", s.queries)
  return &t, nil
}

func (s *server) QueryAuthorizeVote(ctx context.Context, in *pb.MessageToSign) (*pb.SignedMessage, error) {
  if in.Nr >= int32(len(s.queries)) || in.Nr < 0 { //security leak, path that out later!
    fmt.Printf("No such Query: %v\n", in.Nr)
	  return &pb.SignedMessage{Mess: []byte{}, Sign: []byte{}}, nil
  }
  var nT = len(s.queries[in.Nr].Tokens)
  if nT == 0 {
    fmt.Printf("Token not valid\n")
	  return &pb.SignedMessage{Mess: []byte{}, Sign: []byte{}}, nil
  }
  for i:=0; i<nT; i++ {
    if s.queries[in.Nr].Tokens[i].Token == in.Token.Token {
      s.queries[in.Nr].Tokens = append(s.queries[in.Nr].Tokens[:i], s.queries[in.Nr].Tokens[i+1:]...)
      break
    }
    if i == nT-1 {
      fmt.Printf("Token not valid\n")
	    return &pb.SignedMessage{Mess: []byte{}, Sign: []byte{}}, nil
    }
  }
  m := new(big.Int).SetBytes(in.Mess)
  sign := m.Exp(m, s.key.D, s.key.PublicKey.N).Bytes()
  var SM = pb.SignedMessage{
    Mess: in.Mess,
    Sign: sign,
  }
  fmt.Printf("Token valid\n")
  return &SM, nil
}

//for now tokens for votes dont have much informations, no crypto yet
func (s *server) QueryVote(ctx context.Context, in *pb.SignedVote) (*pb.VoteReply, error) {
  var votetohash = pb.Vote{
    Nr: in.Vote.Nr,
    Answer: in.Vote.Answer,
  }
  var hash = sha256.Sum256([]byte(fmt.Sprintf("%v", votetohash)))
  hi := new(big.Int).SetBytes(hash[:])
  bhi := hi.Exp(hi, s.key.D, s.key.PublicKey.N).Bytes()
  if bytes.Compare(bhi, in.Sign) != 0 { //check if sign is ok, rarher ok
    fmt.Printf("Sign invalid: %v\n", in.Sign)
	  return &pb.VoteReply{Mess: "Sign invalid!\n"}, nil
  }
  var v = in.Vote
  if v.Nr >= int32(len(s.queries)) || v.Nr < 0 { //security leak, path that out later!
    fmt.Printf("No such Query: %v\n", v.Nr)
	  return &pb.VoteReply{Mess: "No such Query!\n"}, nil
  }
  var nF = len(s.queries[v.Nr].Fields)
  if nF != len(v.Answer) {
    fmt.Printf("Vote have different number of fields than query\n")
	  return &pb.VoteReply{Mess: "Vote have different number of fields than query!\n"}, nil
  }
  for i:=0; i<nF; i++ {
    if v.Answer[i] >= 1 {
      s.queries[v.Nr].Fields[i].Votes++
    }
  }
  fmt.Printf("Thank you for your vote!\n")
  fmt.Printf("In Memory: %v\n", s.queries)
	return &pb.VoteReply{Mess: "Thank you for your vote!\n"}, nil
}

func serverInit() *server{
  var err error
  s := &server{}
  s.key, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("failed to generate key: %v", err)
	}
  s.key.Precompute()
  return s
}

func main() {
	rec, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterQueryServer(s, serverInit())
	s.Serve(rec)
}

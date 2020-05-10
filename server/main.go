package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
	"math/big"
	"net"

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
		N:    s.key.PublicKey.N.Bytes(),
		E:    int32(s.key.PublicKey.E),
	}, nil
}

func (s *server) QueryInit(stream pb.Query_QueryInitServer) error {
	var q = pb.PollQuestion{
		Id:     int32(len(s.queries)),
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
					Name:  field.Name,
					Votes: 0,
				})
		} else if field.Which < int32(len(q.Fields)) {
			q.Fields[field.Which].Name = field.Name
		} else {
			fmt.Printf("Wrong Query field number\n")
		}
	}
}

func (s *server) QueryGetToken(ctx context.Context, in *pb.TokenRequest) (*pb.VoteToken, error) {
	var t = pb.VoteToken{
		Token: -1}
	if in.Nr < 0 || in.Nr >= int32(len(s.queries)) {
		fmt.Printf("Wrong Query number in Get Token\n")
		return &t, nil
	}
	t.Token = int32(len(s.queries[in.Nr].Tokens) + 1)
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
	for i := 0; i < nT; i++ {
		if s.queries[in.Nr].Tokens[i].Token == in.Token.Token {
			s.queries[in.Nr].Tokens = append(s.queries[in.Nr].Tokens[:i], s.queries[in.Nr].Tokens[i+1:]...)
			break
		}
		if i == nT-1 {
			fmt.Printf("Token not valid\n")
			return &pb.SignedMessage{Mess: []byte{}, Sign: []byte{}}, nil
		}
	}
	//calculate m^d to generate sign for client
	m := new(big.Int).SetBytes(in.Mess)
	sign := m.Exp(m, s.key.D, s.key.PublicKey.N).Bytes()
	//sign := rsa.Decrypt(rand.Reader, s.key, m).Bytes()
	var SM = pb.SignedMessage{
		Mess: in.Mess, //may be not necessary
		Sign: sign,
	}
	fmt.Printf("Token valid\n")
	return &SM, nil
}

//get signed vote from client and check it's validity
func (s *server) QueryVote(ctx context.Context, in *pb.SignedVote) (*pb.VoteReply, error) {
	//first we check if sign is valid
	//to do this check if hash(signm) = (signmd)^e mod N
	//if sign is correct, signmd = hash(signm)^d and equality is satisfied
	hash := sha256.Sum256(in.Signm)
	//calculate signmd^e mod N
	md := new(big.Int).SetBytes(in.Signmd)
	bhi := new(big.Int).Exp(md, big.NewInt(int64(s.key.PublicKey.E)), s.key.PublicKey.N).Bytes()

	if bytes.Compare(bhi, hash[:]) != 0 { //check if sign is ok
		fmt.Printf("Sign invalid: (%v, %v)\n", in.Signm, in.Signmd)
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
	for i := 0; i < nF; i++ {
		if v.Answer[i] >= 1 {
			s.queries[v.Nr].Fields[i].Votes++
		}
	}
	fmt.Printf("Thank you for your vote!\n")
	fmt.Printf("In Memory: %v\n", s.queries)
	return &pb.VoteReply{Mess: "Thank you for your vote!\n"}, nil
}

func serverInit() *server {
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

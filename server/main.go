// Server package for Rada system.
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

	logic "github.com/ememak/Projekt-Rada/logic"
	pb "github.com/ememak/Projekt-Rada/query"
	"google.golang.org/grpc"
)

// In constants we store connection data.
const (
	port = ":12345"
)

// Server type contains server implemented in query/query.proto,
// data used for cryptography and usage of queries.
type server struct {
	pb.UnimplementedQueryServer

	// RSA key for signature validation
	key     *rsa.PrivateKey
	queries []pb.PollQuestion
}

// Hello is function used to exchange server public key.
//
// As an input function takes HelloRequest, which defined in query/query.proto.
func (s *server) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Mess: "Hello World!\n",
		N:    s.key.PublicKey.N.Bytes(),
		E:    int32(s.key.PublicKey.E),
	}, nil
}

// QueryInit generates new query.
//
// Questions are passed from client to server as stream of Field messages
// defined in query.proto.
func (s *server) QueryInit(stream pb.Query_QueryInitServer) error {
	// Make new Query
	q := pb.PollQuestion{
		Id:     int32(len(s.queries)),
		Fields: []*pb.PollQuestion_QueryField{},
	}
	fmt.Printf("QueryInitReceived, id = %v\n", len(s.queries))

	for {
		field, err := stream.Recv()
		// End of stream, we are saving new query.
		if err == io.EOF {
			s.queries = append(s.queries, q)
			fmt.Printf("Sending back: %v\n", q)
			fmt.Printf("In Memory: %v\n", s.queries)
			return stream.SendAndClose(&q)
		}
		if err != nil {
			return err
		}
		// Edit field in query, -1 is a signal of new field.
		if field.Which == -1 {
			q.Fields = append(q.Fields,
				&pb.PollQuestion_QueryField{
					Name:  field.Name,
					Votes: 0,
				})
		} else if field.Which < int32(len(q.Fields)) && field.Which >= 0 {
			q.Fields[field.Which].Name = field.Name
		} else {
			fmt.Printf("Wrong Query field number\n")
		}
	}
}

// QueryGetToken generates token used to authorize ballot.
func (s *server) QueryGetToken(ctx context.Context, in *pb.TokenRequest) (*pb.VoteToken, error) {
	t := pb.VoteToken{}
	if in.Nr < 0 || in.Nr >= int32(len(s.queries)) {
		fmt.Printf("Wrong Query number in Get Token\n")
		return &t, nil
	}
	t.Token = int32(len(s.queries[in.Nr].Tokens) + 1)
	s.queries[in.Nr].Tokens = append(s.queries[in.Nr].Tokens, &t)
	fmt.Printf("GetToken, in Memory: %v\n", s.queries)
	return &t, nil
}

// QueryAuthorizeVote authorizes a ballot if sent with valid token.
//
// Function takes as input message consisting of blinded ballot and
// token returned by function QueryGetToken if such token was not used before.
func (s *server) QueryAuthorizeVote(ctx context.Context, in *pb.MessageToSign) (*pb.SignedMessage, error) {
	// Check if token and number of query are valid.
	ok := logic.AcceptToken(in.Token, in.Nr, s.queries)
	if ok == false {
		return &pb.SignedMessage{Mess: []byte{}, Sign: []byte{}}, nil
	}

	// Calculate m^d to generate sign for client.
	m := new(big.Int).SetBytes(in.Mess)
	sign := m.Exp(m, s.key.D, s.key.PublicKey.N).Bytes()
	var SM = pb.SignedMessage{
		Mess: in.Mess, //may be not necessary
		Sign: sign,
	}
	fmt.Printf("Token valid\n")
	return &SM, nil
}

// QueryVote get signed vote from client and check it's validity.
func (s *server) QueryVote(ctx context.Context, in *pb.SignedVote) (*pb.VoteReply, error) {
	// First we have to check if the sign is valid.
	// To do this check if hash(signm) = (signmd)^e mod N.
	// If sign is correct, signmd = hash(signm)^d and equality is satisfied.
	hash := sha256.Sum256(in.Signm)
	// Calculate signmd^e mod N
	md := new(big.Int).SetBytes(in.Signmd)
	bhi := new(big.Int).Exp(md, big.NewInt(int64(s.key.PublicKey.E)), s.key.PublicKey.N).Bytes()

	if bytes.Compare(bhi, hash[:]) != 0 { //check if sign is ok
		fmt.Printf("Sign invalid: (%v, %v)\n", in.Signm, in.Signmd)
		return &pb.VoteReply{Mess: "Sign invalid!\n"}, nil
	}

	// Vote is properly signed, we proceed to voting.
	return logic.AcceptVote(in.Vote, s.queries)
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

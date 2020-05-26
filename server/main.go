// Server package for Rada system.
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"net"

	"github.com/ememak/Projekt-Rada/bsign"
	"github.com/ememak/Projekt-Rada/logic"
	"github.com/ememak/Projekt-Rada/query"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"
)

// In constants we store connection data.
const (
	port = ":12345"
)

// Server type contains server implemented in query/query.proto,
// data used for cryptography and usage of queries.
type server struct {
	query.UnimplementedQueryServer

	// RSA key for signature validation
	key     *rsa.PrivateKey
	data    *bolt.DB
	queries []query.PollQuestion
}

// Hello is function used to exchange server public key.
//
// As an input function takes HelloRequest, which defined in query/query.proto.
func (s *server) Hello(ctx context.Context, in *query.HelloRequest) (*query.HelloReply, error) {
	return &query.HelloReply{
		Mess: "Hello World!\n",
		N:    s.key.PublicKey.N.Bytes(),
		E:    int32(s.key.PublicKey.E),
	}, nil
}

// QueryInit generates new query.
//
// Questions are passed from client to server as stream of Field messages
// defined in query.proto.
func (s *server) QueryInit(stream query.Query_QueryInitServer) error {
	// Make new Query
	id, err := logic.NewQuery(s.data)
	if err != nil {
		fmt.Printf("failed to create new query in database: %v", err)
	}
	fmt.Printf("QueryInitReceived, id = %v\n", len(s.queries))

	for {
		field, err := stream.Recv()
		// End of stream, we are saving new query.
		if err == io.EOF {
			q := logic.GetQuery(s.data, id)
			// If everything will be stored in database, this is deprecated.
			s.queries = append(s.queries, q)
			fmt.Printf("Sending back: %v\n", q)
			fmt.Printf("In Memory: %v\n", s.queries)
			return stream.SendAndClose(&q)
		}
		if err != nil {
			return err
		}
		err = logic.ModifyQueryField(s.data, id, field.Which, field.Name)
	}
}

// QueryGetToken generates token used to authorize ballot.
func (s *server) QueryGetToken(ctx context.Context, in *query.TokenRequest) (*query.VoteToken, error) {
	t := query.VoteToken{}
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
func (s *server) QueryAuthorizeVote(ctx context.Context, in *query.MessageToSign) (*query.SignedMessage, error) {
	// Check if token and number of query are valid.
	ok := logic.AcceptToken(in.Token, in.Nr, s.queries)
	if ok == false {
		return &query.SignedMessage{}, nil
	}

	// Server is signing authorized message.
	sign := bsign.Sign(s.key, in.Mess)
	SM := query.SignedMessage{
		Mess: in.Mess, //may be not necessary
		Sign: sign.Bytes(),
	}
	fmt.Printf("Token valid\n")
	return &SM, nil
}

// QueryVote get signed vote from client and check it's validity.
func (s *server) QueryVote(ctx context.Context, in *query.SignedVote) (*query.VoteReply, error) {
	// We have to check if the sign is valid.
	if bsign.Verify(&s.key.PublicKey, in.Signm, in.Signmd) == false {
		fmt.Printf("Sign invalid!\n")
		return &query.VoteReply{Mess: "Sign invalid!\n"}, nil
	}

	// Vote is properly signed, we proceed to voting.
	return logic.AcceptVote(in.Vote, s.queries)
}

func serverInit() *server {
	var err error
	s := &server{}

	s.data, err = logic.DBInit("data.db")
	if err != nil {
		fmt.Printf("failed to initialise database: %v", err)
	}

	key := logic.GetKey(s.data)
	if key != nil {
		s.key = key
	} else {
		s.key, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			fmt.Printf("failed to generate key: %v", err)
		}
		err = logic.SaveKey(s.data, s.key)
		if err != nil {
			fmt.Printf("failed to save key in database: %v", err)
		}
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
	service := serverInit()
	query.RegisterQueryServer(s, service)

	defer service.data.Close()
	s.Serve(rec)
}

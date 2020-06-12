// Server package for Rada system.
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io"
	"net"

	"github.com/ememak/Projekt-Rada/bsign"
	"github.com/ememak/Projekt-Rada/query"
	"github.com/ememak/Projekt-Rada/store"
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

	data *bolt.DB
}

// KeyExchange is function used to exchange server public key for specific poll.
//
// As an input function takes KeyRequest, which contains number of query.
// If key is not in database (e.g. requested nonexisting query), reply contains empty byte array.
func (s *server) KeyExchange(ctx context.Context, in *query.KeyRequest) (*query.KeyReply, error) {
	key := store.GetKey(s.data, int(in.Nr))
	var binkey []byte
	if key != nil {
		binkey = x509.MarshalPKCS1PublicKey(&key.PublicKey)
	}

	return &query.KeyReply{
		Key: binkey,
	}, nil
}

// QueryInit generates new query.
//
// Questions are passed from client to server as stream of Field messages
// defined in query.proto.
func (s *server) QueryInit(stream query.Query_QueryInitServer) error {
	// Make new Query
	id, err := store.NewQuery(s.data)
	if err != nil {
		fmt.Printf("failed to create new query in database: %v", err)
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("failed to generate key: %v", err)
	}
	err = store.SaveKey(s.data, id, key)
	if err != nil {
		fmt.Printf("failed to save key in database: %v", err)
	}

	fmt.Printf("QueryInitReceived, id = %v\n", id)

	for {
		field, err := stream.Recv()
		// End of stream, we are saving new query.
		if err == io.EOF {
			q := store.GetQuery(s.data, id)
			fmt.Printf("Sending back: %v\n", q)
			return stream.SendAndClose(&q)
		}
		if err != nil {
			return err
		}
		err = store.ModifyQueryField(s.data, id, field.Which, field.Name)
	}
}

// QueryGetToken generates token used to authorize ballot.
func (s *server) QueryGetToken(ctx context.Context, in *query.TokenRequest) (*query.VoteToken, error) {
	return store.NewToken(s.data, in)
}

// QueryAuthorizeVote authorizes a ballot if sent with valid token.
//
// Function takes as input message consisting of blinded ballot and
// token returned by function QueryGetToken if such token was not used before.
func (s *server) QueryAuthorizeVote(ctx context.Context, in *query.MessageToSign) (*query.SignedMessage, error) {
	// Check if token and number of query are valid.
	ok, _ := store.AcceptToken(s.data, in.Token, in.Nr)
	if ok == false {
		return &query.SignedMessage{}, nil
	}

	// Server is signing authorized message.
	sign := bsign.Sign(store.GetKey(s.data, int(in.Nr)), in.Mess)
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
	if bsign.Verify(&store.GetKey(s.data, int(in.Vote.Nr)).PublicKey, in.Signm, in.Signmd) == false {
		fmt.Printf("Sign invalid!\n")
		return &query.VoteReply{Mess: "Sign invalid!\n"}, nil
	}

	// Vote is properly signed, we proceed to voting.
	vr, err := store.AcceptVote(s.data, in)
	if err != nil {
		err = fmt.Errorf("Error in database during voting: %w\n", err)
		return vr, err
	}
	fmt.Printf("In Memory: %v\n", store.GetQuery(s.data, int(in.Vote.Nr)))
	return vr, nil
}

func serverInit() *server {
	var err error
	s := &server{}

	s.data, err = store.DBInit("data.db")
	if err != nil {
		fmt.Printf("failed to initialise database: %v", err)
	}
	return s
}

func main() {
	rec, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	service := serverInit()
	defer service.data.Close()
	query.RegisterQueryServer(s, service)

	s.Serve(rec)
}

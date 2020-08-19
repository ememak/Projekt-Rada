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
	"os"

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
	key, err := store.GetKey(s.data, int(in.Pollid))
	if err != nil {
		err = fmt.Errorf("Error in KeyExchange while retrieving key from database: %w", err)
		return &query.KeyReply{}, err
	}
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
		err = fmt.Errorf("Error in QueryInit while creating new query in database: %w", err)
		return err
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		err = fmt.Errorf("Error in QueryInit during key generation: %w", err)
		return err
	}
	err = store.SaveKey(s.data, id, key)
	if err != nil {
		err = fmt.Errorf("Error in QueryInit while saving key: %w", err)
		return err
	}

	fmt.Printf("QueryInitReceived, id = %v\n", id)

	for {
		field, err := stream.Recv()
		// End of stream, we are saving new query.
		if err == io.EOF {
			q, errins := store.GetQuery(s.data, id)
			if errins != nil {
				return errins
			}

			fmt.Printf("Sending back: %v\n", q)
			return stream.SendAndClose(&q)
		}
		if err != nil {
			err = fmt.Errorf("Error in QueryInit while streaming: %w", err)
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
func (s *server) QueryAuthorizeVote(ctx context.Context, in *query.BallotToSign) (*query.SignedBallot, error) {
	// Check if token and number of query are valid.
	ok, _ := store.AcceptToken(s.data, in.Token, in.Pollid)
	if ok == false {
		return &query.SignedBallot{}, nil
	}

	key, err := store.GetKey(s.data, int(in.Pollid))
	if err != nil {
		fmt.Printf("Error in QueryAuthorizeVote while retrieving key from database: %w", err)
		return &query.SignedBallot{}, err
	}
	// Server is signing authorized message.
	sign := bsign.Sign(key, in.Ballot)
	SM := query.SignedBallot{
		Ballot: in.Ballot, //may be not necessary
		Sign:   sign.Bytes(),
	}
	fmt.Printf("Token valid\n")
	return &SM, nil
}

// QueryVote get signed vote from client, check it's validity and save it.
//
// SignedVote on input consists of vote and sign. If sign was used before, vote is overwritten.
func (s *server) QueryVote(ctx context.Context, in *query.SignedVote) (*query.VoteReply, error) {
	key, err := store.GetKey(s.data, int(in.Vote.Pollid))
	if err != nil {
		err = fmt.Errorf("Error in QueryVote while retrieving key from database: %w", err)
		return &query.VoteReply{Mess: "Error in QueryVote\n"}, err
	}
	// We have to check if the sign is valid.
	if bsign.Verify(&key.PublicKey, in.Signm, in.Signmd) == false {
		return &query.VoteReply{Mess: "Error in QueryVote\n"}, fmt.Errorf("Sign invalid!")
	}

	// Vote is properly signed, we proceed to voting.
	vr, err := store.AcceptVote(s.data, in)
	if err != nil {
		err = fmt.Errorf("Error in QueryVote while saving key in database: %w", err)
		return vr, err
	}

	q, _ := store.GetQuery(s.data, int(in.Vote.Pollid))
	fmt.Printf("In Memory: %v\n", q)
	return vr, nil
}

func serverInit() (*server, error) {
	var err error
	s := &server{}

	s.data, err = store.DBInit("data.db")
	if err != nil {
		err = fmt.Errorf("Error in serverInit, failed to initialise database: %w", err)
	}
	return s, err
}

func main() {
	rec, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Server failed to listen: %v", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	service, err := serverInit()
	if err != nil {
		fmt.Printf("", err)
		os.Exit(1)
	}

	defer service.data.Close()
	query.RegisterQueryServer(s, service)

	s.Serve(rec)
}

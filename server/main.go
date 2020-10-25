// Server package for Rada system.
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	//"github.com/ememak/Projekt-Rada/bsign"
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

// GetPoll is function used to exchange server public key for specific poll.
//
// As an input function takes GetPollRequest, which contains poll's number.
// If key or poll are not in database (e.g. requested nonexisting query), reply contains empty answer.
func (s *server) GetPoll(ctx context.Context, in *query.GetPollRequest) (*query.PollWithPublicKey, error) {
	key, err := store.GetKey(s.data, int(in.Pollid))
	if err != nil {
		err = fmt.Errorf("Error in GetPoll while retrieving key from database: %w", err)
		return &query.PollWithPublicKey{}, err
	}
	binkey := x509.MarshalPKCS1PublicKey(&key.PublicKey)

	poll, err := store.GetPoll(s.data, int(in.Pollid))
	if err != nil {
		err = fmt.Errorf("Error in GetPoll while retrieving poll from database: %w", err)
		return &query.PollWithPublicKey{}, err
	}

	return &query.PollWithPublicKey{
		Key: &query.PublicKey{
			Key: binkey,
		},
		Poll: poll.Schema,
	}, nil
}

// PollInit generates new poll and saves it to database.
//
// Questions and their types are passed in input parameter.
func (s *server) PollInit(ctx context.Context, in *query.PollSchema) (*query.PollQuestion, error) {
	poll, err := store.NewPoll(s.data, in)
	if err != nil {
		return poll, fmt.Errorf("Error in PollInit while creating new poll in database: %w", err)
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return poll, fmt.Errorf("Error in PollInit during key generation: %w", err)
	}
	err = store.SaveKey(s.data, int(poll.Id), key)
	if err != nil {
		return poll, fmt.Errorf("Error in PollInit while saving key: %w", err)
	}

	fmt.Printf("PollInitReceived, id = %v\n", poll.Id)
	return poll, nil
}

/*
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
}*/

func serverInit(dbfilename string) (*server, error) {
	var err error
	s := &server{}

	s.data, err = store.DBInit(dbfilename)
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
	service, err := serverInit("data.db")
	if err != nil {
		fmt.Printf("Error in serverInit: %v", err)
		os.Exit(1)
	}

	defer service.data.Close()
	query.RegisterQueryServer(s, service)

	s.Serve(rec)
}

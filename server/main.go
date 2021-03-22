// Server package for Rada system.
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ememak/Projekt-Rada/bsign"
	"github.com/ememak/Projekt-Rada/query"
	"github.com/ememak/Projekt-Rada/store"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// In constants we store connection data.
/*const (
	//port = ":12345"
	port = (os.Getenv("PORT") ? (":" + os.Getenv("PORT")) : (":12345")),
)*/

// Server type contains server implemented in query/query.proto,
// data used for cryptography and usage of polls.
type server struct {
	query.UnimplementedQueryServer

	data *bolt.DB
}

// GetPoll is function used to exchange server public key for specific poll.
//
// GetPollRequest contains poll's id. This poll will be returned.
// If key or poll are not in database (e.g. requested nonexisting poll), reply contains empty answer.
func (s *server) GetPoll(ctx context.Context, in *query.GetPollRequest) (*query.PollWithPublicKey, error) {
	key, err := store.GetKey(s.data, in.Pollid)
	if err != nil {
		err = fmt.Errorf("Error in GetPoll while retrieving key from database: %w", err)
		return &query.PollWithPublicKey{}, err
	}
	binkey := x509.MarshalPKCS1PublicKey(&key.PublicKey)

	poll, err := store.GetPoll(s.data, in.Pollid)
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

	return poll, nil
}

// SignBallot authorizes a ballot if sent with valid token.
//
// Function takes as an input message consisting of an envelope (blinded ballot)
// and a token. Envelope is signed if token is valid.
func (s *server) SignBallot(ctx context.Context, in *query.EnvelopeToSign) (*query.SignedEnvelope, error) {
	// Check if token and polls number are valid.
	err := store.AcceptToken(s.data, in.Token, in.Pollid)
	if err != nil {
		return &query.SignedEnvelope{}, err
	}

	key, err := store.GetKey(s.data, in.Pollid)
	if err != nil {
		err = fmt.Errorf("Error in SignBallot while retrieving key from database: %w", err)
		return &query.SignedEnvelope{}, err
	}
	// Token is valid. Server is signing envelope.
	sign := bsign.Sign(key, in.Envelope)
	if len(sign.Bytes()) == 0 {
		return &query.SignedEnvelope{}, fmt.Errorf("Error in SignBallot, envelope shouldn't be null")
	}
	SM := query.SignedEnvelope{
		Envelope: in.Envelope, //may be not necessary
		Sign:     sign.Bytes(),
	}
	return &SM, nil
}

// PollVote get signed vote from client, check it's validity and save it.
//
// VoteRequest on input consists of vote and sign. If sign was used before, vote is overwritten.
func (s *server) PollVote(ctx context.Context, in *query.VoteRequest) (*query.VoteReply, error) {
	key, err := store.GetKey(s.data, in.Pollid)
	if err != nil {
		err = fmt.Errorf("Error in PollVote while retrieving key from database: %w", err)
		return &query.VoteReply{Mess: "Error in PollVote"}, err
	}
	// We have to check if the sign is valid.
	if bsign.Verify(&key.PublicKey, in.Sign.Ballot, in.Sign.Sign) == false {
		err = fmt.Errorf("Error in PollVte, Sign invalid!")
		return &query.VoteReply{Mess: "Error in PollVote"}, err
	}

	// Vote is properly signed, we proceed to voting.
	vr, err := store.SaveVote(s.data, in)
	if err != nil {
		err = fmt.Errorf("Error in PollVote while saving key in database: %w", err)
		return &query.VoteReply{Mess: "Error in PollVote"}, err
	}

	return vr, nil
}

// GetSummary sends all answers for a poll.
func (s *server) GetSummary(ctx context.Context, in *query.SummaryRequest) (*query.PollSummary, error) {
	return store.GetSummary(s.data, in.Pollid)
}

func serverInit(dbfilename string) (*server, error) {
	var err error
	s := &server{}

	s.data, err = store.DBInit(dbfilename)
	if err != nil {
		err = fmt.Errorf("Error in serverInit, failed to initialise database: %w", err)
	}
	return s, err
}

func stringContainSomeElement(s string, match []string) bool {
	for _, v := range match {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}

func main() {
	port := ""
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	} else {
		port = ":12345"
	}
	s := grpc.NewServer()
	service, err := serverInit("data.db")
	if err != nil {
		fmt.Printf("Error in serverInit: %v", err)
		os.Exit(1)
	}

	defer service.data.Close()
	query.RegisterQueryServer(s, service)
	grpclog.SetLogger(log.New(os.Stdout, "exampleserver: ", log.LstdFlags))

	wrappedGrpc := grpcweb.WrapServer(s)
	h := http.FileServer(http.Dir("./client/prodapp"))
	handler := func(resp http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.URL.Path, "query.Query") {
			wrappedGrpc.ServeHTTP(resp, req)
		} else {
			subpages := []string{"pollinit", "vote", "results"}
			if stringContainSomeElement(req.URL.Path, subpages) {
				req.URL.Path = "/"
			}
			h.ServeHTTP(resp, req)
		}
	}
	httpServer := http.Server{
		Addr:    port,
		Handler: http.HandlerFunc(handler),
	}
	fmt.Printf("Server listening on http://localhost" + port + "\n")
	err = httpServer.ListenAndServe()
	if err != nil {
		fmt.Printf("Error while launching server: %v\n", err)
		os.Exit(1)
	}
}

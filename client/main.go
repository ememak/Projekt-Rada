// Example client for Rada system.
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"math/big"
	"time"

	query "github.com/ememak/Projekt-Rada/query"
	"google.golang.org/grpc"
)

var (
	addr = "localhost:12345"
)

func runKeyExchange(ctx context.Context, client query.QueryClient, queryid int32) *rsa.PublicKey {
	r, err := client.KeyExchange(ctx, &query.KeyRequest{Nr: queryid})
	if err != nil {
		fmt.Printf("Client got error on KeyExchange function: %v\n", err)
	}
	key, err := x509.ParsePKCS1PublicKey(r.Key)
	if err != nil {
		fmt.Printf("Error in parsing key: %v\n", err)
	}
	return key
}

func runQueryInit(ctx context.Context, client query.QueryClient) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	stream, err := client.QueryInit(ctx)
	if err != nil {
		log.Fatalf("%v.QueryInit(_) = _, %v", client, err)
	}

	for i := 0; i < len(exampleQuery); i++ {
		if err := stream.Send(&exampleQuery[i]); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("%v.Send(%v) = %v", stream, exampleQuery[i], err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v", stream, err)
	}
	log.Printf("Query: %v", reply)
}

// RunVote performs all necessary measures to send valid vote.
//
// Input consists of client previously connected to query server
// and vote defined in query/query.proto.
//
// Function have to be called after key exchange.
//
// Inside function is getting token from server and anonymously sending vote
// signed using RSA blind signature scheme.
func runVote(ctx context.Context, client query.QueryClient, vote query.Vote) {
	key := runKeyExchange(ctx, client, vote.Nr)
	if key == nil {
		fmt.Printf("Failed to exchange keys\n")
		return
	}

	// Get token to vote, may be changed.
	t, err := client.QueryGetToken(ctx, &query.TokenRequest{Nr: vote.Nr})
	if err != nil {
		fmt.Printf("Client got error on GetToken function: %v", err)
	}
	fmt.Printf("Token: %v\n", t)

	// Generate ballot to be signed.
	// This value will be referred as Mess.
	mess, err := rand.Int(rand.Reader, key.N)
	if err != nil {
		fmt.Printf("Rand.Int error in generating message to sign: %v", err)
	}
	fmt.Printf("m: %v\n", mess)

	// We are hashing ballot.
	// From now on m will be name of this hashed message.
	hash := sha256.Sum256(mess.Bytes())
	m := new(big.Int).SetBytes(hash[:])
	fmt.Printf("Hashed m: %v\n", m)

	// Get random blinding factor.
	r, err := rand.Int(rand.Reader, key.N)
	if err != nil {
		fmt.Printf("Rand.Int error in generating blinding factor: %v", err)
	}

	// We want to send m*r^e mod N to server.
	// To do this we need to calculate few values
	// bfactor = r^e mod N
	bfactor := new(big.Int).Exp(r, big.NewInt(int64(key.E)), key.N)
	// blinded = m*(r^e) mod N
	blinded := bfactor.Mod(bfactor.Mul(bfactor, m), key.N)
	// Now we can send m*(r^e) to server.
	// We are sending it with number of query and token to it so that server could authorize our request.
	var mts = query.MessageToSign{
		Mess:  blinded.Bytes(),
		Nr:    vote.Nr,
		Token: t,
	}

	// Receive (m^d)*r mod N from server.
	sm, err := client.QueryAuthorizeVote(ctx, &mts)
	if err != nil {
		fmt.Printf("Client got error on QueryVote function: %v", err)
	}

	// Having (m^d)*r mod N we are removing blinding factor r,
	// then we send vote with actual sign which is pair (Mess, m^d mod N).
	// To get m^d mod N we need to perform few calculations:
	// smi = (m^d)*r mod N
	smi := new(big.Int).SetBytes(sm.Sign)
	// revr = r^-1 mod N
	revr := new(big.Int).ModInverse(r, key.N)
	// smirevr = smi * revr
	smirevr := new(big.Int).Mul(revr, smi)
	// Now we can calculate second part of sign.
	// sign = smirevr mod N = m^d mod N
	sign := new(big.Int).Mod(smirevr, key.N)

	// We are sending vote with pair (Mess, m^d mod N = hash(Mess)^d mod N)
	var sv = query.SignedVote{
		Vote:   &vote,
		Signm:  mess.Bytes(),
		Signmd: sign.Bytes(),
	}
	vr, err := client.QueryVote(ctx, &sv)
	if err != nil {
		fmt.Printf("Client got error on QueryVote function: %v", err)
	}
	fmt.Printf("Mess: %v\n", vr.Mess)
}

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("grpc.Dial got error %v", err)
	}
	defer conn.Close()
	c := query.NewQueryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	runQueryInit(ctx, c)

	runVote(ctx, c, exampleVote0)

	runVote(ctx, c, exampleVote1)
}

var exampleQuery = []query.Field{
	{Which: -1, Name: "First Option"},
	{Which: -1, Name: "Second Option"},
	{Which: 1, Name: "Edit First Option"},
	{Which: -1, Name: "Third Option"},
}

// First time launching the client first vote should pass, second not.
// Second time both should pass.
// Second vote is asking about query number one, which don't exist during first launch.
var exampleVote0 = query.Vote{
	Nr:     1,
	Answer: []int32{0, 1, 1},
}

var exampleVote1 = query.Vote{
	Nr:     2,
	Answer: []int32{1, 1, 0},
}

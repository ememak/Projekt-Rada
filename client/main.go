package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math/big"
	"time"

	pb "github.com/ememak/Projekt-Rada/query"
	"google.golang.org/grpc"
)

var (
	addr = "localhost:12345"
)

var key *rsa.PublicKey

func runHello(client pb.QueryClient, ctx context.Context) {
	r, err := client.Hello(ctx, &pb.HelloRequest{})
	if err != nil {
		fmt.Printf("Client got error on Hello function: %v", err)
	}
	key = &rsa.PublicKey{
		N: new(big.Int).SetBytes(r.GetN()),
		E: int(r.GetE()),
	}
}

func runQueryInit(client pb.QueryClient, ctx context.Context) {
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

//Function to get token, generate blind signature and vote using it
func runVote(client pb.QueryClient, ctx context.Context, vote pb.Vote) {
	//get token to vote, may be changed
	t, err := client.QueryGetToken(ctx, &pb.TokenRequest{Nr: vote.Nr})
	if err != nil {
		fmt.Printf("Client got error on GetToken function: %v", err)
	}
	fmt.Printf("Token: %v\n", t)

	//generate message to be signed
	mess, err := rand.Int(rand.Reader, key.N)
	if err != nil {
		fmt.Printf("Rand.Int error in generating message to sign: %v", err)
	}
	fmt.Printf("m: %v\n", mess)
	//hash this message
	//from now on m is this hashed message
	hash := sha256.Sum256(mess.Bytes())
	m := new(big.Int).SetBytes(hash[:])
	fmt.Printf("Hashed m: %v\n", m)

	//get random blinding factor
	r, err := rand.Int(rand.Reader, key.N)
	if err != nil {
		fmt.Printf("Rand.Int error in generating blinding factor: %v", err)
	}

	//we want to send m*r^e mod N to server
	//r^e
	bfactor := new(big.Int).Exp(r, big.NewInt(int64(key.E)), key.N)
	//m*r^e mod N
	blinded := bfactor.Mod(bfactor.Mul(bfactor, m), key.N)
	//send m*r^e to server
	//with number of query and token to it
	var mts = pb.MessageToSign{
		Mess:  blinded.Bytes(),
		Nr:    vote.Nr,
		Token: t,
	}

	//receive m^d*r from server
	sm, err := client.QueryAuthorizeVote(ctx, &mts)
	if err != nil {
		fmt.Printf("Client got error on QueryVote function: %v", err)
	}

	//having m^d*r we are removing blinding factor r
	//then we send vote with actual sign
	//smi = m^d*r
	smi := new(big.Int).SetBytes(sm.Sign)
	//revr = r^-1 mod N
	revr := new(big.Int).ModInverse(r, key.N)
	//smirevr = smi * revr
	smirevr := new(big.Int).Mul(revr, smi)
	//sign = smirevr mod N = m^d mod N
	//actual sign
	sign := new(big.Int).Mod(smirevr, key.N).Bytes()

	//send vote with pair m, hash(m)^d
	var sv = pb.SignedVote{
		Vote:   &vote,
		Signm:  mess.Bytes(),
		Signmd: sign,
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
	c := pb.NewQueryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	runHello(c, ctx)

	runQueryInit(c, ctx)

	runVote(c, ctx, exampleVote0)

	runVote(c, ctx, exampleVote1)
}

var exampleQuery = []pb.Field{
	{Which: -1, Name: "First Option"},
	{Which: -1, Name: "Second Option"},
	{Which: 0, Name: "Edit First Option"},
	{Which: -1, Name: "Third Option"},
}

//first time launching the client first vote should pass, second not
//second time both should pass
//second vote is asking about query number one, which don't exist during first launch
var exampleVote0 = pb.Vote{
	Nr:     0,
	Answer: []int32{0, 1, 1},
}

var exampleVote1 = pb.Vote{
	Nr:     1,
	Answer: []int32{1, 1, 0},
}

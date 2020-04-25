package main

import (
  "context"
  "time"
  "fmt"
  "log"
  "io"
  "crypto/rsa"
  "crypto/rand"
  "crypto/sha256"
  "math/big"

  "google.golang.org/grpc"
  pb "github.com/ememak/Projekt-Rada/query"
)

var (
  addr = "localhost:12345"
)

var key *rsa.PublicKey

func runQueryInit (client pb.QueryClient, ctx context.Context) {
  ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
  defer cancel()
  stream, err := client.QueryInit(ctx)
  if err != nil {
    log.Fatalf("%v.QueryInit(_) = _, %v", client, err)
  }

  for i := 0; i<len(exampleQuery); i++ {
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

func runHello(client pb.QueryClient, ctx context.Context){
  r, err := client.Hello(ctx, &pb.HelloRequest{})
  if err != nil {
    fmt.Printf("Client got error on Hello function: %v", err)
  }
  key = &rsa.PublicKey{
    N: new(big.Int).SetBytes(r.GetN()),
    E: int(r.GetE()),
  }
}

func runVote(client pb.QueryClient, ctx context.Context, vote pb.Vote){
  t, err := client.QueryGetToken(ctx, &pb.TokenRequest{Nr: vote.Nr})
  if err != nil {
    fmt.Printf("Client got error on GetToken function: %v", err)
  }
  fmt.Printf("Token: %v\n", t)

  var hash = sha256.Sum256([]byte(fmt.Sprintf("%v", vote)))
  r, err := rand.Int(rand.Reader, key.N);
  if err != nil {
    fmt.Printf("Rand.Int error: %v", err)
  }
  c := new(big.Int).SetBytes(hash[:])
  bfactor := new(big.Int).Exp(r, big.NewInt(int64(key.E)), key.N)
  blinded := bfactor.Mod(bfactor.Mul(bfactor, c), key.N).Bytes()
  var mts = pb.MessageToSign{
    Mess: blinded,
    Nr: vote.Nr,
    Token: t,
  }
  sm, err := client.QueryAuthorizeVote(ctx, &mts)
  if err != nil {
    fmt.Printf("Client got error on QueryVote function: %v", err)
  }

  smi := new(big.Int).SetBytes(sm.Sign)
  revr := new(big.Int).ModInverse(r, key.N)
  smirevr := new(big.Int).Mul(revr, smi)
  sign := new(big.Int).Mod(smirevr, key.N).Bytes()
  var sv = pb.SignedVote{
    Vote: &vote,
    Sign: sign,
  }
  vr, err := client.QueryVote(ctx, &sv)
  if err != nil {
    fmt.Printf("Client got error on QueryVote function: %v", err)
  }
  fmt.Printf("Mess: %v\n", vr.Mess)
}

func main(){
  conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
  if err != nil {
    log.Fatalf("grpc.Dial got error %v", err)
  }
  defer conn.Close()
  c:=pb.NewQueryClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)

  defer cancel()

  runHello(c, ctx)

  runQueryInit(c, ctx);

  runVote(c, ctx, exampleVote0)

  runVote(c, ctx, exampleVote1)
}

var exampleQuery = []pb.Field{
    {Which: -1, Name: "First Option"},
    {Which: -1, Name: "Second Option"},
    {Which: 0, Name: "Edit First Option"},
    {Which: -1, Name: "Third Option"},
  }

var exampleVote0 = pb.Vote{
    Nr: 0,
    Answer: []int32{0, 1, 1},
  }

var exampleVote1 = pb.Vote{
    Nr: 1,
    Answer: []int32{1, 1, 0},
  }

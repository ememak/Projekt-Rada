package main

import (
  "context"
  "time"
  "fmt"
  "log"
  "io"

  "google.golang.org/grpc"
  pb "github.com/ememak/Projekt-Rada/query"
)

var (
  addr = "localhost:12345"
)

func runQueryInit (client pb.QueryClient) {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func main(){
  conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
  if err != nil {
    log.Fatalf("grpc.Dial got error %v", err)
  }
  defer conn.Close()
  c:=pb.NewQueryClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)

  defer cancel()
  r, err := c.Hello(ctx, &pb.HelloRequest{})
  if err != nil {
    fmt.Printf("Client got error on Hello function: %v", err)
  }

  fmt.Printf("Greeting: %s", r.GetMess())

  runQueryInit(c);

  t0, err := c.QueryGetToken(ctx, &pb.TokenRequest{})
  if err != nil {
    fmt.Printf("Client got error on GetToken function: %v", err)
  }
  exampleVote0.Token = t0
  fmt.Printf("Token: %v", t0)

  v0, err := c.QueryVote(ctx, &exampleVote0)
  if err != nil {
    fmt.Printf("Client got error on QueryVote function: %v", err)
  }
  fmt.Printf(v0.Mess)

  t1, err := c.QueryGetToken(ctx, &pb.TokenRequest{})
  if err != nil {
    fmt.Printf("Client got error on GetToken function: %v", err)
  }
  exampleVote1.Token = t1
  fmt.Printf("Token: %v", t1)

  v1, err := c.QueryVote(ctx, &exampleVote1)
  if err != nil {
    fmt.Printf("Client got error on QueryVote function: %v", err)
  }
  fmt.Printf(v1.Mess)
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
    Token: &pb.VoteToken{Token:-1,},
  }

var exampleVote1 = pb.Vote{
    Nr: 1,
    Answer: []int32{1, 1, 0},
    Token: &pb.VoteToken{Token:-1,},
  }

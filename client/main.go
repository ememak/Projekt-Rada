package main

import (
  "context"
  "time"
  "fmt"
  "log"
  "io"

  "google.golang.org/grpc"
  pb "github.com/ememak/Projekt-Rada/hello"
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
	  log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
  }
  log.Printf("Query: %v", reply)
}

func main(){
  conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
  if err != nil {
    //do sth
  }
  defer conn.Close()
  c:=pb.NewQueryClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)

  defer cancel()
  r, err := c.Hello(ctx, &pb.HelloRequest{})

  fmt.Printf("Greeting: %s", r.GetMess())

  runQueryInit(c);
}

var exampleQuery = []pb.Field{
    {Which: -1, Name: "First Option"},
    {Which: -1, Name: "Second Option"},
    {Which: 0, Name: "Edit First Option"},
    {Which: -1, Name: "Third Option"},
  }

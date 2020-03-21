package main

import (
  "context"
  "time"
  "fmt"

  "google.golang.org/grpc"
  pb "Projekt-Rada/hello"
)

var (
  addr = "localhost:12345"
)

func main(){
  conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
  if err != nil {
    //do sth
  }
  defer conn.Close()
  c:=pb.NewHelloClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)

  defer cancel()
  r, err := c.Hello(ctx, &pb.HelloRequest{})

  fmt.Printf("Greeting: %s", r.GetMess())
}

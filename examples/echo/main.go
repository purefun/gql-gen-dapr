package main

import (
	"context"
	"flag"
	"fmt"
)

func main() {
	runClient := flag.Bool("client", false, "run client")
	runServer := flag.Bool("server", false, "run server")
	flag.Parse()

	if !*runClient && !*runServer {
		panic("please add --client or --server flag to run the demo")
	}

	if *runClient {
		NewClient()
	}
	if *runServer {
		NewServer()
	}
}

func NewClient() {
	echo, err := NewEchoClient("echo_server")
	if err != nil {
		panic(err)
	}
	resp, err := echo.Text(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("echo.Text returns:", *resp)
}

func NewServer() {
	h := new(Handlers)
	s, err := NewEchoServer(":6000", h)
	if err != nil {
		panic(err)
	}
	s.Start()
}

type Handlers struct {
}

func (h *Handlers) Echo(ctx context.Context, in *EchoInput) (*EchoOutput, error) {
	panic("not implemented") // TODO: Implement
}

func (h *Handlers) Text(ctx context.Context) (*string, error) {
	resp := "hello world"
	return &resp, nil
}

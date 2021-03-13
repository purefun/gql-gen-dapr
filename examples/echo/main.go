package main

import (
	"context"
	"flag"
	"fmt"

	DaprClient "github.com/dapr/go-sdk/client"
	DaprServer "github.com/dapr/go-sdk/service/grpc"
)

func main() {
	runClient := flag.Bool("client", false, "run client")
	runServer := flag.Bool("server", false, "run server")
	flag.Parse()

	if *runClient {
		NewClient()
	}
	if *runServer {
		NewServer()
	}
}

func NewClient() {
	c, _ := DaprClient.NewClient()
	echo := NewEchoClient(c, "echo_server")
	resp, err := echo.Text(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("echo.Text returns:", *resp)
}

func NewServer() {
	s, err := DaprServer.NewService(":6000")
	if err != nil {
		panic(err)
	}

	h := new(Handlers)

	Register(s, h)

	err = s.Start()
	if err != nil {
		panic(err)
	}
}

type Handlers struct {
}

func (h *Handlers) Echo(ctx context.Context) (*EchoOutput, error) {
	panic("not implemented") // TODO: Implement
}

func (h *Handlers) Text(ctx context.Context) (*string, error) {
	resp := "hello world"
	return &resp, nil
}

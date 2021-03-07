package main

import (
	"context"
	"flag"
	"fmt"

	daprClient "github.com/dapr/go-sdk/client"
	daprServer "github.com/dapr/go-sdk/service/grpc"
	example "github.com/purefun/gql-gen-dapr/generator/testdata"
)

// dapr run -a example_server -p 3000 -P grpc -- go run main.go --server
// dapr run -a example_client -p 3000 -P grpc -- go run main.go --client

func main() {

	runClient := flag.Bool("client", false, "run client")
	runServer := flag.Bool("server", false, "run server")

	flag.Parse()

	if *runClient {
		client()
	}

	if *runServer {
		server()
	}
}

func client() {
	c, _ := daprClient.NewClient()
	client := example.NewExampleClient(c, "example_server")

	resp, err := client.Hello(context.Background())

	if err != nil {
		fmt.Println("call echo error: ", err)
	} else {
		fmt.Println("Hello returns => ", resp)
	}
}

type Handlers struct {
}

func (h *Handlers) Hello(ctx context.Context) (string, error) {
	return "Hello, World!", nil
}

func server() {
	s, err := daprServer.NewService(":3000")

	if err != nil {
		panic(err)
	}

	h := new(Handlers)

	example.Register(s, h)

	s.Start()
}

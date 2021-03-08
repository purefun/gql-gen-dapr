package testdata

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
)

type Example interface {
	Hello(ctx context.Context) (string, error)
}

type _ExampleClient struct {
	cc    client.Client
	appID string
}

func NewExampleClient(cc client.Client, appID string) *_ExampleClient {
	return &_ExampleClient{cc, appID}
}

func (c *_ExampleClient) Hello(ctx context.Context) (string, error) {
	content := &client.DataContent{ContentType: "application/json"}
	resp, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "Hello", "post", content)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

type InvocationHandlerFunc func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)

func _Example_Hello_Handler(srv Example) InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		resp, mErr := srv.Hello(ctx)
		if mErr != nil {
			err = mErr
			return
		}
		data, encErr := json.Marshal(resp)
		if encErr != nil {
			err = encErr
			return
		}
		out = &common.Content{
			ContentType: "application/json",
			Data:        data,
		}
		return
	}
}

func Register(s common.Service, srv ExampleServer) {
	s.AddServiceInvocationHandler("Hello", _Example_Hello_Handler(srv))
}
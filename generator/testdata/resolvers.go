package testdata

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
)

type Example interface {
	Hello(ctx context.Context) (*string, error)
	Hey(ctx context.Context, in *string) (*string, error)
}

type _ExampleClient struct {
	cc    client.Client
	appID string
}

func NewExampleClient(cc client.Client, appID string) *_ExampleClient {
	return &_ExampleClient{cc, appID}
}

func (c *_ExampleClient) Hello(ctx context.Context) (*string, error) {
	content := &client.DataContent{ContentType: "application/json"}
	resp, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "hello", "post", content)
	if err != nil {
		return nil, err
	}
	out := new(string)
	err = json.Unmarshal(resp, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *_ExampleClient) Hey(ctx context.Context, in *string) (*string, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	content := &client.DataContent{ContentType: "application/json", Data: data}
	resp, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "hey", "post", content)
	if err != nil {
		return nil, err
	}
	out := new(string)
	err = json.Unmarshal(resp, out)
	if err != nil {
		return nil, err
	}
	return out, nil
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

func _Example_Hey_Handler(srv Example) InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		req := new(string)
		reqErr := json.Unmarshal(in.Data, req)
		if reqErr != nil {
			return nil, err
		}
		resp, mErr := srv.Hey(ctx, req)
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

func Register(s common.Service, srv Example) {
	s.AddServiceInvocationHandler("hello", _Example_Hello_Handler(srv))
	s.AddServiceInvocationHandler("hey", _Example_Hey_Handler(srv))
}

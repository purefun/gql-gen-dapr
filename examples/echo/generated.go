package echo

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
)

type EchoInput struct {
	Text string `json:"text"`
}

type EchoOutput struct {
	Text string `json:"text"`
}

type Echo interface {
	Echo(ctx context.Context) (EchoOutput, error)
	Text(ctx context.Context) (string, error)
}

type _EchoClient struct {
	cc    client.Client
	appID string
}

func NewEchoClient(cc client.Client, appID string) *_EchoClient {
	return &_EchoClient{cc, appID}
}

func (c *_EchoClient) Echo(ctx context.Context, in *EchoInput) (*EchoOutput, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	content := &client.DataContent{ContentType: "application/json", Data: data}
	resp, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "Echo", "post", content)
	if err != nil {
		return nil, err
	}
	out := new(EchoOutput)
	err = json.Unmarshal(resp, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *_EchoClient) Text(ctx context.Context) (*string, error) {
	content := &client.DataContent{ContentType: "application/json"}
	resp, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "Text", "post", content)
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

func _Echo_Echo_Handler(srv Echo) InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		resp, mErr := srv.Echo(ctx)
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

func _Echo_Text_Handler(srv Echo) InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		resp, mErr := srv.Text(ctx)
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

func Register(s common.Service, srv Echo) {
	s.AddServiceInvocationHandler("echo", _Echo_Echo_Handler(srv))
	s.AddServiceInvocationHandler("text", _Echo_Text_Handler(srv))
}

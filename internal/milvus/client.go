package milvus

import (
	"context"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type Client struct {
	MilvusClient client.Client
}

func NewClient(milvusURI string) *Client {
	c, err := client.NewClient(context.Background(), client.Config{
		Address: "localhost:19530",
	})
	if err != nil {
		// handle error
	}
	defer c.Close()

	return &Client{MilvusClient: c}
}

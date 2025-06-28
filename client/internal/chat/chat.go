package chat

import (
	"context"
	pb "github.com/Git5737/lexchanger/proto/chat/proto"
	"google.golang.org/grpc"
)

type Client struct {
	stream pb.ChatServise_EventStreamClient
	conn   *grpc.ClientConn
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewChatServiseClient(conn)
	stream, err := client.EventStream(context.Background())
	if err != nil {
		return nil, err
	}

	return &Client{stream: stream, conn: conn}, nil
}

func (c *Client) Send(event *pb.Events) error {
	return c.stream.Send(event)
}

func (c *Client) Recv() (*pb.Events, error) {
	return c.stream.Recv()
}

func (c *Client) Close() {
	c.conn.Close()
}

package pubsub

import (
	"context"
	"errors"
	"fmt"
	"os"

	pb "github.com/structx/common/pkg/pubsub"
	"google.golang.org/grpc"
)

// Client
type Client struct {
	address string
	ch      chan *pb.Message
	stream  pb.PubSubService_PublishClient
}

// New returns a new pubsub client
func New() (*Client, error) {

	aa := os.Getenv("PUBSUB_SERVER_ADDRESS")
	if aa == "" {
		return nil, errors.New("PUBSUB_SERVER_ADDRESS is not set")
	}

	return &Client{
		address: aa,
		ch:      make(chan *pb.Message),
	}, nil
}

// Monitor monitors the channel for messages
func (c *Client) Monitor(ctx context.Context) error {

	// connect to server
	cl, co, err := newConn(ctx, c.address)
	if err != nil {
		return err
	}
	defer co.Close()

	// create publish stream
	c.stream, err = cl.Publish(ctx)
	if err != nil {
		return fmt.Errorf("failed to create publish stream: %w", err)
	}

	// start background worker thread
	for i := 0; i < 1; i++ {
		go c.publish(ctx, cl)
	}

	return nil
}

// Publish publishes a message
func (c *Client) Publish(ctx context.Context, topic string, data []byte) error {

	// add message to channel
	c.ch <- &pb.Message{
		Topic: topic,
		Data:  data,
	}

	return nil
}

// Close closes the client publication channel
func (c *Client) Close() {
	close(c.ch)
}

func (c *Client) publish(ctx context.Context, client pb.PubSubServiceClient) error {

	for {

		select {
		case <-ctx.Done():
			return nil
		case msg := <-c.ch:

			// publish message
			err := c.stream.Send(msg)
			if err != nil {
				return fmt.Errorf("failed to publish: %w", err)
			}
		}
	}
}

func newConn(ctx context.Context, address string) (pb.PubSubServiceClient, *grpc.ClientConn, error) {

	co, err := grpc.DialContext(ctx, address)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial: %w", err)
	}

	return pb.NewPubSubServiceClient(co), co, nil
}

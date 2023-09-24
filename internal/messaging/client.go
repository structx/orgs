// Package messaging provides a messaging client
package messaging

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/structx/common/pkg/message_broker"
)

// Client is a messaging client
type Client struct {
	addr string
}

// NewClient creates a new messaging client
func NewClient() (*Client, error) {

	address := os.Getenv("MESSAGE_BROKER_ADDRESS")
	if address == "" {
		return nil, errors.New("$MESSAGE_BROKER_ADDRESS is not set")
	}

	return &Client{
		addr: address,
	}, nil
}

// Publish publishes a message to a topic
func (c *Client) Publish(ctx context.Context, topic string, msg []byte) error {

	conn, err := c.newConn(ctx)
	if err != nil {
		return fmt.Errorf("failed to create connection: %v", err)
	}
	defer conn.Close()

	client := pb.NewMessageBrokerServiceClient(conn)

	_, err = client.Publish(ctx, &pb.Message{
		Topic:   topic,
		Payload: msg,
	})
	if err != nil {
		return fmt.Errorf("failed to publish: %v", err)
	}

	return nil
}

// Subscribe subscribes to a topic
func (c *Client) Subscribe(ctx context.Context, topic string) error {
	return nil
}

// Request sends a message to a topic and waits for a response
func (c *Client) Request(ctx context.Context, topic string, payload []byte) ([]byte, error) {

	// create connection
	conn, err := c.newConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %v", err)
	}
	defer conn.Close()

	// create client
	client := pb.NewMessageBrokerServiceClient(conn)

	// set timeout
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()

	// request response
	resp, err := client.RequestResponse(timeout, &pb.Request{
		Topic:   topic,
		Payload: payload,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to request: %v", err)
	}

	return resp.Payload, nil
}

func (c *Client) newConn(ctx context.Context) (*grpc.ClientConn, error) {

	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()

	conn, err := grpc.DialContext(timeout, c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	return conn, nil
}

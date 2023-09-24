// Package registry provides a registry client
package registry

import (
	"errors"
	"os"
)

// Client is a registry client
type Client struct {
	address string
}

// NewClient creates a new registry client
func NewClient() (*Client, error) {

	address := os.Getenv("REGISTRY_ADDRESS")
	if address == "" {
		return nil, errors.New("$REGISTRY_ADDRESS is not set")
	}

	return &Client{
		address: address,
	}, nil
}

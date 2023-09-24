// Package payment contains payment processor
package payment

import (
	"errors"
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/client"
)

// StripeClient is a payment processor
type StripeClient struct {
	client client.API
}

// NewStripeClient returns a new Stripe payment processor
func NewStripeClient() (*StripeClient, error) {

	apiKey := os.Getenv("STRIPE_API_KEY")
	if apiKey == "" {
		return nil, errors.New("STRIPE_API_KEY is not set")
	}

	stripe := client.API{}
	stripe.Init(apiKey, nil)

	return &StripeClient{
		client: stripe,
	}, nil
}

// CreateAccountHolder creates a new account holder
func (s *StripeClient) CreateAccountHolder() (string, error) {

	accounts := s.client.Accounts

	rsp, err := accounts.New(&stripe.AccountParams{
		Country: stripe.String("US"),
		Type:    stripe.String(string(stripe.AccountTypeCustom)),
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(true),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{
				Requested: stripe.Bool(true),
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create account holder: %w", err)
	}

	return rsp.ID, nil
}

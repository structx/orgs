// Package processor contains the payment processors
package processor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/adyen/adyen-go-api-library/v7/src/adyen"
	"github.com/adyen/adyen-go-api-library/v7/src/common"
	"github.com/adyen/adyen-go-api-library/v7/src/platformsaccount"
)

// Ayden is a payment processor
type Ayden struct {
	client *adyen.APIClient
}

// New returns a new Ayden payment processor
func New() (*Ayden, error) {

	aa := os.Getenv("ADYEN_API_KEY")
	if aa == "" {
		return nil, errors.New("ADYEN_API_KEY is not set")
	}

	cl := adyen.NewClient(&common.Config{
		ApiKey:      aa,
		Environment: common.TestEnv,
	})

	return &Ayden{
		client: cl,
	}, nil
}

// CreateAccountHolder creates a new account holder
func (a *Ayden) CreateAccountHolder(ctx context.Context, code, city, country, houseNumber, postal, state, street string) (string, error) {

	to, ca := context.WithTimeout(ctx, time.Second*3)
	defer ca()

	s := a.client.PlatformsAccount()
	h, _, err := s.CreateAccountHolder(&platformsaccount.CreateAccountHolderRequest{
		AccountHolderCode: code,
		AccountHolderDetails: platformsaccount.AccountHolderDetails{
			Address: &platformsaccount.ViasAddress{
				City:              city,
				Country:           country,
				HouseNumberOrName: houseNumber,
				PostalCode:        postal,
				StateOrProvince:   state,
				Street:            street,
			},
		},
	}, to)
	if err != nil {
		return "", fmt.Errorf("failed to create account holder: %w", err)
	}

	return h.AccountCode, nil
}

// Package domain contains domain models and services
package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/structx/orgs/internal/processor"
	"github.com/structx/orgs/internal/pubsub"
	"github.com/structx/orgs/internal/repository"
)

// NewOrganization ...
type NewOrganization struct {
	Name        string
	City        string
	Country     string
	HouseNumber string
	PostalCode  string
	State       string
	Street      string
}

// Organization ...
type Organization struct {
	UUID      uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OrganizationService is a service for organization
type OrganizationService struct {
	processor *processor.Ayden
	pool      *pgxpool.Pool
	pubsub    *pubsub.Client
}

// Create creates a new organization
func (os *OrganizationService) Create(ctx context.Context, newOrganization *NewOrganization) (*Organization, error) {

	// set timeout
	to, ca := context.WithTimeout(ctx, time.Second*3)
	defer ca()

	// begin transaction
	tx, err := os.pool.Begin(to)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(to)

	// create organization
	o, err := repository.New(tx).CreateOrganization(ctx, newOrganization.Name)
	if err != nil {
		return nil, err
	}

	// create payment processor account holder
	os.processor.CreateAccountHolder(ctx, o.ID.String(), newOrganization.City, newOrganization.Country,
		newOrganization.HouseNumber, newOrganization.PostalCode, newOrganization.State, newOrganization.Street)

	// commit transaction
	err = tx.Commit(to)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &Organization{
		UUID:      o.ID,
		Name:      o.Name,
		CreatedAt: o.CreatedAt.Time,
		UpdatedAt: time.Time{},
	}, nil

}

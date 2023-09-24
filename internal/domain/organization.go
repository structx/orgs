// Package domain contains domain models and services
package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/structx/orgs/internal/event"
	"github.com/structx/orgs/internal/messaging"
	"github.com/structx/orgs/internal/payment"
	"github.com/structx/orgs/internal/repository"
)

// OrganizationStatus is the status of an organization
type OrganizationStatus string

const (
	// Created newly created organization
	Created OrganizationStatus = "created"
	// Updated updated organization
	Updated OrganizationStatus = "updated"
	// NotVerified not verified organization
	NotVerified OrganizationStatus = "not_verified"
	// Verified verified organization
	Verified OrganizationStatus = "verified"
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
	CreatedBy   uuid.UUID
}

// Organization ...
type Organization struct {
	ID          uuid.UUID
	ProcessorID string
	Name        string
	Status      OrganizationStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UpdateOrganization ...
type UpdateOrganization struct {
	ID   uuid.UUID
	Name string
}

// OrganizationService is a service for organization
type OrganizationService struct {
	processor *payment.StripeClient
	pool      *pgxpool.Pool
	messaging *messaging.Client
}

// NewOrganizationService returns a new organization service
func NewOrganizationService(pool *pgxpool.Pool, processor *payment.StripeClient, messaging *messaging.Client) (*OrganizationService, error) {
	return &OrganizationService{
		pool:      pool,
		processor: processor,
		messaging: messaging,
	}, nil
}

// Create creates a new organization
func (os *OrganizationService) Create(ctx context.Context, newOrganization *NewOrganization) (*Organization, error) {

	// set timeout
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	// begin transaction
	tx, err := os.pool.Begin(timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(timeout)

	// create organization account holder through payment processor
	processorID, err := os.processor.CreateAccountHolder()
	if err != nil {
		return nil, fmt.Errorf("failed to create account holder: %w", err)
	}

	// create organization
	sqlOrg, err := repository.New(os.pool).WithTx(tx).CreateOrganization(timeout, &repository.CreateOrganizationParams{
		Name:        newOrganization.Name,
		ProcessorID: processorID,
	})
	if err != nil {
		return nil, err
	}

	// create event for message broker
	evt := &event.OrganizationCreated{
		ID:        sqlOrg.ID.String(),
		CreatedBy: newOrganization.CreatedBy.String(),
	}

	// marshal event
	msg, err := json.Marshal(evt)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	// publish message
	err = os.messaging.Publish(timeout, "organization.created", msg)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message: %w", err)
	}

	// commit transaction
	err = tx.Commit(timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &Organization{
		ID:          sqlOrg.ID,
		Name:        sqlOrg.Name,
		ProcessorID: sqlOrg.ProcessorID.(string),
		Status:      OrganizationStatus(sqlOrg.Status),
		CreatedAt:   sqlOrg.CreatedAt.Time,
		UpdatedAt:   time.Time{},
	}, nil
}

// Get gets an organization
func (os *OrganizationService) Get(ctx context.Context, id uuid.UUID) (*Organization, error) {
	return &Organization{}, nil
}

// Update updates an organization
func (os *OrganizationService) Update(ctx context.Context, organization *Organization) (*Organization, error) {
	return nil, nil
}

// Delete deletes an organization
func (os *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

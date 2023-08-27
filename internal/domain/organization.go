// Package domain contains domain models and services
package domain

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/structx/orgs/internal/payment"
	"github.com/structx/orgs/internal/pubsub"
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

// OrganizationService is a service for organization
type OrganizationService struct {
	processor payment.Processor
	db        *sql.DB
	pubsub    *pubsub.Client
}

// NewOrganizationService returns a new organization service
func NewOrganizationService(db *sql.DB, pubsub *pubsub.Client, processor payment.Processor) (*OrganizationService, error) {
	return &OrganizationService{
		db:        db,
		pubsub:    pubsub,
		processor: processor,
	}, nil
}

// Create creates a new organization
func (os *OrganizationService) Create(ctx context.Context, newOrganization *NewOrganization) (*Organization, error) {

	// set timeout
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	processorID, err := os.processor.CreateAccountHolder()
	if err != nil {
		return nil, fmt.Errorf("failed to create account holder: %w", err)
	}

	// create organization
	sqlOrg, err := repository.New(os.db).CreateOrganization(timeout, &repository.CreateOrganizationParams{
		Name:        newOrganization.Name,
		ProcessorID: processorID,
	})
	if err != nil {
		return nil, err
	}

	return &Organization{
		ID:          sqlOrg.ID,
		Name:        sqlOrg.Name,
		ProcessorID: sqlOrg.ProcessorID.(string),
		Status:      OrganizationStatus(sqlOrg.Status),
		CreatedAt:   sqlOrg.CreatedAt,
		UpdatedAt:   time.Time{},
	}, nil
}

package domain

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/structx/orgs/internal/repository"
)

// NewOrganization
type NewOrganization struct {
	Name string
}

// Organization
type Organization struct {
	UUID      uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OrganizationService is a service for organization
type OrganizationService struct {
	pool *pgxpool.Pool
}

// Create creates a new organization
func (os *OrganizationService) Create(ctx context.Context, newOrganization *NewOrganization) (*Organization, error) {

	o, err := repository.New(os.pool).CreateOrganization(ctx, newOrganization.Name)
	if err != nil {
		return nil, err
	}

	return &Organization{
		UUID:      o.ID,
		Name:      o.Name,
		CreatedAt: o.CreatedAt.Time,
		UpdatedAt: time.Time{},
	}, nil

}

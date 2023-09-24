// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: organizations.sql

package repository

import (
	"context"
)

const createOrganization = `-- name: CreateOrganization :one
INSERT INTO organizations (
    name,
    processor_id
) VALUES (
    $1, $2
) RETURNING id, processor_id, name, status, created_at, updated_at
`

type CreateOrganizationParams struct {
	Name        string
	ProcessorID interface{}
}

func (q *Queries) CreateOrganization(ctx context.Context, arg *CreateOrganizationParams) (*Organization, error) {
	row := q.db.QueryRow(ctx, createOrganization, arg.Name, arg.ProcessorID)
	var i Organization
	err := row.Scan(
		&i.ID,
		&i.ProcessorID,
		&i.Name,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

// Package event provides event definitions
package event

// OrganizationCreated ...
type OrganizationCreated struct {
	ID        string `json:"id"`
	CreatedBy string `json:"created_by"`
}

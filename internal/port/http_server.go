// Package port contains the port adapters
package port

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/structx/orgs/internal/domain"
	"go.uber.org/zap"
)

// HTTPServer is the http server
type HTTPServer struct {
	log     *zap.SugaredLogger
	service *domain.OrganizationService
}

// NewOrganizationPayload is the payload for creating a new organization
type NewOrganizationPayload struct {
	Name        string `json:"name"`
	City        string `json:"city"`
	Country     string `json:"country"`
	HouseNumber string `json:"house_number"`
	PostalCode  string `json:"postal_code"`
	State       string `json:"state"`
	Street      string `json:"street"`
}

// NewOrganizationParams is the params for creating a new organization
type NewOrganizationParams struct {
	Payload *NewOrganizationPayload `json:"payload"`
}

// Bind parse http request into NewOrganizationParams
func (nop *NewOrganizationParams) Bind(_ *http.Request) error {

	if nop.Payload == nil {
		return errors.New("payload is required")
	}

	return nil
}

func (s *HTTPServer) createOrg(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// parse request
	var p NewOrganizationParams
	if err := render.Bind(r, &p); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// create organization
	org, err := s.service.Create(ctx, &domain.NewOrganization{
		Name:        p.Payload.Name,
		City:        p.Payload.City,
		Country:     p.Payload.Country,
		HouseNumber: p.Payload.HouseNumber,
		PostalCode:  p.Payload.PostalCode,
		State:       p.Payload.State,
		Street:      p.Payload.Street,
	})
	if err != nil {
		s.log.Errorf("failed to create organization: %v", err)
		http.Error(w, "failed to create organization", http.StatusInternalServerError)
		return
	}

	// return response
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(org)
	if err != nil {
		s.log.Errorf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Package port contains the port adapters
package port

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"go.uber.org/zap"

	"github.com/structx/orgs/internal/domain"
	"github.com/structx/orgs/internal/middleware"
)

// HTTPServer is the http server
type HTTPServer struct {
	log     *zap.SugaredLogger
	service *domain.OrganizationService
}

// NewHTTPServer creates a new http server
func NewHTTPServer(log *zap.Logger, service *domain.OrganizationService) *HTTPServer {
	return &HTTPServer{
		log:     log.Sugar().Named("http_server"),
		service: service,
	}
}

// NewRouter creates a new router
func NewRouter(auth *middleware.Authenticator, srv *HTTPServer) *chi.Mux {

	router := chi.NewRouter()

	router.Mount("/api/v1", router.Group(func(r chi.Router) {
		r.Use(auth.Authenticate)
		r.Post("/organizations", srv.createOrganization)
		r.Get("/organizations/{id}", srv.fetchOrganization)
		r.Put("/organizations", srv.updateOrganiaztion)
		r.Delete("/organizations/{id}", srv.deleteOrganization)
	}))

	router.Get("/stripe_webhooks", srv.stripeWebhooks)
	router.Get("/health", srv.health)

	return router
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

func (h *HTTPServer) createOrganization(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// parse request
	var p NewOrganizationParams
	if err := render.Bind(r, &p); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// create organization
	org, err := h.service.Create(ctx, &domain.NewOrganization{
		Name:        p.Payload.Name,
		City:        p.Payload.City,
		Country:     p.Payload.Country,
		HouseNumber: p.Payload.HouseNumber,
		PostalCode:  p.Payload.PostalCode,
		State:       p.Payload.State,
		Street:      p.Payload.Street,
	})
	if err != nil {
		h.log.Errorf("failed to create organization: %v", err)
		http.Error(w, "failed to create organization", http.StatusInternalServerError)
		return
	}

	// return response
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(org)
	if err != nil {
		h.log.Errorf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *HTTPServer) fetchOrganization(w http.ResponseWriter, r *http.Request) {}

func (h *HTTPServer) updateOrganiaztion(w http.ResponseWriter, r *http.Request) {}

func (h *HTTPServer) deleteOrganization(w http.ResponseWriter, r *http.Request) {}

func (h *HTTPServer) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.log.Errorf("failed to write response: %v", err)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

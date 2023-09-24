// Package middleware provides middleware for the application
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/structx/orgs/internal/event"
	"github.com/structx/orgs/internal/messaging"
)

// ContextKey is a key for context
type ContextKey string

const (
	// UserID is the key for the user id
	UserID ContextKey = "user_id"
)

// Authenticator middleware for authenticating requests
type Authenticator struct {
	log       *zap.SugaredLogger
	messaging *messaging.Client
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator(log *zap.Logger, messaging *messaging.Client) *Authenticator {
	return &Authenticator{
		log:       log.Sugar().Named("authenticator"),
		messaging: messaging,
	}
}

// Authenticate authenticates a request
func (a *Authenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bearer := r.Header.Get("Authorization")

		token := strings.Trim(bearer, "Bearer: ")
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		evt := &event.AuthVerifyToken{
			Token: token,
		}

		payload, err := json.Marshal(evt)
		if err != nil {
			a.log.Errorf("failed to marshal event: %v", err)
			http.Error(w, "failed to create event", http.StatusInternalServerError)
			return
		}

		rsp, err := a.messaging.Request(ctx, "auth.verify_token", payload)
		if err != nil {
			a.log.Errorf("failed to request: %v", err)
			http.Error(w, "failed to request", http.StatusInternalServerError)
			return
		}

		var ack event.AuthVerifyTokenAck
		err = json.Unmarshal(rsp, &ack)
		if err != nil {
			a.log.Errorf("failed to unmarshal response: %v", err)
			http.Error(w, "failed to unmarshal response", http.StatusInternalServerError)
			return
		}

		if !ack.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, UserID, ack.ID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

package port

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/stripe/stripe-go/v75"
)

func (h *HTTPServer) stripeWebhooks(w http.ResponseWriter, r *http.Request) {

	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("failed to read request body: %v", err)
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		h.log.Errorf("failed to parse stripe event: %v", err)
		http.Error(w, "failed to parse stripe event", http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "":
	default:
		h.log.Infof("Unhandled event type: %s", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

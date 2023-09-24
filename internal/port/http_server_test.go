package port_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
	"github.com/structx/common/database"
	"github.com/structx/common/logging"
	"github.com/structx/orgs/internal/domain"
	"github.com/structx/orgs/internal/payment"
	"github.com/structx/orgs/internal/port"
)

func init() {
	os.Setenv("DB_DSN", "")
	os.Setenv("STRIPE_KEY", "")
	os.Setenv("STRIPE_WEBHOOK_SECRET", "")
}

type HTTPServerSuite struct {
	suite.Suite
	mux *chi.Mux
}

func (suite *HTTPServerSuite) SetupTest() {

	ctx := context.TODO()

	log, err := logging.NewZap()
	suite.NoError(err)

	pool, err := database.NewPGXPool(ctx)
	suite.NoError(err)

	processor, err := payment.NewStripeClient()
	suite.NoError(err)

	s, err := domain.NewOrganizationService(pool, processor, nil)
	suite.NoError(err)

	srv := port.NewHTTPServer(log, s)

	suite.mux = port.NewRouter(srv)
}

func TestHTTPServerSuite(t *testing.T) {
	suite.Run(t, new(HTTPServerSuite))
}

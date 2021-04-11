package routes

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	suite.Run(t, new(TestHealthCheckSuite))
}

type TestHealthCheckSuite struct {
	suite.Suite
	router *mux.Router
}

func (t *TestHealthCheckSuite) SetupTest() {
	route := HealthCheckRoute()
	router := mux.NewRouter()
	router.Path(route.Path()).Methods(route.Method()).Handler(route.Handler())
	t.router = router
}

func (t *TestHealthCheckSuite) TestMetaHandler() {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	t.router.ServeHTTP(rec, req)
	t.Equal(http.StatusOK, rec.Code)
}

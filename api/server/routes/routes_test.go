package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

var meta = &Metadata{
	Version:     "1.0.0",
	BuildCommit: "abcdefg",
	BuildDate:   time.Now().Truncate(time.Second).UTC(),
}

func TestMetadata(t *testing.T) {
	suite.Run(t, new(MetadataTestSuite))
}

type MetadataTestSuite struct {
	suite.Suite
	router *mux.Router
}

func (t *MetadataTestSuite) SetupTest() {
	routes := Routes(meta)
	router := mux.NewRouter()
	for _, route := range routes {
		router.Path(route.Path()).Methods(route.Method()).Handler(route.Handler())
	}
	t.router = router
}

func (t *MetadataTestSuite) TestMetaHandler() {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/meta", nil)
	t.router.ServeHTTP(rec, req)
	t.Equal(http.StatusOK, rec.Code)
	var got metaResponse
	err := json.NewDecoder(rec.Body).Decode(&got)
	t.Require().NoError(err)
	t.Equal(meta, got.Meta)
}

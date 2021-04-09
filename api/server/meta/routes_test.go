package meta

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var meta = &Metadata{
	Version:     "1.0.0",
	BuildCommit: "abcdefg",
	BuildDate:   time.Now().Truncate(time.Second).UTC(),
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

type TestSuite struct {
	suite.Suite
	router *mux.Router
}

func (t *TestSuite) SetupTest() {
	routes := Routes(meta)
	router := mux.NewRouter()
	for _, route := range routes {
		router.Path(route.Path()).Methods(route.Method()).Handler(route.Handler())
	}
	t.router = router
}

func (t *TestSuite) TestMetaHandler() {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/meta", nil)
	t.router.ServeHTTP(rec, req)
	t.Equal(http.StatusOK, rec.Code)
	var got metaResponse
	err := json.NewDecoder(rec.Body).Decode(&got)
	t.Require().NoError(err)
	t.Equal(meta, got.Meta)
}

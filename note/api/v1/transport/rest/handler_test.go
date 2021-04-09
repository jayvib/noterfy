package rest

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"noterfy/note"
	"noterfy/note/service"
	"noterfy/note/store/memory"
	"noterfy/pkg/ptrconv"
	"testing"
)

var dummyCtx = context.TODO()

var dummyNote = &note.Note{
	Title:      ptrconv.StringPointer("Unit Test"),
	Content:    ptrconv.StringPointer("This is a test"),
	IsFavorite: ptrconv.BoolPointer(true),
}

type request struct {
	Note *note.Note `json:"note"`
}

type response struct {
	Note    *note.Note `json:"note"`
	Message string     `json:"message,omitempty"`
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

type HandlerTestSuite struct {
	svc    note.Service
	store  note.Store
	routes http.Handler
	suite.Suite
	require *require.Assertions
}

func (s *HandlerTestSuite) SetupTest() {
	s.store = memory.New()
	s.svc = service.New(s.store)
	s.routes = makeHandler(s.svc)
	s.require = s.Require()
}

func (s *HandlerTestSuite) decodeResponse(rec *httptest.ResponseRecorder) response {
	var resp response
	err := json.NewDecoder(rec.Body).Decode(&resp)
	s.Require().NoError(err)
	return resp
}

func (s *HandlerTestSuite) assertMessage(resp response, want string) {
	s.Equal(want, resp.Message)
}

func (s *HandlerTestSuite) assertStatusCode(rec *httptest.ResponseRecorder, want int) {
	s.Equal(want, rec.Code)
}

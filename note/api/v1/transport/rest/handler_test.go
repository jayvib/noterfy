package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"noterfy/note"
	"noterfy/note/noteutil"
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

func (s *HandlerTestSuite) TestCreate() {

	newNote := noteutil.Copy(dummyNote)

	makeRequest := func(ctx context.Context, n *note.Note) *httptest.ResponseRecorder {
		responseRecorder := httptest.NewRecorder()
		var body bytes.Buffer
		err := json.NewEncoder(&body).Encode(&request{Note: n})
		s.require.NoError(err)
		req := httptest.NewRequest(http.MethodPost, "/note", &body)
		req = req.WithContext(ctx)
		s.routes.ServeHTTP(responseRecorder, req)
		return responseRecorder
	}

	assertNote := func(want, got *note.Note) {
		s.NotNil(got)
		s.NotEqual(uuid.Nil, got.ID)
		s.NotEmpty(got.CreatedTime)
		got.ID = uuid.Nil
		got.CreatedTime = nil
		s.Equal(want, got)
	}

	s.Run("Requesting a create note successfully", func() {
		want := noteutil.Copy(newNote)
		want.ID = uuid.Nil

		responseRecorder := makeRequest(dummyCtx, newNote)
		s.assertStatusCode(responseRecorder, http.StatusOK)
		resp := s.decodeResponse(responseRecorder)
		assertNote(want, resp.Note)
	})

	s.Run("Requesting a create note but the ID is already existing should return an error", func() {
		// TODO: When there's an ID in the note decide if the service need to
		// remove the id.
		inputNote := noteutil.Copy(newNote)
		newNote, err := s.svc.Create(dummyCtx, inputNote)
		s.require.NoError(err)

		responseRecorder := makeRequest(dummyCtx, newNote)
		s.assertStatusCode(responseRecorder, http.StatusConflict)
		resp := s.decodeResponse(responseRecorder)
		s.assertMessage(resp, "Note already exists")
	})

	s.Run("Cancelled request should return an error", func() {
		inputNote := noteutil.Copy(newNote)
		cancelledCtx, cancel := context.WithCancel(dummyCtx)
		cancel()
		responseRecorder := makeRequest(cancelledCtx, inputNote)
		s.assertStatusCode(responseRecorder, StatusClientClosed)
		resp := s.decodeResponse(responseRecorder)
		s.assertMessage(resp, "Request cancelled")
	})
}

func (s *HandlerTestSuite) TestDelete() {

	setup := func() *note.Note {
		newNote, err := s.svc.Create(dummyCtx, noteutil.Copy(dummyNote))
		s.require.NoError(err)
		return newNote
	}

	makeRequest := func(ctx context.Context, id uuid.UUID) *httptest.ResponseRecorder {
		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/note/"+id.String(), nil)
		req = req.WithContext(ctx)
		s.routes.ServeHTTP(responseRecorder, req)
		return responseRecorder
	}

	s.Run("Requesting a delete note successfully", func() {
		newNote := setup()
		responseRecorder := makeRequest(dummyCtx, newNote.ID)
		s.assertStatusCode(responseRecorder, http.StatusOK)
	})

	s.Run("Requesting a note but the ID is nil", func() {
		responseRecorder := makeRequest(dummyCtx, uuid.Nil)
		s.Equal(http.StatusBadRequest, responseRecorder.Code)
		got := s.decodeResponse(responseRecorder)
		want := "Empty note identifier"
		s.assertMessage(got, want)
	})

	s.Run("Cancelled request should return an error", func() {
		newNote := setup()
		cancelledCtx, cancel := context.WithCancel(dummyCtx)
		cancel()
		responseRecorder := makeRequest(cancelledCtx, newNote.ID)
		s.assertStatusCode(responseRecorder, StatusClientClosed)
		resp := s.decodeResponse(responseRecorder)
		s.assertMessage(resp, "Request cancelled")
	})
}

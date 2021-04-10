package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"noterfy/note"
	"noterfy/note/noteutil"
	"noterfy/note/service"
	"noterfy/note/store/memory"
	"noterfy/pkg/ptrconv"
	"noterfy/pkg/timestamp"
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

func (s *HandlerTestSuite) TestFetch() {
	// TODO: Test the ascend option.
	// TODO: Test the sort by ID.
	// TODO: Test the sort by created date.
	s.Run("Fetch successfully", func() {
		// Insert notes
		for i := 0; i < 20; i++ {
			n := new(note.Note)

			n.SetTitle(fmt.Sprintf("Title %d", i)).
				SetContent(fmt.Sprintf("Content %d", i)).
				SetIsFavorite(true)

			_, err := s.svc.Create(dummyCtx, n)
			s.require.NoError(err)
		}

		// Do a fetch request
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/notes?page=1&size=5", nil)

		s.routes.ServeHTTP(rec, req)
		s.require.Equal(http.StatusOK, rec.Code)

		// Assert
		var resp struct {
			Notes      []*note.Note `json:"notes"`
			TotalCount uint64       `json:"total_count"`
			TotalPage  uint64       `json:"total_page"`
		}

		err := json.NewDecoder(rec.Body).Decode(&resp)
		s.require.NoError(err)

		s.Len(resp.Notes, 5)
		s.Equal(uint64(20), resp.TotalCount)
		s.Equal(uint64(4), resp.TotalPage)
	})
}

func (s *HandlerTestSuite) TestGet() {

	makeRequest := func(ctx context.Context, id uuid.UUID) *httptest.ResponseRecorder {
		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/note/"+id.String(), nil)
		req = req.WithContext(ctx)
		s.routes.ServeHTTP(responseRecorder, req)
		return responseRecorder
	}

	setupNewNote := func() *note.Note {
		testNote := noteutil.Copy(dummyNote)
		newNote, err := s.svc.Create(dummyCtx, testNote)
		s.require.NoError(err)
		return newNote
	}

	s.Run("Requesting a note successfully", func() {
		testNote := setupNewNote()

		responseRecorder := makeRequest(dummyCtx, testNote.ID)
		s.Equal(http.StatusOK, responseRecorder.Code)

		want := &note.Note{
			ID:          testNote.ID,
			Title:       testNote.Title,
			Content:     testNote.Content,
			CreatedTime: timestamp.GenerateTimestamp(),
			IsFavorite:  testNote.IsFavorite,
		}

		got := s.decodeResponse(responseRecorder)

		s.Equal(want, got.Note)
	})

	s.Run("Requesting a note that not exists", func() {
		responseRecorder := makeRequest(dummyCtx, uuid.New())
		s.assertStatusCode(responseRecorder, http.StatusNotFound)
		got := s.decodeResponse(responseRecorder)
		want := "Note not found"
		s.assertMessage(got, want)
	})

	s.Run("Requesting a note but the ID is nil", func() {
		responseRecorder := makeRequest(dummyCtx, uuid.Nil)
		s.assertStatusCode(responseRecorder, http.StatusBadRequest)
		got := s.decodeResponse(responseRecorder)
		want := "Empty note identifier"
		s.assertMessage(got, want)
	})

	s.Run("Cancelled request should return an error", func() {
		inputNote := setupNewNote()
		cancelledCtx, cancel := context.WithCancel(dummyCtx)
		cancel()
		responseRecorder := makeRequest(cancelledCtx, inputNote.ID)
		s.assertStatusCode(responseRecorder, StatusClientClosed)
		resp := s.decodeResponse(responseRecorder)
		s.assertMessage(resp, "Request cancelled")
	})
}

func (s *HandlerTestSuite) TestUpdate() {

	newNote := noteutil.Copy(dummyNote)

	setup := func() *note.Note {
		newNote, err := s.svc.Create(dummyCtx, noteutil.Copy(newNote))
		s.require.NoError(err)
		s.require.NotNil(newNote)
		s.require.NotEqual(uuid.Nil, newNote.ID)
		return newNote
	}

	makeRequest := func(ctx context.Context, n *note.Note) *httptest.ResponseRecorder {
		responseRecorder := httptest.NewRecorder()
		var body bytes.Buffer
		err := json.NewEncoder(&body).Encode(&request{Note: n})
		s.require.NoError(err)
		req := httptest.NewRequest(http.MethodPut, "/note", &body)
		req = req.WithContext(ctx)
		s.routes.ServeHTTP(responseRecorder, req)
		return responseRecorder
	}

	assertNote := func(want, got *note.Note) {
		s.Equal(want, got)
	}

	s.Run("Request for update successfully", func() {

		// Update the note via request
		updatedNote := noteutil.Copy(setup())
		updatedNote.Title = ptrconv.StringPointer("Updated Title")

		want := noteutil.Copy(updatedNote)
		want.UpdatedTime = timestamp.GenerateTimestamp()

		responseRecorder := makeRequest(dummyCtx, updatedNote)
		s.assertStatusCode(responseRecorder, http.StatusOK)
		resp := s.decodeResponse(responseRecorder)
		assertNote(want, resp.Note)
	})

	s.Run("Request for update note that is not exist should return an error", func() {
		updatedNote := noteutil.Copy(dummyNote)
		updatedNote.ID = uuid.New()
		responseRecorder := makeRequest(dummyCtx, updatedNote)
		s.assertStatusCode(responseRecorder, http.StatusNotFound)
		resp := s.decodeResponse(responseRecorder)
		s.assertMessage(resp, "Note not found")
	})

	s.Run("Cancelled request should return an error", func() {
		logrus.SetLevel(logrus.DebugLevel)
		updatedNote := noteutil.Copy(setup())
		updatedNote.Title = ptrconv.StringPointer("Updated Title")

		cancelledCtx, cancel := context.WithCancel(dummyCtx)
		cancel()
		responseRecorder := makeRequest(cancelledCtx, updatedNote)
		s.assertStatusCode(responseRecorder, StatusClientClosed)
		resp := s.decodeResponse(responseRecorder)
		s.assertMessage(resp, "Request cancelled")
	})
}

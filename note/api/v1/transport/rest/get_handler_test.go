package rest

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"noterfy/note"
	"noterfy/note/noteutil"
	"noterfy/pkg/timestamp"
)

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

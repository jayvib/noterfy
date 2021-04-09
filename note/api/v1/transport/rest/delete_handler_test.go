package rest

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"noteapp/note"
	"noteapp/note/noteutil"
)

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

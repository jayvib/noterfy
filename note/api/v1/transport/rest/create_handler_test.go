package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"noterfy/note"
	"noterfy/note/noteutil"
)

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

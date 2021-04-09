package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"noterfy/note"
	"noterfy/note/noteutil"
	"noterfy/pkg/ptrconv"
	"noterfy/pkg/timestamp"
)

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

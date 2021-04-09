package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"noterfy/note"
)

func (s *HandlerTestSuite) TestFetch() {
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

package rest

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"noteapp/note"
)

type updateService interface {
	Update(ctx context.Context, n *note.Note) (*note.Note, error)
}

type updateRequest struct {
	Note *note.Note `json:"note"`
}

type updateResponse struct {
	Note *note.Note `json:"note"`
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (reqOut interface{}, err error) {
	var req updateRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := r.Body.Close()
		if cerr != nil && err == nil {
			err = cerr
		}
	}()

	return req, nil
}

func makeUpdateEndpoint(svc updateService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(updateRequest)

		updatedNote, err := svc.Update(ctx, request.Note)
		if err != nil {
			return errorWrapper{
				origErr:    err,
				message:    getMessage(err),
				statusCode: getStatusCode(err),
			}, nil
		}
		return updateResponse{Note: updatedNote}, nil
	}
}

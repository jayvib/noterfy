package rest

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"noteapp/note"
)

// createService is here to follow the interface segregation principle.
type createService interface {
	Create(ctx context.Context, n *note.Note) (*note.Note, error)
}

type createRequest struct {
	Note *note.Note `json:"note"`
}

type createResponse struct {
	Note *note.Note `json:"note"`
}

func decodeCreateRequest(_ context.Context, r *http.Request) (response interface{}, err error) {
	var req createRequest
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

func makeCreateEndpoint(svc createService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(createRequest)
		newNote, err := svc.Create(ctx, request.Note)
		if err != nil {
			return errorWrapper{
				origErr:    err,
				message:    getMessage(err),
				statusCode: getStatusCode(err),
			}, nil
		}
		return createResponse{Note: newNote}, nil
	}
}

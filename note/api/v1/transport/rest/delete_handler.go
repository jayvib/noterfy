package rest

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type deleteService interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

type deleteRequest struct {
	ID uuid.UUID `json:"id"`
}

type deleteResponse struct {
	Message string `json:"message"`
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return deleteRequest{ID: uuid.MustParse(id)}, nil
}

func makeDeleteEndpoint(svc deleteService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(deleteRequest)
		err := svc.Delete(ctx, request.ID)
		if err != nil {
			return errorWrapper{
				origErr:    err,
				message:    getMessage(err),
				statusCode: getStatusCode(err),
			}, nil
		}
		return deleteResponse{"Successfully Deleted"}, nil
	}
}

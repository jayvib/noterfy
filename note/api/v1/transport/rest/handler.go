package rest

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"noterfy/note"
)

// @title Noterfy Note Service
// @version 0.2.1
// @description  Noterfy Note Service.
// @termsOfService http://swagger.io/terms/
//
// @contact.name Jayson Vibandor
// @contact.email jayson.vibandor@gmail.com
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /v1

// makeHandler initializes all the routes for the note service
// handlers and return the routed handler.
func makeHandler(svc note.Service) http.Handler {
	router := mux.NewRouter()
	getHandler := httptransport.NewServer(
		makeGetEndpoint(svc),
		decodeGetRequest,
		encodeResponse,
	)

	createHandler := httptransport.NewServer(
		makeCreateEndpoint(svc),
		decodeCreateRequest,
		encodeResponse,
	)

	updateHandler := httptransport.NewServer(
		makeUpdateEndpoint(svc),
		decodeUpdateRequest,
		encodeResponse,
	)

	deleteHandler := httptransport.NewServer(
		makeDeleteEndpoint(svc),
		decodeDeleteRequest,
		encodeResponse,
	)

	fetchHandler := httptransport.NewServer(
		makeFetchEndpoint(svc),
		decodeFetchRequest,
		encodeResponse,
	)

	router.Handle("/note/{id}", getHandler).Methods(http.MethodGet)
	router.Handle("/note", createHandler).Methods(http.MethodPost)
	router.Handle("/note", updateHandler).Methods(http.MethodPut)
	router.Handle("/note/{id}", deleteHandler).Methods(http.MethodDelete)
	router.Handle("/notes", fetchHandler).Methods(http.MethodGet)

	return router
}

// CreateRequest is a container for the create request.
type CreateRequest struct {
	Note *note.Note `json:"note"`
}

// CreateResponse is a container fo a successful create response.
type CreateResponse struct {
	Note *note.Note `json:"note"`
}

func decodeCreateRequest(_ context.Context, r *http.Request) (response interface{}, err error) {
	var req CreateRequest
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

// CreateRequest godoc
// @Summary Create a new note.
// @Description Creating a new note. The client can assign the note ID with a UUID value but the service will return a conflict error when the note with the ID provided is already exists.
// @Accept json
// @Produce json
// @Param CreateRequest body CreateRequest true "A body containing the new note"
// @Success 200 {object} CreateResponse "Successfully created a new note"
// @Failure 409 {object} ResponseError "Conflict error due to the new note with an ID already exists in the service"
// @Failure 499 {object} ResponseError "Cancel error when the request was aborted"
// @Router /note [post]
func makeCreateEndpoint(svc createService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(CreateRequest)
		newNote, err := svc.Create(ctx, request.Note)
		if err != nil {
			return errorWrapper{
				origErr:    err,
				message:    getMessage(err),
				statusCode: getStatusCode(err),
			}, nil
		}
		return CreateResponse{Note: newNote}, nil
	}
}

type DeleteRequest struct {
	ID uuid.UUID `json:"id"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return DeleteRequest{ID: uuid.MustParse(id)}, nil
}

// DeleteRequest godoc
// @Summary Delete an existing note.
// @Description Delete an existing note.
// @Param id path string true "ID of the note"
// @Success 200 {string} string "Successful deleting a note"
// @Failure 400 {object} ResponseError "Note's ID parameter is not provided in the path"
// @Failure 499 {object} ResponseError "Cancel error when the request was aborted"
// @Router /note/{id} [delete]
func makeDeleteEndpoint(svc deleteService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(DeleteRequest)
		err := svc.Delete(ctx, request.ID)
		if err != nil {
			return errorWrapper{
				origErr:    err,
				message:    getMessage(err),
				statusCode: getStatusCode(err),
			}, nil
		}
		return DeleteResponse{"Successfully Deleted"}, nil
	}
}

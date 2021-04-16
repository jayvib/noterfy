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
	"strconv"
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
//
// @tag.name Note API
// @tag.description Use to interact to the Noterfy note service.
//
// @schemes http https
//
// @query.collection.format multi
//
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

// DeleteRequest is a container for the delete request.
type DeleteRequest struct {
	ID uuid.UUID `json:"id"`
}

// DeleteResponse is a container for the delete response.
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
// @Failure 500 {object} ResponseError "Unexpected server internal error"
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

// FetchRequest is a container for the fetch request API.
type FetchRequest struct {
	Pagination *note.Pagination
}

// FetchResponse is a container for the fetch response API.
type FetchResponse struct {
	Notes      []*note.Note `json:"notes"`
	TotalCount uint64       `json:"total_count" example:"2"`
	TotalPage  uint64       `json:"total_page" example:"5"`
}

func decodeFetchRequest(_ context.Context, r *http.Request) (response interface{}, err error) {

	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")
	sortBy := r.URL.Query().Get("sort_by")
	ascendRaw := r.URL.Query().Get("ascending")
	if ascendRaw == "" {
		// Default will be ascend=true
		ascendRaw = "true"
	}

	ascend, err := strconv.ParseBool(ascendRaw)
	if err != nil {
		ascend = true
	}

	response = FetchRequest{
		Pagination: &note.Pagination{
			Size:      convertAtoU(size),
			Page:      convertAtoU(page),
			SortBy:    note.GetSortBy(sortBy),
			Ascending: ascend,
		},
	}

	return
}

// FetchRequest godoc
// @Summary Fetches notes from the service.
// @Description Fetches notes from the service.
// @Accept json
// @Produce json
// @Param page query int false "The page number of the fetch pagination. Default is page=1."
// @Param size query int false "The page size of the fetch pagination. Default is size=25."
// @Param sort_by query string false "An option for sorting the notes in the response. Default is sort_by=title. [title/id/created_date]"
// @Param ascending query bool false "An option for sorting the results in ascending or descending. Default is ascending=true"
// @Success 200 {object} FetchResponse "Successfully fetches notes"
// @Failure 499 {object} ResponseError "Cancel error when the request was aborted"
// @Failure 500 {object} ResponseError "Unexpected server internal error"
// @Router /notes [get]
func makeFetchEndpoint(svc fetchService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		request := req.(FetchRequest)
		iter, err := svc.Fetch(ctx, request.Pagination)
		if err != nil {
			return newErrorWrapper(err), nil
		}

		defer func() {
			cerr := iter.Close()
			if cerr != nil && err == nil {
				err = cerr
			}
		}()

		var notes []*note.Note
		for iter.Next() {
			n := iter.Note()
			notes = append(notes, n)
		}

		if iter.Error() != nil {
			return newErrorWrapper(err), nil
		}

		resp = FetchResponse{
			Notes:      notes,
			TotalCount: iter.TotalCount(),
			TotalPage:  iter.TotalPage(),
		}

		return
	}
}

// GetRequest is a container for the get request API.
type GetRequest struct {
	ID uuid.UUID `json:"id"`
}

// GetResponse is a container for the get response API.
type GetResponse struct {
	Note *note.Note `json:"note"`
}

// GetRequest godoc
// @Summary Get the note from the service.
// @Description Get the note from the service if exists. When the note is not exists it will return a NotFound response status.
// @Param id path string true "ID of the note"
// @Success 200 {object} GetResponse "Successful getting the note"
// @Failure 404 {object} ResponseError "Note is not found in the service"
// @Failure 400 {object} ResponseError "Note's ID parameter is not provided in the path"
// @Failure 499 {object} ResponseError "Cancel error when the request was aborted"
// @Failure 500 {object} ResponseError "Unexpected server internal error"
// @Router /note/{id} [get]
func makeGetEndpoint(svc getService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(GetRequest)
		v, err := svc.Get(ctx, request.ID)
		if err != nil {
			return errorWrapper{
				origErr:    err,
				statusCode: getStatusCode(err),
				message:    getMessage(err),
			}, nil
		}
		return GetResponse{Note: v}, nil
	}
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return GetRequest{ID: uuid.MustParse(id)}, nil
}

type updateService interface {
	Update(ctx context.Context, n *note.Note) (*note.Note, error)
}

// UpdateRequest is a container for the update request API.
type UpdateRequest struct {
	Note *note.Note `json:"note"`
}

// UpdateResponse is a container for the update response of the API.
type UpdateResponse struct {
	Note *note.Note `json:"note"`
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (reqOut interface{}, err error) {
	var req UpdateRequest
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

// UpdateRequest godoc
// @Summary Update an existing note.
// @Description Updating an existing note. If the note to be updated is not found the API will respond a NotFound status.
// @Accept json
// @Produce json
// @Param UpdateRequest body UpdateRequest true "A body containing the updated note"
// @Success 200 {object} UpdateResponse "Successfully updated the note"
// @Failure 404 {object} ResponseError "Note to be update is not found in the service"
// @Failure 499 {object} ResponseError "Cancel error when the request was aborted"
// @Router /note [put]
func makeUpdateEndpoint(svc updateService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(UpdateRequest)

		updatedNote, err := svc.Update(ctx, request.Note)
		if err != nil {
			return errorWrapper{
				origErr:    err,
				message:    getMessage(err),
				statusCode: getStatusCode(err),
			}, nil
		}
		return UpdateResponse{Note: updatedNote}, nil
	}
}

package rest

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"noterfy/note"
)

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

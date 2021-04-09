package rest

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"noteapp/api"
	"noteapp/note"
	nhttp "noteapp/pkg/http"
)

// Routes returns all the routes that is part of the
// note API service.
func Routes(svc note.Service) []api.Route {
	return getRoutes(svc)
}

func getRoutes(svc note.Service) []api.Route {

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

	routes := []api.Route{
		&nhttp.Route{HandlerValue: getHandler, MethodValue: http.MethodGet, PathValue: "/v1/note/{id}"},
		&nhttp.Route{HandlerValue: createHandler, MethodValue: http.MethodPost, PathValue: "/v1/note"},
		&nhttp.Route{HandlerValue: updateHandler, MethodValue: http.MethodPut, PathValue: "/v1/note"},
		&nhttp.Route{HandlerValue: deleteHandler, MethodValue: http.MethodDelete, PathValue: "/v1/note/{id}"},
		&nhttp.Route{HandlerValue: fetchHandler, MethodValue: http.MethodGet, PathValue: "/v1/notes"},
	}
	return routes
}

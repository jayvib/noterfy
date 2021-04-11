package routes

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"noterfy/api"
	nhttp "noterfy/pkg/http"
	"time"
)

// Routes takes a metadata an return the routes for
// metadata.
func Routes(meta *Metadata) []api.Route {
	return []api.Route{
		HealthCheckRoute(),
		MetadataRoute(meta),
	}
}

// HealthCheckRoute return the route for the health check endpoint.
func HealthCheckRoute() api.Route {

	handler := httptransport.NewServer(
		makeHealthCheckEndpoint(),
		decodeHealthCheckRequest,
		encodeResponse,
	)

	return &nhttp.Route{
		HandlerValue: handler,
		MethodValue:  http.MethodGet,
		PathValue:    "/health",
	}
}

// HealthCheckRequest is a container for the health check request.
type HealthCheckRequest struct{}

// HealthCheckResponse is a container for the health check response.
type HealthCheckResponse struct {
	Message string `json:"message,omitempty"`
}

func decodeHealthCheckRequest(context.Context, *http.Request) (response interface{}, err error) {
	return HealthCheckRequest{}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(resp)
}

func makeHealthCheckEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// Naive-way implementation.
		response = &HealthCheckResponse{
			Message: "OK",
		}
		return
	}
}

// MetadataRoute takes a metadata meta and return its wrapped route.
func MetadataRoute(meta *Metadata) api.Route {
	metaHandler := httptransport.NewServer(
		makeMetaEndpoint(meta),
		decodeMetaRequest,
		encodeMetaResponse,
	)

	return &nhttp.Route{
		HandlerValue: metaHandler,
		MethodValue:  http.MethodGet,
		PathValue:    "/meta",
	}
}

// Metadata contains the information for the API server.
type Metadata struct {
	Version     string    `json:"version,omitempty"`
	BuildCommit string    `json:"build_commit,omitempty"`
	BuildDate   time.Time `json:"build_date,omitempty"`
}

type metaRequest struct{}

type metaResponse struct {
	Meta *Metadata `json:"meta"`
}

func decodeMetaRequest(context.Context, *http.Request) (response interface{}, err error) {
	return metaRequest{}, nil
}

func encodeMetaResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

func makeMetaEndpoint(meta *Metadata) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		resp := &metaResponse{
			Meta: meta,
		}
		return resp, nil
	}
}

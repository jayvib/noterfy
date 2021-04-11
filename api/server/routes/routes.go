package routes

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"noterfy/api"
	nhttp "noterfy/pkg/http"
)

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

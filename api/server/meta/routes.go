package meta

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
	return getRoutes(meta)
}

func getRoutes(meta *Metadata) (routes []api.Route) {

	metaHandler := httptransport.NewServer(
		makeMetaEndpoint(meta),
		decodeMetaRequest,
		encodeResponse,
	)

	routes = append(routes, &nhttp.Route{
		HandlerValue: metaHandler,
		MethodValue:  http.MethodGet,
		PathValue:    "/meta",
	})

	return
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

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
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

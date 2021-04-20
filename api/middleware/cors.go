package middleware

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"noterfy/api"
)

// NewCORSMiddleware initializes the CORS middleware. It takes
// an optional conf for extra configuration. If nil is provided
// it will use the default configuration.
func NewCORSMiddleware(conf *CORSConfig) api.NamedMiddleware {
	return api.NewNamedMiddleware("CORS", CORS(conf))
}

// CORSConfig contains the configuration for the CORS middleware.
type CORSConfig struct {
	AllowedOrigins           []string
	AllowedOriginFunc        func(string) bool
	AllowedOriginRequestFunc func(r *http.Request, origin string) bool
	AllowCredentials         bool
	AllowedMethods           []string
	AllowedHeaders           []string
	ExposedHeaders           []string
	MaxAge                   int
	Passthrough              bool
	Debug                    bool
}

// CORS is a middleware for the CORS handler.
func CORS(conf *CORSConfig) mux.MiddlewareFunc {
	var options cors.Options

	if conf != nil {
		options = cors.Options{
			AllowedOrigins:         conf.AllowedOrigins,
			AllowOriginFunc:        conf.AllowedOriginFunc,
			AllowOriginRequestFunc: conf.AllowedOriginRequestFunc,
			AllowedMethods:         conf.AllowedMethods,
			AllowedHeaders:         conf.AllowedHeaders,
			ExposedHeaders:         conf.ExposedHeaders,
			MaxAge:                 conf.MaxAge,
			AllowCredentials:       conf.AllowCredentials,
			OptionsPassthrough:     conf.Passthrough,
			Debug:                  conf.Debug,
		}
	}

	return func(next http.Handler) http.Handler {
		return cors.New(options).Handler(next)
	}
}

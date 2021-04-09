package api

import "github.com/gorilla/mux"

// NewNamedMiddleware takes the name of the middleware and the middleware
// function and return a named middleware struct.
func NewNamedMiddleware(name string, mw mux.MiddlewareFunc) NamedMiddleware {
	return NamedMiddleware{
		Middleware: mw,
		Name:       name,
	}
}

// NamedMiddleware is a container struct
// for the middleware and its name.
type NamedMiddleware struct {
	Middleware mux.MiddlewareFunc
	Name       string
}

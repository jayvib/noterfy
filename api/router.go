package api

import "net/http"

// Route defines an individual API route in the docker server.
type Route interface {
	// Handler returns the raw function to create the http handler.
	Handler() http.Handler
	// Method returns the http method that the route responds to.
	Method() string
	// Path returns the subpath where the route respond to.
	Path() string
}

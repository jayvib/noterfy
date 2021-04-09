package http

import "net/http"

// Route implements server.Route. It is a container
// for holding the necessary returns of its methods.
type Route struct {
	HandlerValue http.Handler
	MethodValue  string
	PathValue    string
}

// Handler returns the raw function to create the http handler.
func (r Route) Handler() http.Handler {
	return r.HandlerValue
}

// Method returns the http method that the route responds to.
func (r Route) Method() string {
	return r.MethodValue
}

// Path returns the subpath where the route respond to.
func (r Route) Path() string {
	return r.PathValue
}

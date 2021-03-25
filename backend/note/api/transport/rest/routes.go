package rest

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"noteapp/api/server"
	"noteapp/note"
)

// Route implements server.Route. It is a container
// for holding the necessary returns of its methods.
type Route struct {
	handler http.Handler
	method  string
	path    string
}

// Handler returns the raw function to create the http handler.
func (r Route) Handler() http.Handler {
	return r.handler
}

// Method returns the http method that the route responds to.
func (r Route) Method() string {
	return r.method
}

// Path returns the subpath where the route respond to.
func (r Route) Path() string {
	return r.path
}

// Routes returns all the routes that is part of the
// note API service.
func Routes(svc note.Service) []server.Route {
	return getRoutes(svc)
}

func getRoutes(svc note.Service) []server.Route {

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

	routes := []server.Route{
		&Route{getHandler, http.MethodGet, "/note/{id}"},
		&Route{createHandler, http.MethodPost, "/note"},
		&Route{updateHandler, http.MethodPut, "/note"},
		&Route{deleteHandler, http.MethodDelete, "/note/{id}"},
	}
	return routes
}

package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

// New takes config for all the arguments that the server needs and
// return a server instance.
func New(conf *Config) *Server {
	server := &Server{
		Port:        conf.Port,
		Middlewares: conf.Middlewares,
		Routes:      conf.Routes,
	}
	server.init()
	return server
}

// Config contains all the arguments that the server needs.
type Config struct {
	// Port is the port of the server
	Port int
	// Middlewares are the middlewares to be apply to the
	// handler.
	Middlewares []mux.MiddlewareFunc
	// Routes are the API routes that will be register to the server.
	Routes []Route
}

// Server is the wrapper for all the bootstrapping of a typical server.
type Server struct {
	Port        int
	Middlewares []mux.MiddlewareFunc
	server      *http.Server
	Routes      []Route
}

func (s *Server) init() {
	router := mux.NewRouter()

	for _, routes := range s.Routes {
		logrus.Infof("Registering: %s\t%s\n", routes.Method(), routes.Path())
		router.Path(routes.Path()).Methods(routes.Method()).Handler(routes.Handler())
	}

	router.Use(s.Middlewares...)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: router,
	}
}

// ListenAndServe serves clients request by the server.
func (s *Server) ListenAndServe() error {
	logrus.Info("Running in port:", s.Port)
	return s.server.ListenAndServe()
}

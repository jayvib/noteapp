package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"noteapp/note/api/middleware"
)

func New(conf *Config) *Server {
	server := &Server{
		Port:        conf.Port,
		Middlewares: conf.Middlewares,
		Handler:     conf.Handler,
	}
	server.init()
	return server
}

type Config struct {
	Port        int
	Middlewares []middleware.Middleware
	Handler     http.Handler
}

type Server struct {
	Port        int
	Middlewares []middleware.Middleware
	Handler     http.Handler
	server      *http.Server
}

func (s *Server) init() {
	// Apply Middlewares
	handler := s.Handler
	if len(s.Middlewares) > 0 {
		handler = middleware.Apply(handler, s.Middlewares...)
	}
	mux := http.NewServeMux()
	mux.Handle("/", handler)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: mux,
	}
}

func (s *Server) ListenAndServe() error {
	logrus.Info("Running in port:", s.Port)
	return s.server.ListenAndServe()
}

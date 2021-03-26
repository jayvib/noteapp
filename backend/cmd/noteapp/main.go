package main

import (
	"github.com/gorilla/mux"
	"log"
	"noteapp/api/server"
	"noteapp/note/api/middleware"
	"noteapp/note/api/transport/rest"
	noteservice "noteapp/note/service"
	"noteapp/note/store/memory"
)

const port = 50001

func main() {
	svc := noteservice.New(memory.New())
	srv := server.New(&server.Config{
		Port: port,
		Middlewares: []mux.MiddlewareFunc{
			middleware.Logging,
		},
		HTTPRoutes: rest.Routes(svc),
	})

	defer srv.Close()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

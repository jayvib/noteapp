package main

import (
	"log"
	"noteapp/api"
	"noteapp/api/middleware"
	"noteapp/api/server"
	"noteapp/note/api/v1/transport/rest"
	noteservice "noteapp/note/service"
	"noteapp/note/store/memory"
)

const port = 50001

func main() {
	svc := noteservice.New(memory.New())
	srv := server.New(&server.Config{
		Port: port,
		Middlewares: []api.NamedMiddleware{
			middleware.NewLoggingMiddleware(),
		},
		HTTPRoutes: rest.Routes(svc),
	})

	defer srv.Close()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

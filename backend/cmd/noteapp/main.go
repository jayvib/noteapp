package main

import (
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
	handler := rest.MakeHandler(svc)
	srv := server.New(&server.Config{
		Port: port,
		Middlewares: []middleware.Middleware{
			middleware.Logging,
		},
		Handler: handler,
	})
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

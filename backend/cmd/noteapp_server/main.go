package main

import (
	"log"
	"noteapp/api"
	"noteapp/api/middleware"
	"noteapp/api/server"
	"noteapp/api/server/meta"
	"noteapp/note/api/v1/transport/rest"
	noteservice "noteapp/note/service"
	filestore "noteapp/note/store/file"
	"os"
)

// ❎ TODO: Add versioning
// Follow this link as reference:
// https://blog.carlmjohnson.net/post/2021/how-to-use-go-embed/

const (
	port       = 50001
	dbFileName = "note.pb"
)

func main() {

	file, err := os.OpenFile(dbFileName, os.O_CREATE|os.O_RDWR, 0666)
	mustNoError(err)
	defer func() { _ = file.Close() }()

	svc := noteservice.New(filestore.New(file))
	srv := server.New(&server.Config{
		Port: port,
		Middlewares: []api.NamedMiddleware{
			middleware.NewLoggingMiddleware(),
		},
	})

	srv.AddRoutes(meta.Routes(&meta.Metadata{
		Version:     Version,
		BuildCommit: BuildCommit,
		BuildDate:   BuildDate,
	})...)
	srv.AddRoutes(rest.Routes(svc)...)

	defer srv.Close()
	mustNoError(srv.ListenAndServe())
}

func mustNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

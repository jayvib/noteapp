package main

import (
	"log"
	"net/http"
	"noteapp/note/api/transport/rest"
	noteservice "noteapp/note/service"
	"noteapp/note/store/memory"
)

func main() {
	svc := noteservice.New(memory.New())
	http.Handle("/", rest.MakeHandler(svc))
	if err := http.ListenAndServe(":50001", nil); err != nil {
		log.Fatal(err)
	}
}

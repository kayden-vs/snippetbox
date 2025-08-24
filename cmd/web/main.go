package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	r.Get("/", home)
	r.Get("/snippet/view", snippetView)
	r.Post("/snippet/create", snippetCreate)

	err := http.ListenAndServe(":4000", r)
	log.Fatal(err)
}

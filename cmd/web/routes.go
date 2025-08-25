package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	r.Get("/", app.home)
	r.Get("/snippet/view", app.snippetView)
	r.Post("/snippet/create", app.snippetCreate)

	return r
}

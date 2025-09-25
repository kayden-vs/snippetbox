package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(secureHeaders)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Serve static files
	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// for session
	dynamic := app.sessionManager.LoadAndSave

	r.With(dynamic).Get("/", app.home)

	r.Route("/snippet", func(r chi.Router) {
		// Apply middleware to everything inside this group
		r.Use(dynamic)

		r.Get("/view/{id}", app.snippetView)
		r.Get("/create", app.snippetCreateForm)
		r.Post("/create", app.snippetCreate)
	})

	return r
}

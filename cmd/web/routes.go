package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	// --- Global Middleware Stack ---
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(app.secureHeaders)

	// --- Public Routes ---
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// --- Dynamic & Session-Enabled Routes ---
	r.Group(func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave)
		r.Use(noSurf)
		r.Use(app.authenticate)

		r.Get("/", app.home)
		r.Get("/snippet/view/{id}", app.snippetView)
		r.Get("/user/signup", app.userSignup)
		r.Post("/user/signup", app.userSignupPost)
		r.Get("/user/login", app.userLogin)
		r.Post("/user/login", app.userLoginPost)

		// --- Authenticated-Only Routes ---
		r.Group(func(r chi.Router) {
			r.Use(app.requireAuthentication)

			r.Get("/snippet/create", app.snippetCreate)
			r.Post("/snippet/create", app.snippetCreatePost)
			r.Post("/user/logout", app.userLogoutPost)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	return r
}

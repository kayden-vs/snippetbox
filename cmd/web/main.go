package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()

	r.Get("/", home)
	r.Get("/snippet/view", snippetView)
	r.Get("/snippet/create", snippetCreate)

	err := http.ListenAndServe(":4000", r)
	log.Fatal(err)
}

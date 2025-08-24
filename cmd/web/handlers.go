package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kayden-vs/snippetbox/ui/html/pages"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	component := pages.HomePage()
	component.Render(context.Background(), w)
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Create a new snippet..."))
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kayden-vs/snippetbox/ui/html/pages"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	component := pages.HomePage()
	err := component.Render(context.Background(), w)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Create a new snippet..."))
}

package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kayden-vs/snippetbox/internal/models"
	"github.com/kayden-vs/snippetbox/ui/html/pages"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.Render(w, pages.HomePage(snippets))
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	createdStr := snippet.Created.Format("02 Jan 2006 at 15:04")
	expiresStr := snippet.Expires.Format("02 Jan 2006 at 15:04")
	component := pages.ViewSnippet(snippet.ID, snippet.Title, snippet.Content, createdStr, expiresStr)
	app.Render(w, component)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	//will be taking this from user later
	title := "0 Snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.infoLog.Println("New Data added, Id: ", id)

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

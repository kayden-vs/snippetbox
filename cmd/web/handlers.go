package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kayden-vs/snippetbox/internal/models"
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
	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// component := pages.HomePage()
	// err := component.Render(context.Background(), w)
	// if err != nil {
	// 	app.errorLog.Println(err.Error())
	// 	app.serverError(w, err)
	// 	return
	// }
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
	fmt.Fprintf(w, "%+v", snippet)
	// fmt.Fprintf(w, "ID: %d\nTitle: %s\nContent: %s\nCreated: %v\nExpires: %v\n", snippet.ID, snippet.Title, snippet.Content, snippet.Created, snippet.Expires)
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

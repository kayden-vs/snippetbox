package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kayden-vs/snippetbox/internal/models"
	"github.com/kayden-vs/snippetbox/internal/validator"
	"github.com/kayden-vs/snippetbox/ui/html/pages"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.Render(w, pages.HomePage(snippets))
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
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

	flash := app.sessionManager.PopString(r.Context(), "flash")

	createdStr := snippet.Created.Format("02 Jan 2006 at 15:04")
	expiresStr := snippet.Expires.Format("02 Jan 2006 at 15:04")
	component := pages.ViewSnippet(snippet.ID, snippet.Title, snippet.Content, createdStr, expiresStr, flash)
	app.Render(w, component)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Create form instance with validator
	var form pages.SnippetCreateForm
	form.Validator = validator.Validator{FieldErrors: make(map[string]string)}

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Perform validation
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		app.Render(w, pages.SnippetForm(form))
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.infoLog.Println("New Data added, Id: ", id)

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	form := pages.SnippetCreateForm{
		Expires:   365,
		Validator: validator.Validator{FieldErrors: make(map[string]string)},
	}
	app.Render(w, pages.SnippetForm(form))
}

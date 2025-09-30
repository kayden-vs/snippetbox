package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/a-h/templ"
	"github.com/go-playground/form"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace) //to get the correct line number of err and avoid err reference to this file

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// RenderPage injects flash and isAuthenticated into the page component.
func (app *application) RenderPage(
	w http.ResponseWriter,
	r *http.Request,
	renderFunc func(flash string, isAuthenticated bool) templ.Component,
) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	flash := app.sessionManager.PopString(r.Context(), "flash")
	isAuth := app.isAuthenticated(r)
	err := renderFunc(flash, isAuth).Render(r.Context(), w)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}

	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	w.WriteHeader(http.StatusInternalServerError)
	app.render(w, r, "error.page.tmpl", &templateData{
		ErrorMessage: http.StatusText(http.StatusInternalServerError),
	})
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
	// http.Error(w, http.StatusText(status), status)
	w.WriteHeader(status)
	app.render(w, r, "error.page.tmpl", &templateData{
		ErrorMessage: http.StatusText(status),
	})
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	app.clientError(w, r, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()
	return td
}

func (app *application) authenticatedUser(r *http.Request) int {
	user_id := 0
	cookies := r.Cookies()
	for _, c := range cookies {
		if c.Name == "session_cookie" {
			user_id = app.sessions.GetUser(c.Value)
			break
		}
	}
	return user_id
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, r, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	buf.WriteTo(w)
}

package main

import (
	"net/http"

	"github.com/justinas/alice"
	"snippetbox.peter/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux() // creating a http-router

	//fileserverfs created http handler which serves the embedded
	// files in ui.Files. And our static files are contained
	// in the "static" folder, and we don't need to strip
	// prefixes now -> any requests that start with /static/
	// can be just passed directly to the file server
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /ping", ping)

	// unprotected routes use dynamic middleware chain
	dynamic := alice.New(app.sessionManager.LoadAndSave, app.noSurf, app.authenticate)

	mux.Handle("GET /about", dynamic.ThenFunc(app.about))

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// new protected middleware chain that requires user to be authenticated
	// includes requireAuthentication middleware
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	standart := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standart.Then(mux)
}

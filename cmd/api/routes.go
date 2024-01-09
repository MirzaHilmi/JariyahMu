package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(app.notFound)
	mux.MethodNotAllowed(app.methodNotAllowed)

	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)
	mux.Use(app.authenticate)

	v1 := chi.NewRouter()

	v1.Get("/status", app.status)

	v1.Route("/auth", func(r chi.Router) {
		r.Post("/signup", app.signUserUp)
		r.Post("/login", app.logUserIn)
	})

	v1.Group(func(v1 chi.Router) {
		v1.Use(app.requireAuthenticatedUser)

		v1.Get("/protected", app.protected)
	})

	v1.Group(func(v1 chi.Router) {
		v1.Use(app.requireBasicAuthentication)

		v1.Get("/basic-auth-protected", app.protected)
	})

	mux.Mount("/api/v1", v1)

	return mux
}

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) Routes() http.Handler {
	cr := chi.NewRouter()
	//register middleware
	cr.Use(middleware.Recoverer)
	cr.Use(app.enableCORS)
	cr.Handle("/", http.StripPrefix("/", http.FileServer((http.Dir("./html/")))))

	cr.Route("/web", func(r chi.Router) {
		r.Post("/auth", app.Authenticate)
		// /refresh token
		r.Get("/refresh-token", app.RefreshUsingCookie)
		// /logout
		r.Get("/logout", app.deleteRefreshCookie)
	})
	//authentication routes - auth handler, refresh
	cr.Post("/auth", app.Authenticate)
	cr.Post("/refresh-token", app.Refresh)

	cr.Route("/users", func(ur chi.Router) {
		//use auth middleware
		ur.Use(app.authRequired)
		ur.Get("/", app.AllUsers)
		ur.Get("/{UserID}", app.GetUser)
		ur.Delete("/{UserID}", app.DeleteUser)
		ur.Put("/", app.InsertUser)
		ur.Patch("/", app.UpdateUser)
	})
	return cr
}

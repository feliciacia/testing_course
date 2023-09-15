package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) Routes() http.Handler {
	mux := chi.NewRouter()
	//register middleware
	mux.Use(middleware.Recoverer)
	//mux.Use(app.enableCORS)
	//authentication routes - auth handler, refresh
	mux.Post("/auth", app.Authenticate)
	mux.Post("/refresh-token", app.Refresh)
	//test handler
	mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		var payload = struct {
			Message string `json:"messsage"`
		}{
			Message: "hello, world",
		}
		_ = app.writeJSON(w, http.StatusOK, payload)
	})
	//protected route
	mux.Route("/users", func(mux chi.Router) {
		//use auth middleware
		mux.Get("/", app.AllUsers)
		mux.Get("/{UserID}", app.GetUser)
		mux.Delete("/{UserID}", app.DeleteUser)
		mux.Put("/{UserID}", app.InsertUser)
		mux.Patch("/{UserID}", app.UpdateUser)
	})
	return mux
}

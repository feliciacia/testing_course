package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) routes() http.Handler {
	mux := chi.NewRouter()
	//register middleware
	mux.Use(middleware.Recoverer)
	//mux.Use(app.enableCORS)
	//authentication routes - auth handler, refresh
	//test handler
	//protected route
	return mux
}

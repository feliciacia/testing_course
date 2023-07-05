package routes

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	Session *scs.SessionManager
}

func (app *Application) Routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer) //register middleware
	mux.Use(app.addIPToContext)
	mux.Use(app.Session.LoadAndSave)
	mux.Get("/", app.Home) //register routes
	mux.Post("/login", app.Login)
	//static assets
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

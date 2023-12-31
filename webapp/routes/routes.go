package routes

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/felicia/testing_course/webapp/pkg/db/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	Session *scs.SessionManager
	DB      repository.DatabaseRepo
	DSN     string
}

func (app *Application) Routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer) //register middleware
	mux.Use(app.addIPToContext)
	mux.Use(app.Session.LoadAndSave)
	mux.Get("/", app.Home) //register routes
	mux.Post("/login", app.Login)

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(app.auth)
		mux.Get("/profile", app.Profile)
		mux.Post("/upload-profile-pic", app.UploadProfilePic)
	})
	//static assets
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

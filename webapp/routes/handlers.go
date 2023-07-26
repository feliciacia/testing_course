package routes

import (
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/felicia/testing_course/webapp/pkg/data"
)

var PathtoTemplate = "../templates/"

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	td := make(map[string]interface{})
	if app.Session.Exists(r.Context(), "test") {
		msg := app.Session.GetString(r.Context(), "test")
		td["test"] = msg
	} else {
		app.Session.Put(r.Context(), "test", "Hit this page at"+time.Now().UTC().String())
	}
	templateData := &TemplateData{
		IP:   app.ipFromContext(r.Context()),
		Data: td,
	}
	_ = app.render(w, r, "home.page.gohtml", templateData)
}

func (app *Application) Profile(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "profile.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP    string
	Data  map[string]interface{}
	Error string
	Flash string
	User  data.User
}

func (app *Application) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) error {
	parsedTemplate, err := template.ParseFiles(path.Join(PathtoTemplate, t), path.Join(PathtoTemplate, "base.layout.gohtml"))

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	td.IP = app.ipFromContext(r.Context())
	//information to pass through templates
	//pull something from session with key error if error exists
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Flash = app.Session.PopString(r.Context(), "flash")

	//execute template and pass the data
	err = parsedTemplate.Execute(w, td)

	if err != nil {
		return err
	}
	return nil
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	form := NewForm(r.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		//redirect to login page with error message
		app.Session.Put(r.Context(), "error", "Invalid login credentials") //including error message
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email") //"email" declared at html

	password := r.Form.Get("password")

	user, err := app.DB.GetUserByEmail(email)

	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid login") //error message for no email
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//authenticate
	//if not authenticate then redirect with error
	if !app.authenticate(r, user, password) {
		app.Session.Put(r.Context(), "error", "Invalid login")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//prevent fixation attack //regenerate session
	_ = app.Session.RenewToken(r.Context())

	app.Session.Put(r.Context(), "flash", "successfully log in")
	//redirect to other page
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (app *Application) authenticate(r *http.Request, user *data.User, pswd string) bool {
	if valid, err := user.PasswordMatches(pswd); err != nil || !valid {
		return false
	}
	app.Session.Put(r.Context(), "user", user)
	return true
}

package routes

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
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
	data := &TemplateData{
		IP:   app.ipFromContext(r.Context()),
		Data: td,
	}
	_ = app.render(w, r, "home.page.gohtml", data)
}

func (app *Application) Profile(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "profile.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *Application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	parsedTemplate, err := template.ParseFiles(path.Join(PathtoTemplate, t), path.Join(PathtoTemplate, "base.layout.gohtml"))

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	data.IP = app.ipFromContext(r.Context())
	//execute template and pass the data
	err = parsedTemplate.Execute(w, data)

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
	log.Println(password, user.FirstName)
	//authenticate
	//if not authenticate then redirect with error
	//prevent fixation attack //regenerate session
	_ = app.Session.RenewToken(r.Context())

	app.Session.Put(r.Context(), "flash", "successfully log in")
	//redirect to other page
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

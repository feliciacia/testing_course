package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
)

var pathtoTemplate = "../templates/"

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "home.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *Application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	parsedTemplate, err := template.ParseFiles(path.Join(pathtoTemplate, t), path.Join(pathtoTemplate, "base.layout.gohtml"))

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
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	form := NewForm(r.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		fmt.Fprint(w, "failed validation")
		return
	}

	email := r.Form.Get("email") //"email" declared at html

	password := r.Form.Get("password")

	log.Println(email, password)
	fmt.Fprint(w, email)
}

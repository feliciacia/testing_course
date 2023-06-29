package routes

import (
	"html/template"
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
	parsedTemplate, err := template.ParseFiles(path.Join(pathtoTemplate, t))

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

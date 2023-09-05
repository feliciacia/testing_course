package routes

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/felicia/testing_course/webapp/pkg/data"
)

var PathtoTemplate = "../templates/"

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	td := make(map[string]interface{})
	if app.Session.Exists(r.Context(), "test") {
		msg := app.Session.GetString(r.Context(), "test")
		td["test"] = msg
		log.Printf("Retrieved session data: %s", msg)
	} else {
		log.Println("Session data not found")
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

	if app.Session.Exists(r.Context(), "user") {
		td.User = app.Session.Get(r.Context(), "user").(data.User)
	}
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
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	log.Printf("Received login request with email: %s, password: %s", email, password)

	if !form.Valid() {
		//redirect to login page with error message
		app.Session.Put(r.Context(), "error", "Invalid login credentials") //including error message
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user, err := app.DB.GetUserByEmail(email)

	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid login") //error message for no email
		log.Printf("User %s not found in the database", email)
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
	flashmsg := app.Session.GetString(r.Context(), "flash")
	log.Printf("Flash message after setting: %s", flashmsg)
	//redirect to other page
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (app *Application) authenticate(r *http.Request, user *data.User, pswd string) bool {
	if valid, err := user.PasswordMatches(pswd); err != nil || !valid {
		log.Printf("Authentication failed for user: %s. Password matches: %v, Error: %v", user.Email, valid, err)
		return false
	}
	app.Session.Put(r.Context(), "user", user)
	return true
}

func (app *Application) UploadProfilePic(w http.ResponseWriter, r *http.Request) {
	//call a function that extracts a file from an upload

	//get the user from session

	//create a var from type data.UserImage

	//insert the user image into user_image

	//refresh the sessional variable "user"

	//redirect back to profile page
}

type UploadedFile struct {
	OriginalFileName string
	FileSize         int64
}

func (app *Application) UploadFiles(r *http.Request, uploadDir string) ([]*UploadedFile, error) {
	var uploadedfiles []*UploadedFile

	err := r.ParseMultipartForm(int64(1024 * 1024 * 5)) //5 means 5 Gb //no exceed of this 1024 * 1024 * 5

	if err != nil {
		return nil, fmt.Errorf("the uploaded file is too big and must be less than %d files", 1024*1024*5)
	}
	for _, fHeaders := range r.MultipartForm.File {
		for _, header := range fHeaders {
			uploadedfiles, err := func(uploadedfiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile
				infile, err := header.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()
				uploadedFile.OriginalFileName = header.Filename

				var outfile *os.File
				defer outfile.Close()

				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.OriginalFileName)); nil != err {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, infile)
					if err != nil {
						return nil, err
					}
					uploadedFile.FileSize = fileSize
				}
				uploadedfiles = append(uploadedfiles, &uploadedFile)
				return uploadedfiles, nil
			}(uploadedfiles)
			if err != nil {
				return uploadedfiles, nil
			}
		}
	}
	return uploadedfiles, nil
}

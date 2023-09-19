package main

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Credential struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *Application) Authenticate(w http.ResponseWriter, r *http.Request) {
	var creds Credential
	//read a json payload
	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}
	//look up the user by email address
	user, err := app.DB.GetUserByEmail(creds.Username)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}
	//check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}
	//generate token
	tokenPairs, err := app.GenerateTokenPair(user)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}
	//send token
	_ = app.writeJSON(w, http.StatusOK, tokenPairs)
}

func (app *Application) Refresh(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) AllUsers(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) GetUser(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) InsertUser(w http.ResponseWriter, r *http.Request) {

}

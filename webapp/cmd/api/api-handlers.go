package main

import "net/http"

func (app *Application) Authenticate(w http.ResponseWriter, r *http.Request) {
	//read a json payload

	//look up the user by email address

	//check password
	//generate token
	//send token
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
package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/felicia/testing_course/webapp/pkg/data"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
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
	http.SetCookie(w, &http.Cookie{
		Name:     "__Host-refresh_token",
		Path:     "/",
		Value:    tokenPairs.RefreshToken,
		Expires:  time.Now().Add(refreshTokenExpiry),
		MaxAge:   int(refreshTokenExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   true,
	})
	//send token
	_ = app.writeJSON(w, http.StatusOK, tokenPairs)
}

func (app *Application) Refresh(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	refreshToken := r.Form.Get("refresh_token")
	claims := &Claims{}
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.JWTSecret), nil
	})
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	if time.Unix(claims.ExpiresAt.Unix(), 0).Sub(time.Now()) > 30*time.Second { //compared the difference time of expired token time and the time now
		app.errorJSON(w, errors.New("refreshed token no need to be renewed yet"), http.StatusTooEarly)
		return
	}
	//get the user id from claims
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	user, err := app.DB.GetUser(userID)
	if err != nil {
		app.errorJSON(w, errors.New("unknown users"), http.StatusBadRequest)
		return
	}
	tokenPairs, err := app.GenerateTokenPair(user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "__Host-refresh_token",
		Path:     "/",
		Value:    tokenPairs.RefreshToken,
		Expires:  time.Now().Add(refreshTokenExpiry),
		MaxAge:   int(refreshTokenExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   true,
	})
	_ = app.writeJSON(w, http.StatusOK, tokenPairs)
}

func (app *Application) AllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.DB.AllUsers()
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		fmt.Printf("Error in AllUsers handler: %v\n", err)
		return
	}
	_ = app.writeJSON(w, http.StatusOK, users)
}

func (app *Application) GetUser(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "UserID")
	fmt.Printf("UserID from URL: %s\n", userIDParam)

	userID, err := strconv.Atoi(userIDParam)
	fmt.Printf("UserID: %d", userID)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		fmt.Printf("Error converting UserID to int: %v\n", err)
		return
	}
	user, err := app.DB.GetUser(userID)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		fmt.Printf("Error fetching user by ID: %v\n", err)
		return
	}
	_ = app.writeJSON(w, http.StatusOK, user)
}

func (app *Application) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := app.readJSON(w, r, &user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	err = app.DB.UpdateUser(user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent) //done update and no response body
}

func (app *Application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "UserID")
	fmt.Printf("UserID from URL: %s\n", userIDParam)

	userID, err := strconv.Atoi(userIDParam)
	fmt.Printf("UserID: %d\n", userID)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		fmt.Println("cant convert")
		return
	}
	err = app.DB.DeleteUser(userID)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		fmt.Println("cant delete")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) InsertUser(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := app.readJSON(w, r, &user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		fmt.Printf("Error reading JSON: %v\n", err)
		return
	}
	_, err = app.DB.InsertUser(user)
	fmt.Printf("Request Body: %+v\n", user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

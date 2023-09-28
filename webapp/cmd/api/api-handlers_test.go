package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/felicia/testing_course/webapp/pkg/data"
	"github.com/go-chi/chi"
)

func Test_app_authenticate(t *testing.T) {
	var theTests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{"valid user", `{"email":"admin@example.com", "password":"secret"}`, http.StatusOK},
		{"not json", `I'm not JSON`, http.StatusUnauthorized},
		{"empty json", `{}`, http.StatusUnauthorized},
		{"empty email", `{"email":""}`, http.StatusUnauthorized},
		{"empty password", `{"email":"admin@example.com"}`, http.StatusUnauthorized},
		{"invalid user", `{"email":"admin@someotherdomain.com", "password":"secret"}`, http.StatusUnauthorized},
	}
	for _, e := range theTests {
		var reader io.Reader
		reader = strings.NewReader(e.requestBody)
		req, _ := http.NewRequest("POST", "/auth", reader)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.Authenticate)

		handler.ServeHTTP(rr, req)
		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}

func Test_app_refreshToken(t *testing.T) {
	var tests = []struct {
		name               string
		token              string
		expectedStatusCode int
		resetRefreshTime   bool
	}{
		{"valid", "", http.StatusOK, true},
		{"valid but not yet ready to expire", "", http.StatusTooEarly, false},
		{"expired token", expiredToken, http.StatusBadRequest, false},
	}
	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}
	oldRefreshTime := refreshTokenExpiry
	for _, e := range tests {
		var refreshtkn string
		if e.token == "" {
			if e.resetRefreshTime {
				refreshTokenExpiry = time.Second * 1
			}
			tokens, _ := app.GenerateTokenPair(&testUser)
			refreshtkn = tokens.RefreshToken
		} else {
			refreshtkn = e.token
		}
		postedData := url.Values{
			"refresh_token": {refreshtkn},
		}
		req, _ := http.NewRequest("POST", "/refresh-token", strings.NewReader(postedData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.Refresh)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: expected status of %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
		refreshTokenExpiry = oldRefreshTime
	}
}

func Test_app_userHandler(t *testing.T) {
	var tests = []struct {
		name           string
		method         string
		json           string
		paramID        string
		handler        http.HandlerFunc
		expectedStatus int
	}{
		{"AllUsers", "GET", "", "", app.AllUsers, http.StatusOK},
		{"deleteUsers", "DELETE", "", "1", app.DeleteUser, http.StatusNoContent},
		{"deleteUsers bad URL param", "DELETE", "", "Y", app.DeleteUser, http.StatusBadRequest},
		{"getUsers valid", "GET", "", "1", app.GetUser, http.StatusOK},
		{"getUsers bad URL param", "GET", "", "Y", app.GetUser, http.StatusBadRequest},
		{
			"updateUser valid",
			"PATCH",
			`{"id":1, "first_name":"Administrator", "last_name":"User", "email":"admin@example.com"}`,
			"",
			app.UpdateUser,
			http.StatusNoContent,
		},
		{
			"updateUser invalid",
			"PATCH",
			`{"id:100, "first_name":"Administrator", "last_name":"User", "email":"admin@example.com"}`,
			"",
			app.UpdateUser,
			http.StatusBadRequest,
		},
		{
			"updateUser invalid json",
			"PATCH",
			`{"id":1, first_name:"Administrator", "last_name":"User","email":"admin@example.com"}`,
			"",
			app.UpdateUser,
			http.StatusBadRequest,
		},
		{
			"insertUser valid",
			"PUT",
			`{"first_name":"Jack", "last_name":"Smith", "email":"jack@example.com"}`,
			"",
			app.InsertUser,
			http.StatusNoContent,
		},
		{
			"insertUser invalid",
			"PUT",
			`{"foo":"bar","first_name":"Jack", "last_name":"Smith", "email":"jack@example.com"}`,
			"",
			app.InsertUser,
			http.StatusBadRequest,
		},
		{
			"insertUser invalid json",
			"PUT",
			`{"first_name:"Jack", "last_name":"Smith", "email":"jack@example.com"}`,
			"",
			app.InsertUser,
			http.StatusBadRequest,
		},
	}
	for _, e := range tests {
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("UserID", e.paramID)

		var req *http.Request

		if e.json == "" {
			req, _ = http.NewRequest(e.method, fmt.Sprintf("/users/%s", e.paramID), nil)
		} else {
			req, _ = http.NewRequest(e.method, fmt.Sprintf("/users/%s", e.paramID), strings.NewReader(e.json))
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(e.handler)
		handler.ServeHTTP(rr, req)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		log.Printf("Test case: %s, Method: %s, URL: %s", e.name, e.method, req.URL.String())
		if e.json != "" {
			log.Printf("Request Body:\n%s", e.json)
		}

		if e.paramID != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("UserID", e.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}
		log.Printf("ParamID set to: %s", e.paramID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		rr = httptest.NewRecorder()
		handler = http.HandlerFunc(e.handler)
		handler.ServeHTTP(rr, req)
		responseString := rr.Body.String()
		if responseString != "" {
			log.Printf("Response Body:\n%s", responseString)
		}
		log.Printf("Response Status Code: %d", rr.Code)
		if rr.Code != e.expectedStatus {
			t.Errorf("%s: wrong status returned; expected %d, but got %d", e.name, e.expectedStatus, rr.Code)
		}
	}
}

func Test_app_refresh_cookie(t *testing.T) {
	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}
	tokens, _ := app.GenerateTokenPair(&testUser)
	testCookie := &http.Cookie{
		Name:     "__Host-refresh_token",
		Path:     "/",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(refreshTokenExpiry),
		MaxAge:   int(refreshTokenExpiry),
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   true,
	}

	badCookie := &http.Cookie{
		Name:     "__Host-refresh_token",
		Path:     "/",
		Value:    "somebadstring",
		Expires:  time.Now().Add(refreshTokenExpiry),
		MaxAge:   int(refreshTokenExpiry),
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   true,
	}
	var tests = []struct {
		name           string
		addCookie      bool
		cookie         *http.Cookie
		expectedStatus int
	}{
		{"valid cookie", true, testCookie, http.StatusOK},
		{"invalid cookie", true, badCookie, http.StatusBadRequest},
	}
	for _, e := range tests {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		if e.addCookie {
			req.AddCookie(e.cookie)
		}
		handler := http.HandlerFunc(app.RefreshUsingCookie)
		handler.ServeHTTP(rr, req)
		if rr.Code != e.expectedStatus {
			t.Errorf("%s: wrong status code returned; expected %d, but got %d", e.name, e.expectedStatus, rr.Code)
		}
	}
}

func Test_deleteRefreshCookie(t *testing.T) {
	
}
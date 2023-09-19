package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_EnableCORS(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
	var tests = []struct {
		name         string
		method       string
		expectHeader bool
	}{
		{"preflight", "OPTIONS", true},
		{"get", "GET", false},
	}

	for _, e := range tests {
		handlertoTest := app.enableCORS(nextHandler)
		req := httptest.NewRequest(e.method, "http://testing", nil)
		rr := httptest.NewRecorder()
		handlertoTest.ServeHTTP(rr, req)
		if e.expectHeader && rr.Header().Get("Access-Control-Allow-Credentials") == "" {
			t.Errorf("%s: expected header, but did not find it", e.name)
		}
		if e.expectHeader && rr.Header().Get("Access-Control-Allow-Credentials") != "" {
			t.Errorf("%s: expected no header, but got one", e.name)
		}
	}
}

package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_app_routes(t *testing.T) {
	var registered = []struct {
		route  string
		method string
	}{
		{"/auth", "POST"},
		{"/refresh-token", "POST"},
		{"/users/", "GET"},
		{"/users/{UserID}", "GET"},
		{"/users/{UserID}", "DELETE"},
		{"/users/", "PUT"},
		{"/users/", "PATCH"},
	}
	mux := app.Routes()
	chiRoutes := mux.(chi.Routes)
	for _, routes := range registered {
		if !routeExists(routes.route, routes.method, chiRoutes) {
			t.Errorf("route %s is not registered", routes.route)
		}
	}
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false
	_ = chi.Walk(chiRoutes, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})
	return found
}

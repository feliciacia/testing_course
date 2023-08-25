package routes

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

type ContextKey string

const ContextUserKey ContextKey = "user_ip"

func (app *Application) ipFromContext(ctx context.Context) string {
	return ctx.Value(ContextUserKey).(string)
}

func (app *Application) addIPToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forwardedFor := r.Header.Get("X-Forwarded-For")
		log.Printf("X-Forwarded-For header: %s", forwardedFor)
		var ctx = context.Background()
		ip, err := getIP(r)
		if err != nil {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
			if len(ip) == 0 {
				ip = "unknown"
			}
			ctx = context.WithValue(r.Context(), ContextUserKey, ip)
		} else {
			ctx = context.WithValue(r.Context(), ContextUserKey, ip)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown", err
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
	}
	forward := r.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		ip = strings.TrimSpace(strings.Split(forward, ",")[0])
	}
	if len(forward) == 0 {
		ip = "forward"
	}
	log.Printf("Final IP: %s", ip)
	return ip, nil
}

func (app *Application) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userData := app.Session.GetString(r.Context(), "user")
		log.Println("Auth middleware triggered")
		if !app.Session.Exists(r.Context(), "user") {
			app.Session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		log.Printf("User session data: %s", userData)
		next.ServeHTTP(w, r)
	})
}

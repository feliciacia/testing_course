package main

import (
	"log"
	"net/http"
	"primeapp/testing_course/webapp/routes"
)

func main() {
	app := &routes.Application{}
	mux := app.Routes()
	app.Session = getSession()
	log.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"

	"github.com/felicia/testing_course/webapp/routes"
)

func main() {
	app := &routes.Application{}
	app.Session = routes.GetSession()
	log.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", app.Routes())
	if err != nil {
		log.Fatal(err)
	}
}
 
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/felicia/testing_course/webapp/pkg/db/repository"
	"github.com/felicia/testing_course/webapp/pkg/db/repository/dbrepo"
)

const port = 8090

type Application struct {
	DSN       string
	DB        repository.DatabaseRepo
	Domain    string
	JWTSecret string
}

func main() {
	var app Application
	flag.StringVar(&app.Domain, "domain", "example.com", "Domain for application, e.g. company.com")
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Posgtres connection")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "signing secret")
	flag.Parse()

	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	log.Printf("Starting api on port %d\n", port)

	err = http.ListenAndServe(":8090", app.Routes())
	if err != nil {
		log.Fatal(err)
	}
}

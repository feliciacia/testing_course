package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/felicia/testing_course/webapp/pkg/db/repository"
	"github.com/felicia/testing_course/webapp/pkg/db/repository/dbrepo"
)

const port = 8090

type Application struct {
	domain    string
	db        repository.DatabaseRepo
	DSN       string
	JWTsecret string
}

func main() {
	var app Application
	flag.StringVar(&app.domain, "domain", "example.com", "domain for application, e.g company.com")
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.StringVar(&app.JWTsecret, "jwt-secret", "774ef6cc9ec7f62fe885c04849d9469e6aaf3f8e3568581ebe6d08ad25f15f0e", "signing secret")
	flag.Parse()
	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	app.db = &dbrepo.PostgresDBRepo{DB: conn}
	log.Printf("Starting api on port : %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}

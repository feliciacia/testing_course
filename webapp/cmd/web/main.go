package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/felicia/testing_course/webapp/pkg/data"
	"github.com/felicia/testing_course/webapp/pkg/db/repository/dbrepo"
	"github.com/felicia/testing_course/webapp/routes"
)

func main() {
	gob.Register(data.User{})
	app := &routes.Application{}
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	fmt.Println("DSN:", &app.DSN)
	flag.Parse() //read value where it has to be
	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	app.Session = routes.GetSession()
	log.Println("Starting server on port 8000...")
	err = http.ListenAndServe(":8000", app.Routes())
	if err != nil {
		log.Fatal(err)
	}
}

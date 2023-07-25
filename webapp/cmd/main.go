package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/felicia/testing_course/webapp/pkg/db"
	"github.com/felicia/testing_course/webapp/routes"
)

func main() {
	app := &routes.Application{}
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.Parse() //read value where it has to be
	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	app.DB = db.PostgresConn{DB: conn}
	app.Session = routes.GetSession()
	log.Println("Starting server on port 8080...")
	err = http.ListenAndServe(":8080", app.Routes())
	if err != nil {
		log.Fatal(err)
	}
}

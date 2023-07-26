package main

import (
	"log"
	"os"
	"testing"

	"github.com/felicia/testing_course/webapp/pkg/db"
	"github.com/felicia/testing_course/webapp/routes"
)

var app routes.Application

func TestMain(m *testing.M) {
	routes.PathtoTemplate = "../templates/"
	app.Session = routes.GetSession()
	app.DSN = "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5"
	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	app.DB = db.PostgresConn{DB: conn}
	os.Exit(m.Run())
}

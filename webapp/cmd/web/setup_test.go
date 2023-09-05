package main

import (
	"os"
	"testing"

	"github.com/felicia/testing_course/webapp/pkg/db/repository/dbrepo"
	"github.com/felicia/testing_course/webapp/routes"
)

var app routes.Application

func TestMain(m *testing.M) {
	routes.PathtoTemplate = "../templates/"
	app.Session = routes.GetSession()
	app.DB = &dbrepo.PostgresDBRepo{}
	os.Exit(m.Run())
}

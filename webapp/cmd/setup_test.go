package main

import (
	"os"
	"testing"

	"github.com/felicia/testing_course/webapp/routes"
)

var app routes.Application

func Test_main(m *testing.M) {
	app.Session = routes.GetSession()
	os.Exit(m.Run())
}

package main

import (
	"os"
	"testing"

	"github.com/felicia/testing_course/webapp/pkg/db/repository/dbrepo"
)

var app Application

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBRepo{}
	app.Domain = "example.com"
	app.JWTSecret = "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160"
	os.Exit(m.Run())
}

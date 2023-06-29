package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handlers(t *testing.T) {
	var theTest = []struct {
		name       string
		url        string
		statuscode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/notfound", http.StatusNotFound},
	}
	var app Application
	routes := app.Routes()
	//create test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTest {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.statuscode {
			t.Errorf("for %s: expected status %d, but got %d", e.name, e.statuscode, resp.StatusCode)
		}
	}
}

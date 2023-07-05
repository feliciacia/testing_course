package routes

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handlers(t *testing.T) {

	var app Application
	app.Session = GetSession()
	routes := app.Routes()
	//create test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()
	var theTest = []struct {
		name       string
		url        string
		statuscode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/notfound", http.StatusNotFound},
	}
	pathtoTemplate = "./../templates/"
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

func Test_Home_Old(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	var app Application
	app.Session = GetSession()
	req = AddContextAndSessionToRequest(req, app)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.Home)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("TestHome returned wrong status code; expected 200 but got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if !strings.Contains(string(body), `<small>From Session =`) {
		t.Error("did not find correct text in html")
	}
}

func Test_Home(t *testing.T) {
	var tests = []struct {
		name         string
		putInSession string
		expectedHTML string
	}{
		{"first visit", "", "<small>From Session ="},
		{"second visit", "hello, world", "<small>From Session = hello, world"},
	}
	for _, e := range tests {
		req, _ := http.NewRequest("GET", "/", nil)
		var app Application
		app.Session = GetSession()
		req = AddContextAndSessionToRequest(req, app)
		_ = app.Session.Destroy(req.Context())
		if e.putInSession != "" {
			app.Session.Put(req.Context(), "test", e.putInSession)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.Home)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("TestHome returned wrong status code; expected 200 but got %d", rr.Code)
		}
		body, _ := io.ReadAll(rr.Body)
		if !strings.Contains(string(body), e.expectedHTML) {
			t.Errorf("%s: did not find %s in response body", e.name, e.expectedHTML)
		}
	}
}

func Get_Context(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), ContextUserKey, "unknown")
	return ctx
}

func AddContextAndSessionToRequest(req *http.Request, app Application) *http.Request {
	req = req.WithContext(Get_Context(req))
	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))
	return req.WithContext(ctx)
}

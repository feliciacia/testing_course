package routes

import (
	"context"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/felicia/testing_course/webapp/pkg/db/repository/dbrepo"
)

func Test_handlers(t *testing.T) {

	var app Application
	app.Session = GetSession()

	conn, _ := app.ConnectToDB()
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	defer conn.Close()

	// Set the application's DB field to the connected test database
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	routes := app.Routes()
	//create test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	} //create http transport with tls configuration disabling the verification for testing purposes
	client := &http.Client{ //create client http to request to server
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse //not redirecting
		},
	}
	var theTest = []struct {
		name                    string
		url                     string
		statuscode              int
		expectedURL             string
		expectedfirststatuscode int //status code after redirect
	}{
		{"home", "/", http.StatusOK, "/", http.StatusOK},
		{"404", "/notfound", http.StatusNotFound, "/notfound", http.StatusNotFound},
		{"profile", "/user/profile", http.StatusOK, "/", http.StatusTemporaryRedirect},
	}
	for _, e := range theTest {
		resp, err := ts.Client().Get(ts.URL + e.url)
		t.Logf("Testing %s: %s", e.name, ts.URL+e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.statuscode {
			t.Errorf("for %s: expected status %d, but got %d", e.name, e.statuscode, resp.StatusCode)
		}
		if resp.Request.URL.Path != e.expectedURL {
			t.Errorf("for %s: expected final url of %s but got %s", e.name, e.expectedURL, resp.Request.URL.Path)
		}
		t.Logf("Response status code: %d", resp.StatusCode)
		t.Logf("Response location header: %s", resp.Header.Get("Location"))
		resp2, _ := client.Get(ts.URL + e.url) //concatenate test server url with the expected one in order to get the construct complete url to send http request
		if resp2.StatusCode != e.expectedfirststatuscode {
			t.Errorf("%s: expected first returned status code to be %d but got %d", e.name, e.expectedfirststatuscode, resp2.StatusCode)
		}
	}
}

func Test_Home(t *testing.T) {
	var tests = []struct {
		name         string
		putInSession string
		expectedHTML string
	}{
		{"first visit", "", "<small>From Session:"},
		{"second visit", "hello, world", "<small>From Session: hello, world"},
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

func Test_render_with_bad_template(t *testing.T) {
	PathtoTemplate = "./testdata/"

	var app Application
	app.Session = GetSession()

	req, _ := http.NewRequest("GET", "/", nil)
	req = AddContextAndSessionToRequest(req, app)
	rr := httptest.NewRecorder()

	err := app.render(rr, req, "bad.page.gohtml", &TemplateData{})
	if err == nil { //expect error but said no error
		t.Error("expected error from bad template, but did not get one")
	}
	if rr.Code != http.StatusBadRequest {
		t.Error("expect bad request but not getting the bad request")
	}
	PathtoTemplate = "../templates/"
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

func Test_login(t *testing.T) {
	var app Application
	app.Session = GetSession()
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.Parse() //read value where it has to be
	conn, err := app.ConnectToDB()
	if err != nil {
		t.Fatalf("Error connecting to database: %s", err)
	}
	defer conn.Close()

	// Set the application's DB field to the connected test database
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	var tests = []struct {
		name       string
		postedData url.Values //handling form post to web browser
		//url.Values for handling various value
		expectedStatusCode int
		expectedPage       string
	}{
		{
			name: "valid login",
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedPage:       "/user/profile",
		},
		{
			name: "missing form data",
			postedData: url.Values{
				"email":    {""},
				"password": {""},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedPage:       "/",
		},
		{
			name: "bad credentials",
			postedData: url.Values{
				"email":    {"you@there.com"},
				"password": {"password"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedPage:       "/",
		},
		{
			name: "user not found",
			postedData: url.Values{
				"email":    {"admin2@example.com"},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedPage:       "/",
		},
		{
			name: "bad credentials",
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"password"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedPage:       "/",
		},
	}
	appSessionManager := app.Session
	//for valid login
	for _, e := range tests {
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(e.postedData.Encode())) //request and encoded to the format

		req = AddContextAndSessionToRequest(req, app)                       //add context and the session
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //encodead the format as the html form data
		rr := httptest.NewRecorder()                                        //capture the response test
		handler := http.HandlerFunc(app.Login)                              //handler for test http request
		handler.ServeHTTP(rr, req)                                          //send the req to handler and capture the response in rr

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: returned wrong status code: expected %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		log.Printf("Received login request with email: %s, password: %s", e.postedData.Get("email"), e.postedData.Get("password"))

		actualPage := rr.Header().Get("Location")
		t.Log("location:", actualPage)

		if actualPage != e.expectedPage {
			t.Errorf("%s: expected location %s, but got %s", e.name, e.expectedPage, actualPage)
		}
		sessionData := appSessionManager.GetString(req.Context(), "test")
		t.Logf("Retrieved session data: %v", sessionData)
	}

}

package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_addIPToContext(t *testing.T) {
	tests := []struct {
		headerName  string
		headerValue string
		addr        string
		emptyAddr   bool
	}{
		{"", "", "", false},
		{"", "", "", true},
		{"X-Forwarded-For", "192.0.1.2", "", false},
		{"", "", "hello:world", false},
	}

	var app Application

	//create a dummy handler that use to check the context
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//make sure the value exist
		val := r.Context().Value(ContextUserKey)
		if val == nil {
			t.Error(ContextUserKey, "not present")
		}
		// make sure get the string back
		ip, ok := val.(string)
		if !ok {
			t.Error("not string")
		}
		t.Log(ip)
	})
	for _, e := range tests {
		//create handler to test
		handlertotest := app.addIPToContext(nextHandler)
		req := httptest.NewRequest("GET", "http://testing", nil)
		if e.emptyAddr {
			req.RemoteAddr = ""
		}

		if len(e.headerName) > 0 { //it means the header name exists
			req.Header.Add(e.headerName, e.headerValue)
		}

		if len(e.addr) > 0 { // it means the address exists
			req.RemoteAddr = e.addr
		}

		//call the handler so dummy handler can perform the test
		handlertotest.ServeHTTP(httptest.NewRecorder(), req)
	}
}

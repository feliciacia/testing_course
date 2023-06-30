package routes

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_has(t *testing.T) {
	form := NewForm(nil) //1 st scenario
	has := form.Has("bdhbcbwjkc")
	if has {
		t.Error("form shows have field when it shouldn't")
	}
	postedData := url.Values{} //2 nd scenario
	postedData.Add("name", "value")
	form = NewForm(postedData)
	has = form.Has("name")
	if !has {
		t.Error("form shows that does not have field when it should")
	}
}

func Test_required(t *testing.T) {
	r := httptest.NewRequest("POST", "/request", nil) //1 st scenario
	form := NewForm(r.PostForm)
	form.Required("a", "b", "c") //requiring these 3 forms
	if form.Valid() {            //the form we required are not being there but it said being there, then it should return error
		t.Error("form shows valid when required are missing")
	}
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")
	r, _ = http.NewRequest("POST", "/request", nil)
	r.PostForm = postedData
	form = NewForm(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows post does not have required fields, when it does")
	}
}

func Test_check(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "password is required")
	if form.Valid() {
		t.Error("Valid() returns false and it should be true when calling check") //valid must return false because have password error message
	}
}

func Test_GetError(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "password is required")
	s := form.Errors.Get("password") //1 st scenario
	if len(s) == 0 {
		t.Error("should have returned error from Get() but do not")
	}
	s = form.Errors.Get("hdbfshb") //2 nd scenario
	if len(s) != 0 {
		t.Error("should have not returned error but do")
	}
}

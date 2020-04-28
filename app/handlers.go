package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// render determined template
func render(template string, w http.ResponseWriter) error {
	f, e := ioutil.ReadFile(template)
	if e != nil {
		return e
	}
	fmt.Fprint(w, string(f))
	return nil
}

// Hsign "/"
func Hsign(w http.ResponseWriter, r *http.Request) {
	e := render("static/template/login_signup.html", w)
	if e != nil {
		http.Error(w, "internal server error", 500)
	}
}

// Hprofile "/profile"
func Hprofile(w http.ResponseWriter, r *http.Request) {
	e := render("static/template/profile.html", w)
	if e != nil {
		http.Error(w, "internal server error", 500)
	}
}

// Hverification "/verification"
func Hverification(w http.ResponseWriter, r *http.Request) {
	e := render("static/template/verification.html", w)
	if e != nil {
		http.Error(w, "internal server error", 500)
	}
}

// Hcontact "/contact"
func Hcontact(w http.ResponseWriter, r *http.Request) {
	e := render("static/template/contact.html", w)
	if e != nil {
		http.Error(w, "internal server error", 500)
	}
}

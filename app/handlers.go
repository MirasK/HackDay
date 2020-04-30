package app

import (
	"encoding/json"
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

// Hsign "/" && "/sign"
func Hsign(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ok := CheckIsLogged(w, r)
		if !ok {
			e = render("static/template/login_signup.html", w)
			if e != nil {
				WriteLog(e.Error())
				http.Error(w, "internal server error", 500)
			}
			return
		}
		http.Redirect(w, r, "/profile", 302)
		return
	} else if r.Method == "POST" {
		data := struct {
			Msg  string `json:"msg"`
			Type string `json:"type"`
			E    string `json:"err"`
		}{"redirect", "in", ""}
		sign := r.FormValue("sign")
		if sign == "in" {
			e = signIn(w, r)
		} else {
			data.Type = "up"
			e = signUp(w, r)
		}
		if e != nil {
			WriteLog(e.Error())
			data.E = e.Error()
			data.Msg = ""
		}
		js, e := json.Marshal(data)
		if e != nil {
			WriteLog(e.Error())
			http.Error(w, "internal server error", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// Hlogout "/logout"
func Hlogout(w http.ResponseWriter, r *http.Request) {
	e = logout(w, r)
	if e != nil {
		WriteLog(e.Error())
	}
}

// Hprofile "/profile"
func Hprofile(w http.ResponseWriter, r *http.Request) {
	ok := CheckIsLogged(w, r)
	if ok {
		e := render("static/template/profile.html", w)
		if e != nil {
			WriteLog(e.Error())
			http.Error(w, "internal server error", 500)
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// Hsettings "/profile/setting"
func Hsettings(w http.ResponseWriter, r *http.Request) {
	ok := CheckIsLogged(w, r)
	if ok {
		e := render("static/template/settings.html", w)
		if e != nil {
			WriteLog(e.Error())
			http.Error(w, "internal server error", 500)
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// Hverification "/verification"
func Hverification(w http.ResponseWriter, r *http.Request) {
	e := render("static/template/verification.html", w)
	if e != nil {
		WriteLog(e.Error())
		http.Error(w, "internal server error", 500)
	}
}

// Hcontact "/contact"
func Hcontact(w http.ResponseWriter, r *http.Request) {
	e := render("static/template/contact.html", w)
	if e != nil {
		WriteLog(e.Error())
		http.Error(w, "internal server error", 500)
	}
}

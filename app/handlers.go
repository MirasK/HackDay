package app

import (
	"hackday/db"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JSONAns type to send frontend async
type JSONAns struct {
	Msg  string      `json:"msg"`
	Type string      `json:"type"`
	E    string      `json:"err"`
	Data interface{} `json:"data"`
}

// render determined template
func render(templates []string, w http.ResponseWriter, data interface{}) error {
	t, e := template.ParseFiles(templates...)
	if e != nil {
		return e
	}

	e = t.Execute(w, data)
	if e != nil {
		return e
	}
	return nil
}

// Hsign "/"
func Hsign(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ok := CheckIsLogged(w, r)
		if !ok {
			e = render([]string{"static/template/login_signup.html"}, w, nil)
			if e != nil {
				WriteLog(e.Error())
				http.Error(w, "internal server error", 500)
			}
			return
		}
		http.Redirect(w, r, "/profile", 302)
		return
	} else if r.Method == "POST" {
		data := &JSONAns{"redirect", "in", "", nil}
		sign := r.FormValue("sign")
		if sign == "in" {
			e = signIn(w, r)
		} else {
			e = signUp(w, r)
			if e != nil {
				data.E = e.Error()
				data.Msg = ""
			} else {
				em := r.FormValue("email")
				msg := `http://localhost:8080/s/` + toCrypt(em)
				mes := "To: " + em + "\nFrom: " + "hackday20@mail.ru" + "\nSubject: Verification\n\n" +
					"You will be going to register. Follow by link, to submit your registration\n\nlink: " + msg
				e = SendMail("hackday20@mail.ru", em, mes)
				data.Type = "up"
				data.Msg = "mail sended to your email, check it"
			}
		}
		if e != nil {
			data.E = e.Error()
			data.Msg = ""
		}
		doJS(w, data)
	}
}

// HsaveUser "/s/"
func HsaveUser(w http.ResponseWriter, r *http.Request) {
	arr := strings.Split(r.URL.Path, "/")
	crypt := ""
	for _, v := range arr[2:] {
		crypt += v + "/"
	}
	crypt = crypt[:len(crypt)-1]
	email := fromCrypt(crypt)
	org := strings.Index(email, "@")
	if email[org+1] == 'g' {
		email = email[:org+1] + "gmail.com"
	} else if email[org+1] == 'm' {
		email = email[:org+1] + "mail.ru"
	}
	e = db.Update(db.GetUsersColl(), bson.D{{Key: "email", Value: email}}, bson.D{{Key: "$set", Value: bson.M{"ok": true, "expire": primitive.Null{}}}})
	if e != nil {
		WriteLog(e.Error())
		http.Error(w, "internal server error", 500)
		return
	}
	http.Redirect(w, r, "/profile/settings", 302)
}

// Hlogout "/logout"
func Hlogout(w http.ResponseWriter, r *http.Request) {
	e = logout(w, r)
	if e != nil {
		WriteLog(e.Error())
		http.Error(w, "internal seerver error", 500)
		return
	}
	http.Redirect(w, r, "/", 302)
}

// Hprofile "/profile"
func Hprofile(w http.ResponseWriter, r *http.Request) {
	ok := CheckIsLogged(w, r)
	if ok {
		e = render([]string{"static/template/profile.html"}, w, nil)
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
	if r.Method == "GET" {
		ok := CheckIsLogged(w, r)
		if ok {
			e = render([]string{"static/template/settings.html"}, w, nil)
			if e != nil {
				WriteLog(e.Error())
				http.Error(w, "internal server error", 500)
			}
		} else {
			http.Redirect(w, r, "/", 302)
		}
	} else if r.Method == "POST" {

	}
}

// Hcontact "/contact"
func Hcontact(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		e = render([]string{"static/template/contact.html"}, w, nil)
		if e != nil {
			WriteLog(e.Error())
			http.Error(w, "internal server error", 500)
		}
	} else if r.Method == "POST" {
		data := &JSONAns{"Sended", "Contact", "", ""}
		username := r.FormValue("username")
		email := r.FormValue("email")
		msg := r.FormValue("text")
		mes := "To: " + "miron.arystan@mail.ru" + "\nFrom: " + "hackday20@mail.ru" + "\nSubject: " + email + "(" + username + ") sended mail\n\n" + msg
		e = SendMail("hackday20@mail.ru", "miron.arystan@mail.ru", mes)
		if e != nil {
			data.E = e.Error()
			data.Msg = ""
		}
		doJS(w, data)
	}
}

// Hphoto handle "/profile/change-photo"
func Hphoto(w http.ResponseWriter, r *http.Request) {

}

// Hforgot handle "/forgot"
func Hforgot(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		e = render([]string{"static/template/forgot.html"}, w, nil)
		if e != nil {
			http.Error(w, "internal server error", 500)
			WriteLog(e.Error())
		}
		return
	} else if r.Method == "POST" {
		data := &JSONAns{"redirect", "forgot", "", ""}
		email := r.FormValue("email")
		e = checkEmail(true, email)
		if e != nil {
			data.E = e.Error()
			data.Msg = ""
		}

		msg := StringWithCharset(12)
		codes[msg] = msg

		mes := "To: " + "miron.arystan@mail.ru" + "\nFrom: " + "hackday20@mail.ru" +
			"\nSubject: Restore password\n\nEnter this code to verification and restore password\nCode: " + msg
		e = SendMail("hackday20@mail.ru", email, mes)
		doJS(w, data)
	}
}

// Hverification "/verification"
func Hverification(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		e = render([]string{"static/template/verification.html"}, w, nil)
		if e != nil {
			WriteLog(e.Error())
			http.Error(w, "internal server error", 500)
		}
	} else if r.Method == "POST" {
		code := r.FormValue("code")
		data := &JSONAns{"redirect", "verification", "", ""}

		if _, ok := codes[code]; !ok {
			data.E = "wrong code"
			data.Msg = ""
		}
		delete(codes, code)
		doJS(w, data)
	}
}

// Hrestore handle "/restore"
func Hrestore(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		e = render([]string{"static/template/restore.html"}, w, nil)
		if e != nil {
			http.Error(w, "internal server error", 500)
			WriteLog(e.Error())
		}
		return
	} else if r.Method == "POST" {
		data := &JSONAns{"redirect", "restore", "", ""}
		email := r.FormValue("email")
		password := r.FormValue("pass")
		rep := r.FormValue("repPass")

		if rep != password {
			data.E = "password is mismatch"
			data.Msg = ""
		} else {
			e = checkEmail(true, email)
			e = checkPassword(false, password, email)
		}
		if e != nil {
			data.E = e.Error()
			data.Msg = ""
		}
		pass, e := bcrypt.GenerateFromPassword([]byte(password), 4)
		if e != nil {
			http.Error(w, "internal server error", 500)
			WriteLog(e.Error())
			return
		}
		e = db.Update(db.GetUsersColl(), bson.D{{Key: "email", Value: email}}, bson.D{{Key: "$set", Value: bson.M{"password": string(pass)}}})
		if e != nil {
			http.Error(w, "internal server error", 500)
			WriteLog(e.Error())
			return
		}

		doJS(w, data)
	}
}

// HworkCreate handle "/create-work"
func HworkCreate(w http.ResponseWriter, r *http.Request) {

}

// Hworks handle "/works"
func Hworks(w http.ResponseWriter, r *http.Request) {

}

// Hwork handle "/work"
func Hwork(w http.ResponseWriter, r *http.Request) {

}

// HworkReq handle "/work/req"
func HworkReq(w http.ResponseWriter, r *http.Request) {

}

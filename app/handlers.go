package app

import (
	"errors"
	"hackday/db"
	"hackday/models"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// JSONAns type to send frontend async
type JSONAns struct {
	Msg  string      `json:"msg"`
	Type string      `json:"type"`
	E    string      `json:"err"`
	Data interface{} `json:"data"`
}

// Profile data to profile page
type Profile struct {
	ID         string
	Username   string
	University string
	Socials    primitive.A
	Email      string
	Photo      string
	IsStudent  bool
	IsUser     bool
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
			data.Type = "up"
			data.Msg = "mail sended to your email, check it"
			em := r.FormValue("email")
			msg := `https://hackday2020.herokuapp.com/s/` + toCrypt(em)
			mes := "To: " + em + "\nFrom: " + "hackday20@mail.ru" + "\nSubject: Verification\n\n" +
				"You will be going to register. Follow by link, to submit your registration\n\nlink: " + msg
			e = SendMail("hackday20@mail.ru", em, mes)

			if e != nil {
				e = errors.New("wrong email")
				data.Msg = ""
			} else {
				e = signUp(w, r)
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
	e = db.Update(db.GetUsersColl(), bson.D{{Key: "email", Value: email}}, bson.D{{Key: "$set", Value: bson.M{"ok": true, "expire": primitive.Null{}}}})
	delete(users, email)
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
		cook, _ := r.Cookie("sid")
		sid, _ := url.QueryUnescape(cook.Value)
		UpdateSession(sid)
		updateCooks(w, r)
		arr := strings.Split(r.URL.Path, "/")
		var res bson.M
		ID := ""
		isUser := true
		if arr[1] == "profile" {
			cook, _ = r.Cookie("sem")
			em, _ := url.QueryUnescape(cook.Value)
			res, e = db.GetOneByFilter(db.GetUsersColl(), bson.M{"email": em})
		} else {
			isUser = false
			iddb, _ := primitive.ObjectIDFromHex(arr[2])
			ID = arr[2]
			res, e = db.GetOneByFilter(db.GetUsersColl(), bson.M{"_id": iddb})
		}
		if e != nil {
			WriteLog(e.Error())
			http.Error(w, "wrong user", 500)
			return
		}

		isStudent := false
		if res["role"].(string) == "student" {
			isStudent = true
		} else if res["role"].(string) != "student" && !isUser {
			http.Error(w, "wrong user", 500)
			return
		}

		data := &Profile{ID, res["username"].(string), res["appertain"].(string), res["socials"].(primitive.A),
			res["email"].(string), res["photo"].(string), isStudent, isUser}

		e = render([]string{"static/template/profile.html"}, w, data)
		if e != nil {
			WriteLog(e.Error())
			http.Error(w, "internal server error", 500)
			return
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// Settings ...
type Settings struct {
	Tab         int
	User        *models.User
	Resume      *models.Resume
	MedCard     *models.MedCard
	StudentInfo *models.StudentInfo
	IsUser      bool
	IsStudent   bool
}

// Hsettings "/profile/settings/"
func Hsettings(w http.ResponseWriter, r *http.Request) {
	ok := CheckIsLogged(w, r)
	if ok {
		if r.Method == "GET" {
			cook, _ := r.Cookie("sid")
			sid, _ := url.QueryUnescape(cook.Value)
			UpdateSession(sid)
			updateCooks(w, r)

			arr := strings.Split(r.URL.Path, "/")
			data := &Settings{}

			switch arr[3] {
			case "i":
				data.Tab = 1
			case "a":
				data.Tab = 3
			case "s":
				data.Tab = 4
			case "m":
				data.Tab = 2
			case "r":
				data.Tab = 5
			}

			var res bson.M
			email := ""
			isUser := true
			if arr[1] == "profile" {
				cook, _ = r.Cookie("sem")
				email, _ = url.QueryUnescape(cook.Value)
				res, e = db.GetOneByFilter(db.GetUsersColl(), bson.M{"email": email})
			} else {
				isUser = false
				ID, _ := primitive.ObjectIDFromHex(arr[4])
				res, e = db.GetOneByFilter(db.GetUsersColl(), bson.M{"_id": ID})
			}
			if e != nil {
				http.Error(w, "wrong user", 500)
				WriteLog(e.Error())
				return
			}
			data.IsUser = isUser
			data.User = &models.User{
				ID:        res["_id"].(primitive.ObjectID).Hex(),
				DOB:       res["dob"].(string),
				Appertain: res["appertain"].(string),
				Username:  res["username"].(string),
				Gender:    res["gender"].(string),
				Phone:     res["phone"].(string),
				Email:     email,
			}
			if res["role"] == "student" {
				data.IsStudent = true
			}

			email = res["email"].(string)
			if data.IsStudent {
				res, e = db.GetOneByFilter(db.GetResumesColl(), bson.M{"email": email})
				if res != nil {
					data.Resume = &models.Resume{
						ID:         res["_id"].(primitive.ObjectID).Hex(),
						Skills:     res["skills"].(primitive.A),
						WhereWorks: res["whereWorks"].(primitive.A),
						AboutMe:    res["aboutMe"].(string),
						Date:       res["date"].(string),
						Link:       res["link"].(string),
					}
				}

				res, e = db.GetOneByFilter(db.GetMedCardsColl(), bson.M{"email": email})
				if res != nil {
					data.MedCard = &models.MedCard{
						ID:         res["_id"].(primitive.ObjectID).Hex(),
						Allergs:    res["allergs"].(primitive.A),
						Ills:       res["ills"].(primitive.A),
						Phychs:     res["phychs"].(primitive.A),
						Chrons:     res["chrons"].(primitive.A),
						BloodGroup: res["bloodGroup"].(string),
						Invalid:    res["invalid"].(string),
					}
				}

				res, e = db.GetOneByFilter(db.GetStudInfosColl(), bson.M{"email": email})
				if res != nil {
					data.StudentInfo = &models.StudentInfo{
						ID:           res["_id"].(primitive.ObjectID).Hex(),
						Sertificates: res["sertificates"].(primitive.A),
						Achievs:      res["achievs"].(primitive.A),
					}
				}
			}
			e = render([]string{"static/template/settings.html"}, w, data)
			if e != nil {
				WriteLog(e.Error())
				http.Error(w, "internal server error", 500)
				return
			}
		} else if r.Method == "POST" {
			arr := strings.Split(r.URL.Path, "/")
			tab := -1
			switch arr[3] {
			case "i":
				tab = 1
			case "a":
				tab = 3
			case "s":
				tab = 4
			case "m":
				tab = 2
			case "r":
				tab = 5
			case "acc":
				tab = 6
			}
			data := &JSONAns{"updated!", "settings", "", ""}
			e = updateUser(tab, arr[4], r)
			if e != nil {
				data.E = e.Error()
				data.Msg = ""
			}
			doJS(w, data)
		}
	} else {
		http.Redirect(w, r, "/", 302)
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

// HphotoAndSocials handle "/profile/ch"
func HphotoAndSocials(w http.ResponseWriter, r *http.Request) {
	ok := CheckIsLogged(w, r)
	if ok {
		cook, _ := r.Cookie("sid")
		sid, _ := url.QueryUnescape(cook.Value)
		UpdateSession(sid)
		updateCooks(w, r)
		if r.Method == "GET" {
			cook, _ := r.Cookie("sem")
			sem, _ := url.QueryUnescape(cook.Value)
			res, _ := db.GetOneByFilter(db.GetUsersColl(), bson.M{"email": sem})
			data := &models.User{
				ID:     res["_id"].(primitive.ObjectID).Hex(),
				Social: res["socials"].(primitive.A),
				Photo:  res["photo"].(string),
			}
			e = render([]string{"static/template/change.html"}, w, data)
			if e != nil {
				WriteLog(e.Error())
				http.Error(w, "internal server error", 500)
				return
			}
		} else if r.Method == "POST" {
			cook, _ := r.Cookie("sem")
			sem, _ := url.QueryUnescape(cook.Value)
			res, _ := db.GetOneByFilter(db.GetUsersColl(), bson.M{"email": sem})
			data := &JSONAns{"updated!", "change", "", ""}
			e = updateChange(res["_id"].(primitive.ObjectID), r)
			if e != nil {
				data.E = e.Error()
				data.Msg = ""
			}
			doJS(w, data)
		}
		return
	}
	http.Redirect(w, r, "/", 302)
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
	ok := CheckIsLogged(w, r)
	if ok {
		cook, _ := r.Cookie("sid")
		sid, _ := url.QueryUnescape(cook.Value)
		UpdateSession(sid)
		updateCooks(w, r)
		if r.Method == "GET" {
			e = render([]string{"static/template/create-work.html"}, w, nil)
			if e != nil {
				WriteLog(e.Error())
				http.Error(w, "internal server error", 500)
			}
		} else if r.Method == "POST" {
			data := &JSONAns{"Created!", "Work", "", ""}
			e = createWork(w, r)
			if e != nil {
				data.E = e.Error()
				data.Msg = ""
				WriteLog(e.Error())
			}
			doJS(w, data)
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// Hfilter "/filter"
func Hfilter(w http.ResponseWriter, r *http.Request) {
	ok := CheckIsLogged(w, r)
	if ok {
		cook, _ := r.Cookie("sid")
		sid, _ := url.QueryUnescape(cook.Value)
		UpdateSession(sid)
		updateCooks(w, r)
		if r.Method == "GET" {
			e = render([]string{"static/template/filter.html"}, w, nil)
			if e != nil {
				WriteLog(e.Error())
				http.Error(w, "internal server error", 500)
			}
			return
		}
	}
	http.Redirect(w, r, "/", 302)
}

// Hworks handle "/works"
func Hworks(w http.ResponseWriter, r *http.Request) {
	ok := CheckIsLogged(w, r)
	if ok {
		cook, _ := r.Cookie("sid")
		sid, _ := url.QueryUnescape(cook.Value)
		UpdateSession(sid)
		updateCooks(w, r)
		if r.Method == "GET" {
			arr := strings.Split(r.URL.Path, "/")
			filter := bson.M{}
			class := "profile/my-vacantions"
			// filter
			if arr[1] == "works" {
				class = "works"
				company := r.FormValue("c")
				requires := r.FormValue("r")
				if company != "" {
					filter["company"] = company
				}
				if requires != "" {
					filter["requirements"] = bson.M{"$elemMatch": bson.M{"$eq": requires}}
				}
			} else {
				cook, _ := r.Cookie("sem")
				em, _ := url.QueryUnescape(cook.Value)
				filter["email"] = em
			}
			// options
			opt := options.Find()
			opt.SetLimit(10)
			opt.SetSort(bson.M{"date": -1})

			res, e := db.GetAllByFilter(db.GetWorksColl(), filter, opt)
			if e != nil {
				http.Error(w, "internal server error", 500)
				WriteLog(e.Error())
				return
			}
			data := []*models.Work{}
			for _, v := range res {
				cur := &models.Work{
					ID:      v["_id"].(primitive.ObjectID).Hex(),
					Date:    v["date"].(string),
					Company: v["company"].(string),
					Info:    v["info"].(string),
					Class:   class,
				}
				data = append(data, cur)
			}

			e = render([]string{"static/template/works.html"}, w, data)
			if e != nil {
				WriteLog(e.Error())
				http.Error(w, "internal server error", 500)
			}
			return
		}
	}
	http.Redirect(w, r, "/", 302)
}

// Hwork handle "/works"
func Hwork(w http.ResponseWriter, r *http.Request) {
	ok := CheckIsLogged(w, r)
	if ok {
		cook, _ := r.Cookie("sid")
		sid, _ := url.QueryUnescape(cook.Value)
		UpdateSession(sid)
		updateCooks(w, r)
		if r.Method == "GET" {
			arr := strings.Split(r.URL.Path, "/")
			var ID primitive.ObjectID
			if arr[1] == "profile" {
				ID, e = primitive.ObjectIDFromHex(arr[3])
			} else {
				ID, e = primitive.ObjectIDFromHex(arr[2])
			}
			if e != nil {
				http.Error(w, "wrong work", 500)
				WriteLog(e.Error())
				return
			}
			res, e := db.GetOneByFilter(db.GetWorksColl(), bson.M{"_id": ID})
			if e != nil {
				http.Error(w, "wrong work", 500)
				WriteLog(e.Error())
				return
			}
			data := &models.Work{
				ID:           ID.Hex(),
				Date:         res["date"].(string),
				Company:      res["company"].(string),
				Info:         res["info"].(string),
				Phone:        res["phone"].(string),
				Email:        res["email"].(string),
				Type:         res["type"].(string),
				Requirements: res["requirements"].(primitive.A),
				Class:        "employer",
			}
			if arr[1] == "works" {
				data.Class = "student"
				cook, _ = r.Cookie("sem")
				sem, _ := url.QueryUnescape(cook.Value)
				res, _ := db.GetOneByFilter(db.GetUsersColl(), bson.M{"email": sem})
				if res != nil {
					res2, _ := db.GetOneByFilter(db.GetMsgsColl(), bson.M{"$or": bson.A{bson.M{"senderId": ID}, bson.M{"ownerId": ID}}})
					if res2 != nil {
						data.IsResponsed = true
					} else {
						data.IsResponsed = false
					}
				} else {
					WriteLog(e.Error())
					http.Error(w, "internal server error", 500)
				}
			} else {
				res, e := db.GetAllByFilter(db.GetMsgsColl(), bson.M{"$or": bson.A{bson.M{"senderId": res["_id"]}, bson.M{"ownerId": res["_id"]}}}, nil)
				if e != nil {
					panic(e)
				}
				users := []*models.User{}
				for _, v := range res {
					v2, _ := db.GetOneByFilter(db.GetUsersColl(), bson.M{"$or": bson.A{bson.M{"_id": v["senderId"]}, bson.M{"_id": v["ownerId"]}}})
					cur := &models.User{
						ID:       v2["_id"].(primitive.ObjectID).Hex(),
						Email:    v2["email"].(string),
						Username: v2["username"].(string),
						Text:     v["text"].(string),
					}
					users = append(users, cur)
				}
				data.Users = users
			}
			e = render([]string{"static/template/work.html"}, w, data)
			if e != nil {
				WriteLog(e.Error())
				http.Error(w, "internal server error", 500)
			}
		} else if r.Method == "POST" {
			data := &JSONAns{"Responsed!", "to employer", "", ""}
			tom, e := response(w, r)
			data.Type = tom
			if e != nil {
				data.E = e.Error()
				data.Msg = ""
				WriteLog(e.Error())
			}
			doJS(w, data)
		}
		return
	}
	http.Redirect(w, r, "/", 302)
}

// Hsubs "/profile/my-subscription"
func Hsubs(w http.ResponseWriter, r *http.Request) {
	ok := CheckIsLogged(w, r)
	if ok {
		cook, _ := r.Cookie("sid")
		sid, _ := url.QueryUnescape(cook.Value)
		UpdateSession(sid)
		updateCooks(w, r)
		if r.Method == "GET" {
			cook, _ := r.Cookie("sem")
			sem, _ := url.QueryUnescape(cook.Value)
			res, _ := db.GetOneByFilter(db.GetUsersColl(), bson.M{"email": sem})
			ID := res["_id"]
			msgs, _ := db.GetAllByFilter(db.GetMsgsColl(), bson.M{"$or": bson.A{bson.M{"senderId": ID}, bson.M{"ownerId": ID}}}, nil)
			data := []*models.Msg{}
			for _, v := range msgs {
				cur := &models.Msg{
					ID:     v["_id"].(primitive.ObjectID).Hex(),
					Type:   v["type"].(string),
					Status: v["status"].(bool),
					Text:   v["text"].(string),
				}

				var v2 bson.M
				if cur.Type == "to student" {
					v2, e = db.GetOneByFilter(db.GetWorksColl(), bson.M{"_id": v["senderId"]})
				} else {
					v2, e = db.GetOneByFilter(db.GetWorksColl(), bson.M{"_id": v["ownerId"]})
				}
				if e != nil {
					WriteLog(e.Error())
					http.Error(w, "internal server error", 500)
					return
				}
				cur.Info = v2["info"].(string)
				cur.Company = v2["company"].(string)
				data = append(data, cur)
			}
			e = render([]string{"static/template/my-subs.html"}, w, data)
			if e != nil {
				WriteLog(e.Error())
				http.Error(w, "internal server error", 500)
			}
			return
		}
	}
	http.Redirect(w, r, "/", 302)
}

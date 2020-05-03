package app

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"hackday/db"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// checkEmail check if email is empty or not
func checkEmail(mode bool, email string) error {
	filter := bson.M{"email": email}
	res, e := db.GetOneByFilter(db.GetUsersColl(), filter)
	if res != nil && e == nil && res["email"] == email {
		if !mode {
			return errors.New("this email is not empty")
		} else if mode && !res["ok"].(bool) {
			return errors.New("you are not valid your account")
		} else if mode && res["ok"].(bool) {
			return nil
		}
	}
	if mode {
		return errors.New("this email is not correct")
	}
	return nil
}

// checkPassword check is password is valid(up) or correct password(in)
func checkPassword(mode bool, pass, login string) error {
	if !mode {
		var validPassTmpl = regexp.MustCompile(`[A-Z]`)
		ok := validPassTmpl.MatchString(pass)
		if !ok {
			return errors.New("password must have A-Z")
		}
		validPassTmpl = regexp.MustCompile(`[a-z]`)
		ok = validPassTmpl.MatchString(pass)
		if !ok {
			return errors.New("password must have a-z(small)")
		}
		validPassTmpl = regexp.MustCompile(`[0-9]`)
		ok = validPassTmpl.MatchString(pass)
		if !ok {
			return errors.New("password must have 0-9")
		}
		if len(pass) < 8 {
			return errors.New("password must have at least 8 chars")
		}
	} else {
		filter := bson.M{"email": login}
		res, e := db.GetOneByFilter(db.GetUsersColl(), filter)
		if res != nil {
			e = bcrypt.CompareHashAndPassword([]byte(res["password"].(string)), []byte(pass))
			if e != nil {
				return errors.New("wrong password")
			}
		} else {
			return e
		}
	}
	return nil
}

// check email and password
func checkPersonData(email, password string, mode bool) error {
	e = checkEmail(mode, email)
	if e != nil {
		return e
	}
	e = checkPassword(mode, password, email)
	if e != nil {
		return e
	}
	return nil
}

// signIn check and start
// 	e = checkPersonData(email, password, true)
// 	sid, e := SessionStart(w, r, email, "sid")
// 	e = db.Update(db.GetUsersColl(), bson.D{{Key: "email", Value: email}}, bson.D{{Key: "sesId", Value: sid}})
func signIn(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	pass := r.FormValue("password")
	e = checkPersonData(email, pass, true)
	if e != nil {
		return e
	}
	sid, e := SessionStart(w, r, email, "sid")
	if e != nil {
		return e
	}
	db.Update(db.GetUsersColl(), bson.D{{Key: "email", Value: email}}, bson.D{{Key: "$set", Value: bson.M{"sesId": sid}}})

	SetSesVal("login", email, sid)
	return nil
}

// signUp check person data by
// 	if rep != pass {return}
// 	e = checkPersonData(email, password, false)
// 	sid, e := SessionStart(w, r, email, "sid")
// 	password, e := bcrypt.GenerateFromPassword([]byte(pass), 4)
// 	_, e = db.Create(db.GetUsersColl(), bson.M{"email": email, "password": password, "sesId": sid, "username": name})
func signUp(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	pass := r.FormValue("password")
	rep := r.FormValue("repPassword")
	name := r.FormValue("username")
	role := r.FormValue("role")
	if rep != pass {
		return errors.New("password mismatch")
	}
	e = checkPersonData(email, pass, false)
	if e != nil {
		return e
	}
	sid, e := SessionStart(w, r, email, "sid")
	if e != nil {
		return e
	}
	password, e := bcrypt.GenerateFromPassword([]byte(pass), 4)
	if e != nil {
		return e
	}
	_, e = db.Create(db.GetUsersColl(), bson.M{"email": email, "password": string(password), "sesId": sid, "username": name, "gender": primitive.Null{},
		"dob": primitive.Null{}, "photo": primitive.Null{}, "phone": primitive.Null{}, "userInfoId": primitive.Null{},
		"role": role, "ok": false, "expire": TimeExpire(1 * time.Hour)})
	if e != nil {
		return e
	}
	SetSesVal("login", email, sid)
	return nil
}

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func toCrypt(text string) string {
	key := "123456789012345678901234"
	return Encrypt(key, text)
}

func fromCrypt(text string) string {
	key := "123456789012345678901234"
	return Decrypt(key, text)
}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Encrypt ...
func Encrypt(key, text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return encodeBase64(ciphertext)
}

// Decrypt ...
func Decrypt(key, text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	ciphertext := decodeBase64(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext)
}

// logout ...
func logout(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("sid")
	if err != nil || cookie.Value == "" {
		return err
	}

	SessionDestroy(cookie.Value)
	if err != nil {
		return err
	}
	c := &http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	c2 := &http.Cookie{
		Name:     "sem",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
	http.SetCookie(w, c2)
	return nil
}

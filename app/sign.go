package app

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"hackday/db"
	"io"
	"net/http"
	"regexp"
	"strings"
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
	password, e := bcrypt.GenerateFromPassword([]byte(pass), 4)
	if e != nil {
		return e
	}
	if role == "student" {
		_, e = db.Create(db.GetStudInfosColl(), bson.M{"email": email, "achievs": []string{"None"}, "sertificates": []string{"None"}})
		_, e = db.Create(db.GetMedCardsColl(), bson.M{"email": email, "bloodGroup": "None", "ills": []string{"None"}, "phychs": []string{"None"},
			"chrons": []string{"None"}, "allergs": []string{"None"}, "invalid": "None"})
		_, e = db.Create(db.GetResumesColl(), bson.M{"email": email, "skills": []string{"None"}, "whereWorks": []string{"None"}, "aboutMe": "None",
			"link": "None", "date": TimeExpire(1 * time.Nanosecond)})
	}
	_, e = db.Create(db.GetUsersColl(), bson.M{"email": email, "password": string(password), "sesId": primitive.Null{}, "username": name, "gender": "None", "socials": []string{},
		"dob": "None", "photo": "None", "phone": "None", "appertain": "None", "role": role, "ok": false, "expire": TimeExpire(1 * time.Hour)})
	if e != nil {
		return e
	}
	users[email] = email
	return nil
}

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func toCrypt(text string) string {
	key := []byte("hack20MirasMiron")
	res, _ := encrypt(key, text)
	return res
}

func fromCrypt(text string) string {
	key := []byte("hack20MirasMiron")
	res, _ := decrypt(key, text)
	return res
}

// bottom 6 func is use to enc and decr message
func addBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}

	return value
}

func removeBase64Padding(value string) string {
	return strings.Replace(value, "=", "", -1)
}

func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}

func encrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	msg := pad([]byte(text))
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))
	return finalMsg, nil
}

func decrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(text))
	if err != nil {
		return "", err
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return "", errors.New("blocksize must be multipe of decoded message length")
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := unpad(msg)
	if err != nil {
		return "", err
	}

	return string(unpadMsg), nil
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

package app

import (
	"errors"
	"hackday/db"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// create a new vacantion
func createWork(w http.ResponseWriter, r *http.Request) error {
	company := r.FormValue("company")
	info := r.FormValue("info")
	typeOfWork := r.FormValue("type")
	requires := strings.Split(r.FormValue("requirements"), ",")
	phone := r.FormValue("phone")
	cook, _ := r.Cookie("sem")
	email, _ := url.QueryUnescape(cook.Value)

	_, e = db.Create(db.GetWorksColl(), bson.M{"company": company, "info": info, "date": TimeExpire(1 * time.Nanosecond),
		"phone": phone, "requirements": requires, "email": email, "type": typeOfWork, "expire": TimeExpire(30 * 24 * time.Hour)})
	if e != nil {
		return e
	}
	return nil
}

// response to one work
func response(w http.ResponseWriter, r *http.Request) (string, error) {
	text := r.FormValue("text")
	if text == "" {
		return "", errors.New("empty text")
	}
	arr := strings.Split(r.URL.Path, "/")
	owner, e := primitive.ObjectIDFromHex(arr[2])
	if e != nil {
		return "", e
	}

	typeOfMsg := r.FormValue("type")
	if typeOfMsg == "" {
		typeOfMsg = "to employer"
	}
	statusR := r.FormValue("status")
	status := false
	if statusR != "" && statusR == "true" {
		status = true
	}
	var ID primitive.ObjectID
	if statusR == "" {
		cook, _ := r.Cookie("sem")
		email, _ := url.QueryUnescape(cook.Value)
		res, e := db.GetOneByFilter(db.GetUsersColl(), bson.M{"email": email})
		if e != nil {
			return "", e
		}
		ID = res["_id"].(primitive.ObjectID)
	} else {
		ID, e = primitive.ObjectIDFromHex(r.FormValue("senderId"))
		if e != nil {
			return "", e
		}
	}

	if typeOfMsg == "to employer" {
		_, e = db.Create(db.GetMsgsColl(), bson.M{"type": typeOfMsg, "status": status, "text": text, "ownerId": owner, "senderId": ID})
	} else {
		e = db.Update(db.GetMsgsColl(), bson.D{{Key: "$or", Value: bson.A{bson.M{"senderId": ID}, bson.M{"ownerId": ID}}}},
			bson.D{{Key: "$set", Value: bson.M{"status": status, "type": typeOfMsg, "text": text, "ownerId": owner, "senderId": ID}}})
	}

	if e != nil {
		return "", e
	}
	return typeOfMsg, nil
}

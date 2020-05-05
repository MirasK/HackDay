package app

import (
	"errors"
	"hackday/db"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

func updateI(ID primitive.ObjectID, r *http.Request) error {
	username := r.FormValue("username")
	appertain := r.FormValue("appertain")
	dob := r.FormValue("dob")
	gender := r.FormValue("gender")
	phone := r.FormValue("phone")
	e = db.Update(db.GetUsersColl(), bson.D{{Key: "_id", Value: ID}}, bson.D{{Key: "$set",
		Value: bson.M{"username": username, "appertain": appertain, "dob": dob, "gender": gender, "phone": phone}}})
	if e != nil {
		return e
	}
	return nil
}

func trim(inp []string) []string {
	arr := []string{}
	for _, v := range inp {
		if v != "" {
			arr = append(arr, v)
		}
	}
	return arr
}

func updateA(ID primitive.ObjectID, r *http.Request) error {
	achievs := trim(strings.Split(r.FormValue("achievs"), ","))
	e = db.Update(db.GetStudInfosColl(), bson.D{{Key: "_id", Value: ID}}, bson.D{{Key: "$set",
		Value: bson.M{"achievs": achievs}}})
	if e != nil {
		return e
	}
	return nil
}

func updateM(ID primitive.ObjectID, r *http.Request) error {
	bloodGroup := r.FormValue("blood")
	allergs := trim(strings.Split(r.FormValue("allergs"), ","))
	ills := trim(strings.Split(r.FormValue("ills"), ","))
	phychs := trim(strings.Split(r.FormValue("phychs"), ","))
	chrons := trim(strings.Split(r.FormValue("chrons"), ","))
	invalid := r.FormValue("inval")
	e = db.Update(db.GetMedCardsColl(), bson.D{{Key: "_id", Value: ID}}, bson.D{{Key: "$set",
		Value: bson.M{"bloodGroup": bloodGroup, "allergs": allergs, "ills": ills, "phychs": phychs, "chrons": chrons, "invalid": invalid}}})
	if e != nil {
		return e
	}
	return nil
}

func updateAcc(ID primitive.ObjectID, r *http.Request) error {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == password && email == "" {
		return errors.New("not updated, it nil")
	}
	if email != "" {
		e = checkEmail(false, email)
		if e != nil {
			return e
		}
		e = db.Update(db.GetUsersColl(), bson.D{{Key: "_id", Value: ID}}, bson.D{{Key: "$set", Value: bson.M{"email": email}}})
		if e != nil {
			return e
		}
	}
	if password != "" {
		e = checkPassword(false, password, "")
		if e != nil {
			return e
		}
		ps, e := bcrypt.GenerateFromPassword([]byte(password), 4)
		if e != nil {
			return e
		}
		e = db.Update(db.GetUsersColl(), bson.D{{Key: "_id", Value: ID}}, bson.D{{Key: "$set", Value: bson.M{"password": string(ps)}}})
		if e != nil {
			return e
		}
	}

	return nil
}

func uploadFile(file *multipart.File, filename string) error {
	f, e := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if e != nil {
		return e
	}
	defer f.Close()
	io.Copy(f, *file)
	return nil
}

func updateS(ID primitive.ObjectID, r *http.Request) error {
	file, fh, e := r.FormFile("sertificate")
	link := ""
	if fh != nil {
		defer file.Close()
		if e != nil {
			return e
		}
		str := StringWithCharset(8)
		link = "static/img/sertificates/" + str + fh.Filename
		e = uploadFile(&file, link)
		if e != nil {
			return e
		}
	}
	if link != "" {
		e = db.Update(db.GetStudInfosColl(), bson.D{{Key: "_id", Value: ID}}, bson.D{{Key: "$push", Value: bson.M{"sertificates": link}}})
		if e != nil {
			return e
		}
	}
	return nil
}

func updateR(ID primitive.ObjectID, r *http.Request) error {
	skills := trim(strings.Split(r.FormValue("skills"), ","))
	whereWorks := trim(strings.Split(r.FormValue("where"), ","))
	aboutMe := r.FormValue("about")

	file, fh, e := r.FormFile("link")
	link := ""
	if fh != nil {
		defer file.Close()
		if e != nil {
			return e
		}
		str := StringWithCharset(8)
		link = "static/resume/" + str + fh.Filename
		e = uploadFile(&file, link)
		if e != nil {
			return e
		}
	}

	date := TimeExpire(1 * time.Nanosecond)
	if link != "" {
		e = db.Update(db.GetResumesColl(), bson.D{{Key: "_id", Value: ID}}, bson.D{{Key: "$set", Value: bson.M{"link": link,
			"skills": skills, "whereWorks": whereWorks, "aboutMe": aboutMe, "date": date}}})
	} else {
		e = db.Update(db.GetResumesColl(), bson.D{{Key: "_id", Value: ID}}, bson.D{{Key: "$set", Value: bson.M{"skills": skills,
			"whereWorks": whereWorks, "aboutMe": aboutMe, "date": date}}})
	}

	if e != nil {
		return e
	}
	return nil
}

func updateUser(tab int, idr string, r *http.Request) error {
	ID, e := primitive.ObjectIDFromHex(idr)
	if e != nil {
		return e
	}
	switch tab {
	case 1:
		e = updateI(ID, r)
	case 2:
		e = updateM(ID, r)
	case 3:
		e = updateA(ID, r)
	case 4:
		e = updateS(ID, r)
	case 5:
		e = updateR(ID, r)
	case 6:
		e = updateAcc(ID, r)
	}
	if e != nil {
		return e
	}

	return nil
}

func updateChange(ID primitive.ObjectID, r *http.Request) error {
	face := r.FormValue("facebook")
	linked := r.FormValue("linkedin")
	insta := r.FormValue("instagram")
	arr := []string{}
	if face != "" {
		arr = append(arr, face)
	}
	if linked != "" {
		arr = append(arr, linked)
	}
	if insta != "" {
		arr = append(arr, insta)
	}
	file, fh, e := r.FormFile("photo")
	link := ""
	if fh != nil {
		defer file.Close()
		if e != nil {
			return e
		}
		str := StringWithCharset(8)
		link = "static/img/avatars/" + str + fh.Filename
		e = uploadFile(&file, link)
		if e != nil {
			return e
		}
	}

	if link != "" {
		e = db.Update(db.GetUsersColl(), bson.D{{Key: "_id", Value: ID}}, bson.D{{Key: "$set", Value: bson.M{"photo": link, "socials": arr}}})
	} else {
		e = db.Update(db.GetUsersColl(), bson.D{{Key: "_id", Value: ID}}, bson.D{{Key: "$set", Value: bson.M{"socials": arr}}})
	}

	if e != nil {
		return e
	}
	return nil
}

package app

import (
	"fmt"
	"hackday/db"
	"net/http"
	"net/url"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// vars
var (
	e          error
	logFile    *os.File
	timeLayout = "2006-01-02 15:04:05"
)

// WriteLog write to logs file
// 	logFile.WriteString(time.Now().Format(timeLayout) + "| " + msg + "\n")
func WriteLog(msg string) {
	logFile.WriteString(time.Now().Format(timeLayout) + "| " + msg + "\n")
}

// return new UUID for sessions
// 	u1 := uuid.Must(uuid.NewV4(), err)
// 	return fmt.Sprint(u1)
func newSessID() string {
	u1 := uuid.Must(uuid.NewV4(), e)
	return fmt.Sprint(u1)
}

// CheckIsLogged check if user is logged
// 	cookie, err := r.Cookie(cookName)
// 	res, e := db.GetOneByFilter(db.GetUsersColl(), bson.M{"sesId": cookie.Value})
// 	return res["email"].(string)
func CheckIsLogged(w http.ResponseWriter, r *http.Request) bool {
	sidCook, e := r.Cookie("sid")
	semCook, e := r.Cookie("sem")
	if e != nil || sidCook.Value == "" || semCook.Value == "" {
		return false
	}
	em, e := url.QueryUnescape(semCook.Value)
	if e != nil {
		return false
	}
	res, e := db.GetOneByFilter(db.GetUsersColl(), bson.M{"email": em})
	if e != nil {
		return false
	}
	if res["sesId"] == bson.TypeNull {
		logout(w, r)
		return false
	}
	return true
}

// return time.Now().Add(add).Format(timeLayout)
func timeExpire(add time.Duration) string {
	return time.Now().Add(add).Format(timeLayout)
}

// UpdateSession update db data session
// 	e = db.Update(db.GetSessColl(), bson.D{{Key: "filename", Value: sid}}, bson.D{{Key: "expire", Value: timeExpire(1 * time.Hour)}})
// 	if e != nil {
// 		return e
// 	}
// 	return nil
func UpdateSession(sid string) error {
	e = db.Update(db.GetSessColl(), bson.D{{Key: "filename", Value: sid}}, bson.D{{Key: "$set", Value: bson.M{"expire": timeExpire(1 * time.Hour)}}})
	if e != nil {
		return e
	}
	return nil
}

// SessionStart ...
func SessionStart(w http.ResponseWriter, r *http.Request, login, cookieName string) (string, error) {
	cookie, e := r.Cookie(cookieName)
	sid := ""
	if e != nil || cookie.Value == "" {
		res, e := db.GetOneByFilter(db.GetUsersColl(), bson.M{"email": login})
		if res == nil || res != nil && res["sesId"] == nil {
			sid = newSessID()
			e = SessionInit(sid)
			_, e = db.Create(db.GetSessColl(), bson.M{"filename": sid, "expire": timeExpire(1 * time.Hour)})
			if e != nil {
				return "", e
			}
		} else {
			sid = res["sesId"].(string)
			e = SessionRead(sid)
			e = UpdateSession(sid)
			if e != nil {
				return "", e
			}
		}
	} else {
		sid, _ = url.QueryUnescape(cookie.Value)
		e = SessionRead(sid)
		if e != nil {
			return "", e
		}
		e = UpdateSession(sid)
		if e != nil {
			return "", e
		}
	}
	sidCook := http.Cookie{Name: cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: 3600}
	semCook := http.Cookie{Name: "sem", Value: url.QueryEscape(login), Path: "/", HttpOnly: true, MaxAge: 3600}
	http.SetCookie(w, &sidCook)
	http.SetCookie(w, &semCook)
	return sid, nil
}

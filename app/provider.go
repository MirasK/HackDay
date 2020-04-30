/*
	This file define and describe provider interface. Provider is used to control and audit sessions: Session: init, read, destroy, destroy by expire
*/

package app

import (
	"hackday/db"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SessionInit init new session
// 	_, e = os.Create("sess/" + sid)
// 	if e != nil {
// 		return e
// 	}
// 	return nil
func SessionInit(sid string) error {
	_, e = os.Create("sess/" + sid)
	if e != nil {
		return e
	}
	return nil
}

// SessionRead return/create session
// 	_, e = os.Open("sess/" + sid)
// 	if e != nil {
// 		_, e = os.Create("sess/" + sid)
// 		if e != nil {
// 			return e
// 		}
// 	}
// 	return nil
func SessionRead(sid string) error {
	_, e = os.Open("sess/" + sid)
	if e != nil {
		_, e = os.Create("sess/" + sid)
		if e != nil {
			return e
		}
	}
	return nil
}

// SessionDestroy delete session by sid
// 	e = os.Remove("sess/" + sid)
// 	e = db.Delete(db.GetSessColl(), bson.D{{Key: "filename", Value: sid}})
// 	e = db.Update(db.GetUsersColl(), bson.D{{Key: "sesId", Value: sesID}}, bson.D{{Key: "sesId", Value: bson.TypeNull}})
func SessionDestroy(sid string) {
	// remove from filesystem and provcontrol
	e = os.Remove("sess/" + sid)
	if e != nil {
		WriteLog(e.Error())
		return
	}

	// remove from Sessions
	e = db.Delete(db.GetSessColl(), bson.D{{Key: "filename", Value: sid}})
	if e != nil {
		WriteLog(e.Error())
		return
	}

	// update Users
	e = db.Update(db.GetUsersColl(), bson.D{{Key: "sesId", Value: sid}}, bson.D{{Key: "$set", Value: bson.M{"sesId": primitive.Null{}}}})
	if e != nil {
		WriteLog(e.Error())
		return
	}
}

// SessionGC delete expired session
// 	res, e := db.GetAllByFilter(db.GetSessColl(), bson.M{"expire": bson.M{"$lte": timeExpire(1 * time.Nanosecond)}})
// 	if e != nil {
// 		logFile.WriteString(time.Now().Format(timeLayout) + "| " + e.Error() + "\n")
// 		return
// 	}
// 	for _, v := range res {
// 		go SessionDestroy(v["filename"].(string))
// 	}
func SessionGC() {
	res, e := db.GetAllByFilter(db.GetSessColl(), bson.M{"expire": bson.M{"$lte": timeExpire(1 * time.Nanosecond)}})
	if e != nil {
		WriteLog(e.Error())
		return
	}
	for _, v := range res {
		go SessionDestroy(v["filename"].(string))
	}
}

// CheckPerMin call SessionGC per minute that delete expired sessions
// 	for {
// 		timer := time.NewTimer(1 * time.Minute)
// 		<-timer.C
// 		go SessionGC()
// 	}
func CheckPerMin() {
	for {
		timer := time.NewTimer(1 * time.Minute)
		<-timer.C
		go SessionGC()
	}
}

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MedCard is users medcards
// type MedCard struct {
// 	ID         primitive.ObjectID `json:"_id"`
// 	BloodGroup string             `json:"bloodGroup"`
// 	Ills       []string           `json:"ills"`
// 	Phobies    []string           `json:"phobies"`
// }

// MedCard is users medcards
type MedCard struct {
	ID         string
	BloodGroup string
	Allergs    primitive.A
	Ills       primitive.A
	Phychs     primitive.A
	Chrons     primitive.A
	Invalid    string
}

// Msg is message between employee & employer
// type Msg struct {
// 	ID       primitive.ObjectID `json:"_id"`
// 	Type     string             `json:"type"`
// 	Status   bool               `json:"status"`
// 	Text     string             `json:"text"`
// 	OwnerID  primitive.ObjectID `json:"ownerId"`
// 	SenderID primitive.ObjectID `json:"senderId"`
// }

// Msg one msg
type Msg struct {
	ID      string
	Type    string
	Status  bool
	Text    string
	Info    string
	Company string
}

// Resume is users resume
// type Resume struct {
// 	ID         primitive.ObjectID `json:"_id"`
// 	Skills     []string           `json:"skills"`
// 	WhereWorks []string           `json:"whereWorks"`
// 	AboutMe    string             `json:"aboutMe"`
// 	Date       time.Time          `json:"date"`
// 	Link       string             `json:"link"`
// }

// Resume is users resume
type Resume struct {
	ID         string
	Skills     primitive.A
	WhereWorks primitive.A
	AboutMe    string
	Date       string
	Link       string
}

// Session in db
// type Session struct {
// 	ID       primitive.ObjectID `json:"_id"`
// 	FileName string             `json:"filename"`
// 	Expire   string             `json:"expire"`
// }

// StudentInfo additional student info
// type StudentInfo struct {
// 	ID           primitive.ObjectID `json:"_id"`
// 	MedID        primitive.ObjectID `json:"medId"`
// 	ResumeID     primitive.ObjectID `json:"resumeId"`
// 	Sertificates []string           `json:"sertificates"`
// 	Achievs      []string           `json:"achievs"`
// }

// StudentInfo additional student info
type StudentInfo struct {
	ID           string
	Sertificates primitive.A
	Achievs      primitive.A
}

// Work is one vacantion
// type Work struct {
// 	ID           primitive.ObjectID `json:"_id"`
// 	Date         time.Time          `json:"date"`
// 	Company      string             `json:"company"`
// 	Info         string             `json:"info"`
// 	Requirements []string           `json:"requirements"`
// 	Type         string             `json:"type"`
// 	Phone        string             `json:"phone"`
// 	Email        string             `json:"email"`
// }

// Work one work
type Work struct {
	ID           string
	Date         string
	Phone        string
	Requirements primitive.A
	Email        string
	Type         string
	Company      string
	Info         string
	Class        string
	IsResponsed  bool
	Users        []*User
}

// User in one user
// type User struct {
// 	ID       primitive.ObjectID `json:"_id"`
// 	Photo    string             `json:"photo"`
// 	Username string             `json:"username"`
// 	Email    string             `json:"email"`
// 	DOB      time.Time          `json:"dob"`
// 	Password string             `json:"password"`
// 	Phone    string             `json:"phone"`
// 	Role     string             `json:"role"`
// 	UserInfo primitive.ObjectID `json:"userInfoId,omitempty"`
// 	SesID    primitive.ObjectID `json:"sesId,omitempty"`
// }

// User in one user
type User struct {
	ID        string
	Photo     string
	Username  string
	Email     string
	DOB       string
	Gender    string
	Appertain string
	Social    primitive.A
	Phone     string
	Role      string
	Text      string
}

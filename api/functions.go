package api

import (
	"errors"
	"fmt"
	"hackday/app"
	"hackday/db"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret []byte = []byte("hackday")
)

// APISchema graphql
var APISchema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

// ExecuteQuery ...
func ExecuteQuery(query string, schema graphql.Schema, token string) *graphql.Result {
	if token == "" {
		return &graphql.Result{Data: "authorize before\nsend auth 'email' and 'password' to /auth?email=email&password=password"}
	}
	e := ValidateJWT(token)
	if e != nil {
		return &graphql.Result{Data: e.Error()}
	}

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	return result
}

// CreateTokenEndpoint create api token
func CreateTokenEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	email := r.FormValue("email")
	pass := r.FormValue("password")
	res, e := db.GetOneByFilter(db.GetUsersColl(), bson.M{"email": email})
	if e != nil {
		app.WriteLog(e.Error())
		w.Write([]byte(`{ "error": "` + e.Error() + `" }`))
		return
	}
	e = bcrypt.CompareHashAndPassword([]byte(res["password"].(string)), []byte(pass))
	if e != nil {
		app.WriteLog(e.Error())
		w.Write([]byte(`{ "error": "` + e.Error() + `" }`))
		return
	}
	timeEx := app.TimeExpire(1 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"expire": timeEx,
	})
	tokenString, e := token.SignedString(jwtSecret)
	if e != nil {
		app.WriteLog(e.Error())
		w.Write([]byte(`{ "error": "` + e.Error() + `" }`))
		return
	}
	_, e = db.Create(db.GetTokenColl(), bson.M{"token": tokenString, "expire": timeEx})
	if e != nil {
		w.Write([]byte(`{ "error": "` + e.Error() + `" }`))
		return
	}
	w.Write([]byte(`{ "token": "` + tokenString + `", "expiration time": "` + timeEx + `" }`))
}

// ValidateJWT check token
func ValidateJWT(t string) error {
	if t == "" {
		return errors.New("Authorization token must be present")
	}
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return jwtSecret, nil
	})
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	expire, ok := token.Claims.(jwt.MapClaims)["expire"]
	if ok && app.TimeExpire(time.Nanosecond) > expire.(string) {
		return errors.New("token is not valid more")
	}
	return nil
}

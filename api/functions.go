package api

import (
	"errors"
	"fmt"
	"hackday/app"
	"hackday/db"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
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
	_, e := ValidateJWT(token)
	if e != nil {
		return &graphql.Result{Data: "authorize before\nsend auth 'email' and 'password' to /auth?email=email&password=password"}
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": email,
		"password": pass,
	})
	tokenString, e := token.SignedString(jwtSecret)
	if e != nil {
		app.WriteLog(e.Error())
		w.Write([]byte(`{ "error": "` + e.Error() + `" }`))
		return
	}
	w.Write([]byte(`{ "token": "` + tokenString + `" }`))
}

// ValidateJWT check token
func ValidateJWT(t string) (interface{}, error) {
	if t == "" {
		return nil, errors.New("Authorization token must be present")
	}
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return jwtSecret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var decodedToken interface{}
		mapstructure.Decode(claims, &decodedToken)
		return decodedToken, nil
	}
	return nil, errors.New("Invalid authorization token")

}

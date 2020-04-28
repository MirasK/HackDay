package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hackday/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client *mongo.Client
)

// APISchema graphql
var APISchema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

// AppSchema graphql
var AppSchema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

var jwtSecret []byte = []byte("hackday")

// SetClient get client from main.go
// 	client = cl
func SetClient(cl *mongo.Client) {
	client = cl
}

// this func is get one document from mongo server
//	filter := bson.M{"_id": ID}
// 	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
// 	e := collection.FindOne(ctx, filter).Decode(&filler)
// 	if e != nil {
// 		log.Fatal(e)
// 	}
func queryGetByID(collection *mongo.Collection, id interface{}) (primitive.M, error) {
	filter := bson.M{"_id": id}
	var result bson.M
	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cncl()
	e := collection.FindOne(ctx, filter).Decode(&result)
	if e != nil {
		return nil, e
	}
	return result, nil
}

// this func get all documents from collection
// 	var arr []interface{}
// 	var ids []interface{}
// 	ids = append(ids, result["_id"])
// 	arr = append(arr, filler)
func queryGetAll(collection *mongo.Collection) ([]primitive.M, error) {
	var arr []primitive.M
	ctx, cncl := context.WithTimeout(context.Background(), 30*time.Second)
	defer cncl()
	cur, e := collection.Find(ctx, bson.D{})
	if e != nil {
		return nil, e
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		e = cur.Decode(&result)
		if e != nil {
			return nil, e
		}
		arr = append(arr, result)
	}
	if e := cur.Err(); e != nil {
		return nil, e
	}
	return arr, nil
}

// ExecuteQuery ...
func ExecuteQuery(query string, schema graphql.Schema, token string) *graphql.Result {
	_, err := ValidateJWT(token)
	if err != nil {
		return &graphql.Result{Data: "authorize before\\nsend auth data to /auth"}
	}

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	return result
}

// CreateTokenEndpoint create api token
func CreateTokenEndpoint(response http.ResponseWriter, request *http.Request) {
	var user models.User
	_ = json.NewDecoder(request.Body).Decode(&user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})
	tokenString, error := token.SignedString(jwtSecret)
	if error != nil {
		fmt.Println(error)
	}
	response.Header().Set("content-type", "application/json")
	response.Write([]byte(`{ "token": "` + tokenString + `" }`))
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

// create a new document in db
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	res, e := collection.InsertOne(ctx, data)
// 	if e != nil {
// 		return nil, e
// 	}
// 	id := res.InsertedID
// 	return id, nil
func create(collection *mongo.Collection, data bson.M) (interface{}, error) {
	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cncl()
	res, e := collection.InsertOne(ctx, data)
	if e != nil {
		return nil, e
	}
	id := res.InsertedID
	return id, nil
}

// update document by filter and set data
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	_, e := collection.UpdateMany(ctx, filter, data)
// 	if e != nil {
// 		return e
// 	}
// 	return nil
func update(collection *mongo.Collection, filter, data bson.D) error {
	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cncl()
	res, e := collection.UpdateMany(ctx, filter, data)
	if e != nil {
		return e
	}
	if res.ModifiedCount == 0 {
		return errors.New("not updated or not exist")
	}
	return nil
}

// delete all matchs
// 	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cncl()
// 	_, e := collection.DeleteMany(ctx, filter)
// 	if e != nil {
// 		return e
// 	}
// 	return nil
func delete(collection *mongo.Collection, filter bson.D) error {
	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cncl()
	res, e := collection.DeleteMany(ctx, filter)
	if e != nil {
		return e
	}
	if res.DeletedCount == 0 {
		return errors.New("not deleted or not exist")
	}
	return nil
}

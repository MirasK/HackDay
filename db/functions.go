package db

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName = "HHCustom"
var client *mongo.Client
var e error

// Conn is conn to server db
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(URL))
// 	return client, nil
func Conn() error {
	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cncl()
	client, e = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://miron:89f90g@webdev-mk9fx.mongodb.net/test?retryWrites=true&w=majority"))
	if e != nil {
		return e
	}
	return nil
}

// GetOneByFilter get from db one document
// 	ctx, cncl := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cncl()
// 	var result bson.M
// 	e := collection.FindOne(ctx, filter).Decode(&result)
// 	if e != nil {
// 		return nil, e
// 	}
// 	return result, nil
func GetOneByFilter(collection *mongo.Collection, filter bson.M) (primitive.M, error) {
	ctx, cncl := context.WithTimeout(context.Background(), 30*time.Second)
	defer cncl()
	var result bson.M
	e := collection.FindOne(ctx, filter).Decode(&result)
	if e != nil {
		return nil, e
	}
	return result, nil
}

// GetAllByFilter this func get all documents from collection
// 	var arr []primitive.M
// 	cur, e := collection.Find(ctx, filter)
// 	e = cur.Decode(&result)
// 	arr = append(arr, result)
func GetAllByFilter(collection *mongo.Collection, filter bson.M, opt *options.FindOptions) ([]primitive.M, error) {
	var arr []primitive.M
	ctx, cncl := context.WithTimeout(context.Background(), 30*time.Second)
	defer cncl()
	cur, e := collection.Find(ctx, filter, opt)
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

// Create a new document in db
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	res, e := collection.InsertOne(ctx, data)
// 	if e != nil {
// 		return nil, e
// 	}
// 	id := res.InsertedID
// 	return id, nil
func Create(collection *mongo.Collection, data bson.M) (interface{}, error) {
	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cncl()
	res, e := collection.InsertOne(ctx, data)
	if e != nil {
		return nil, e
	}
	id := res.InsertedID
	return id, nil
}

// Update document by filter and set data
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	_, e := collection.UpdateMany(ctx, filter, data)
// 	if e != nil {
// 		return e
// 	}
// 	return nil
func Update(collection *mongo.Collection, filter, data bson.D) error {
	ctx, cncl := context.WithTimeout(context.Background(), 30*time.Second)
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

// Delete all matchs
// 	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cncl()
// 	_, e := collection.DeleteMany(ctx, filter)
// 	if e != nil {
// 		return e
// 	}
// 	return nil
func Delete(collection *mongo.Collection, filter bson.D) error {
	ctx, cncl := context.WithTimeout(context.Background(), 30*time.Second)
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

// GetUsersColl return Users collection
func GetUsersColl() *mongo.Collection {
	return client.Database(dbName).Collection("Users")
}

// GetSessColl return Sessions collection
func GetSessColl() *mongo.Collection {
	return client.Database(dbName).Collection("Sessions")
}

// GetWorksColl return Works collection
func GetWorksColl() *mongo.Collection {
	return client.Database(dbName).Collection("Works")
}

// GetMedCardsColl return MedCards collection
func GetMedCardsColl() *mongo.Collection {
	return client.Database(dbName).Collection("MedCards")
}

// GetResumesColl return Resumes collection
func GetResumesColl() *mongo.Collection {
	return client.Database(dbName).Collection("Resumes")
}

// GetStudInfosColl return StudentInfos collection
func GetStudInfosColl() *mongo.Collection {
	return client.Database(dbName).Collection("StudentInfos")
}

// GetMsgsColl return Msgs collection
func GetMsgsColl() *mongo.Collection {
	return client.Database(dbName).Collection("Msgs")
}

// GetTokenColl return Msgs collection
func GetTokenColl() *mongo.Collection {
	return client.Database(dbName).Collection("Tokens")
}

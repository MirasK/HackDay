package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Conn is conn to server db
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	client, _ := mongo.Connect(ctx, options.Client().ApplyURI())
// 	return client, nil
func Conn() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, e := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://miron:89f90g@webdev-mk9fx.mongodb.net/test?retryWrites=true&w=majority"))
	if e != nil {
		return nil, e
	}
	return client, nil
}

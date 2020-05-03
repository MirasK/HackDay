package api

import (
	"log"

	"hackday/db"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// APIqueryType is graphql query handler for api requests
var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"work": &graphql.Field{
				Type:        WorkType,
				Description: "Get work by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					ID, e := primitive.ObjectIDFromHex(id)
					if e != nil {
						return nil, e
					}
					if ok {
						res, e := db.GetOneByFilter(db.GetWorksColl(), bson.M{"_id": ID})
						if e != nil {
							return nil, e
						}
						return res, nil
					}

					return nil, nil
				},
			},
			"works": &graphql.Field{
				Type:        graphql.NewList(WorkType),
				Description: "Get works list",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					arr, e := db.GetAllByFilter(db.GetWorksColl(), bson.M{})
					if e != nil {
						return nil, e
					}
					return arr, nil
				},
			},
			"user": &graphql.Field{
				Type:        UserType,
				Description: "Get user by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					ID, e := primitive.ObjectIDFromHex(id)
					if e != nil {
						return nil, e
					}
					if ok {
						res, e := db.GetOneByFilter(db.GetUsersColl(), bson.M{"_id": ID})
						if e != nil {
							return nil, e
						}
						return res, nil
					}
					return nil, nil
				},
			},
			"users": &graphql.Field{
				Type:        graphql.NewList(UserType),
				Description: "Get users list",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					arr, e := db.GetAllByFilter(db.GetUsersColl(), bson.M{})
					if e != nil {
						log.Fatal(e)
						return nil, e
					}
					return arr, nil
				},
			},
		},
	},
)

package api

import (
	"log"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// APIqueryType is graphql query handler for api requests
var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"msg": &graphql.Field{
				Type:        MsgType,
				Description: "Get message by id",
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
						collection := client.Database("HHCustom").Collection("Msgs")
						res, e := queryGetByID(collection, ID)
						if e != nil {
							return nil, e
						}
						return res, nil
					}
					return nil, nil
				},
			},
			"msgs": &graphql.Field{
				Type:        graphql.NewList(MsgType),
				Description: "Get message list",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					collection := client.Database("HHCustom").Collection("Msgs")
					arr, e := queryGetAll(collection)
					if e != nil {
						return nil, e
					}
					return arr, nil
				},
			},
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
						collection := client.Database("HHCustom").Collection("Works")
						res, e := queryGetByID(collection, ID)
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
					collection := client.Database("HHCustom").Collection("Works")
					arr, e := queryGetAll(collection)
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
						collection := client.Database("HHCustom").Collection("Users")
						res, e := queryGetByID(collection, ID)
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
					collection := client.Database("HHCustom").Collection("Users")
					arr, e := queryGetAll(collection)
					if e != nil {
						log.Fatal(e)
						return nil, e
					}
					return arr, nil
				},
			},
			"medcard": &graphql.Field{
				Type:        MedCardType,
				Description: "Get medcard by id",
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
						collection := client.Database("HHCustom").Collection("MedCards")
						res, e := queryGetByID(collection, ID)
						if e != nil {
							return nil, e
						}
						return res, nil
					}
					return nil, nil
				},
			},
			"medcards": &graphql.Field{
				Type:        graphql.NewList(MedCardType),
				Description: "Get medcards list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					collection := client.Database("HHCustom").Collection("MedCards")
					arr, e := queryGetAll(collection)
					if e != nil {
						return nil, e
					}
					return arr, nil
				},
			},
			"resume": &graphql.Field{
				Type:        ResumeType,
				Description: "Get resume by id",
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
						collection := client.Database("HHCustom").Collection("Resumes")
						res, e := queryGetByID(collection, ID)
						if e != nil {
							return nil, e
						}
						return res, nil
					}
					return nil, nil
				},
			},
			"resumes": &graphql.Field{
				Type:        graphql.NewList(ResumeType),
				Description: "Get resumes list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					collection := client.Database("HHCustom").Collection("Resumes")
					arr, e := queryGetAll(collection)
					if e != nil {
						return nil, e
					}
					return arr, nil
				},
			},
		},
	},
)

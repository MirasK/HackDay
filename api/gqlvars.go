package api

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// graphql vars
var (
	MedCardType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "MedCard",
			Fields: graphql.Fields{
				"_id": &graphql.Field{
					Type: graphql.ID,
				},
				"bloodGroup": &graphql.Field{
					Type: graphql.String,
				},
				"ills": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
				"phobies": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
			},
		},
	)

	MsgType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Msg",
			Fields: graphql.Fields{
				"_id": &graphql.Field{
					Type: graphql.ID,
				},
				"type": &graphql.Field{
					Type: graphql.String,
				},
				"status": &graphql.Field{
					Type: graphql.Boolean,
				},
				"text": &graphql.Field{
					Type: graphql.String,
				},
				"ownerId": &graphql.Field{
					Type: graphql.ID,
				},
				"senderId": &graphql.Field{
					Type: graphql.ID,
				},
			},
		},
	)

	ResumeType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Resume",
			Fields: graphql.Fields{
				"_id": &graphql.Field{
					Type: graphql.ID,
				},
				"skills": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
				"whereWorks": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
				"aboutMe": &graphql.Field{
					Type: graphql.String,
				},
				"date": &graphql.Field{
					Type: graphql.DateTime,
				},
				"link": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	SessionType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Msg",
			Fields: graphql.Fields{
				"_id": &graphql.Field{
					Type: graphql.ID,
				},
				"filename": &graphql.Field{
					Type: graphql.String,
				},
				"expire": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	UserInfoType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "UserInfo",
			Fields: graphql.Fields{
				"_id": &graphql.Field{
					Type: graphql.ID,
				},
				"medId": &graphql.Field{
					Type: graphql.ID,
				},
				"resumeId": &graphql.Field{
					Type: graphql.ID,
				},
				"sertificates": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
				"achievs": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
			},
		},
	)

	WorkType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Work",
			Fields: graphql.Fields{
				"_id": &graphql.Field{
					Type: graphql.ID,
				},
				"date": &graphql.Field{
					Type: graphql.DateTime,
				},
				"company": &graphql.Field{
					Type: graphql.String,
				},
				"info": &graphql.Field{
					Type: graphql.String,
				},
				"requirements": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
				"type": &graphql.Field{
					Type: graphql.String,
				},
				"phone": &graphql.Field{
					Type: graphql.String,
				},
				"email": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	UserType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "User",
			Fields: graphql.Fields{
				"_id": &graphql.Field{
					Type: graphql.ID,
				},
				"dob": &graphql.Field{
					Type: graphql.DateTime,
				},
				"username": &graphql.Field{
					Type: graphql.String,
				},
				"email": &graphql.Field{
					Type: graphql.String,
				},
				"photo": &graphql.Field{
					Type: graphql.String,
				},
				"phone": &graphql.Field{
					Type: graphql.String,
				},
				"role": &graphql.Field{
					Type: graphql.String,
				},
				"userInfoId": &graphql.Field{
					Type: UserInfoType,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Source.(primitive.M)[p.Info.FieldName]
						if ok && id != "null" {
							collection := client.Database("HHCustom").Collection("StudentInfos")
							res, e := queryGetByID(collection, id)
							if e != nil {
								return nil, e
							}
							return res, nil
						}
						return nil, nil
					},
				},
				"sesId": &graphql.Field{
					Type: SessionType,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Source.(primitive.M)[p.Info.FieldName]
						if ok && id != "null" {
							collection := client.Database("HHCustom").Collection("Sessions")
							res, e := queryGetByID(collection, id)
							if e != nil {
								return nil, e
							}
							return res, nil
						}
						return nil, nil
					},
				},
			},
		},
	)
)

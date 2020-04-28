package api

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createMedCard": &graphql.Field{
				Type:        MedCardType,
				Description: "Create new MedCard",
				Args: graphql.FieldConfigArgument{
					"bloodGroup": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"ills": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"phobies": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, e := create(client.Database("HHCustom").Collection("MedCards"),
						bson.M{"bloodGroup": p.Args["bloodGroup"], "ills": p.Args["ills"], "phobies": p.Args["phobies"]})
					if e != nil {
						return nil, e
					}
					return bson.M{"_id": id, "bloodGroup": p.Args["bloodGroup"], "ills": p.Args["ills"], "phobies": p.Args["phobies"]}, nil
				},
			},
			"createMsg": &graphql.Field{
				Type:        MsgType,
				Description: "Create new Msg",
				Args: graphql.FieldConfigArgument{
					"type": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"status": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Boolean),
					},
					"text": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"ownerId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"senderId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					IDO, e := primitive.ObjectIDFromHex(p.Args["ownerId"].(string))
					IDS, e := primitive.ObjectIDFromHex(p.Args["senderId"].(string))
					if e != nil {
						return nil, e
					}
					id, e := create(client.Database("HHCustom").Collection("Msgs"),
						bson.M{"type": p.Args["type"], "status": p.Args["status"], "text": p.Args["text"],
							"ownerId": IDO, "senderId": IDS})
					if e != nil {
						return nil, e
					}
					return bson.M{"_id": id, "type": p.Args["type"], "status": p.Args["status"], "text": p.Args["text"],
						"ownerId": IDO, "senderId": IDS}, nil
				},
			},
			"createResume": &graphql.Field{
				Type:        ResumeType,
				Description: "Create new Resume",
				Args: graphql.FieldConfigArgument{
					"skills": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
					},
					"whereWorks": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
					},
					"aboutMe": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"date": &graphql.ArgumentConfig{
						Type: graphql.DateTime,
					},
					"link": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, e := create(client.Database("HHCustom").Collection("Resumes"),
						bson.M{"skills": p.Args["skills"], "whereWorks": p.Args["whereWorks"], "aboutMe": p.Args["aboutMe"],
							"date": p.Args["date"], "link": p.Args["link"]})
					if e != nil {
						return nil, e
					}
					return bson.M{"_id": id, "skills": p.Args["skills"], "whereWorks": p.Args["whereWorks"], "aboutMe": p.Args["aboutMe"],
						"date": p.Args["date"], "link": p.Args["link"]}, nil
				},
			},
			"createSession": &graphql.Field{
				Type:        SessionType,
				Description: "Create new Session",
				Args: graphql.FieldConfigArgument{
					"filename": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"expire": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, e := create(client.Database("HHCustom").Collection("Sessions"),
						bson.M{"filename": p.Args["filename"], "expire": p.Args["expire"]})
					if e != nil {
						return nil, e
					}
					return bson.M{"_id": id, "filename": p.Args["filename"], "expire": p.Args["expire"]}, nil
				},
			},
			"createUserInfo": &graphql.Field{
				Type:        UserInfoType,
				Description: "Create new StudentInfo",
				Args: graphql.FieldConfigArgument{
					"medId": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
					"resumeId": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
					"sertificates": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"achievs": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					IDM, _ := primitive.ObjectIDFromHex(p.Args["medId"].(string))
					IDR, _ := primitive.ObjectIDFromHex(p.Args["resumeId"].(string))
					id, e := create(client.Database("HHCustom").Collection("StudentInfos"),
						bson.M{"medId": IDM, "resumeId": IDR,
							"sertificates": p.Args["sertificates"], "achievs": p.Args["achievs"]})
					if e != nil {
						return nil, e
					}
					return bson.M{"_id": id, "medId": IDM, "resumeId": IDR,
						"sertificates": p.Args["sertificates"], "achievs": p.Args["achievs"]}, nil
				},
			},
			"createWork": &graphql.Field{
				Type:        WorkType,
				Description: "Create new Work",
				Args: graphql.FieldConfigArgument{
					"date": &graphql.ArgumentConfig{
						Type: graphql.DateTime,
					},
					"company": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"info": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"requirements": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
					},
					"type": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"phone": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, e := create(client.Database("HHCustom").Collection("Works"),
						bson.M{"date": p.Args["date"], "company": p.Args["company"],
							"info": p.Args["info"], "requirements": p.Args["requirements"],
							"type": p.Args["type"], "phone": p.Args["phone"], "email": p.Args["email"]})
					if e != nil {
						return nil, e
					}
					return bson.M{"_id": id, "date": p.Args["date"], "company": p.Args["company"],
						"info": p.Args["info"], "requirements": p.Args["requirements"],
						"type": p.Args["type"], "phone": p.Args["phone"], "email": p.Args["email"]}, nil
				},
			},
			"createUser": &graphql.Field{
				Type:        UserType,
				Description: "Create new Users",
				Args: graphql.FieldConfigArgument{
					"dob": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.DateTime),
					},
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"photo": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"phone": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"role": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"userInfoId": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
					"sesId": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					IDUI, _ := primitive.ObjectIDFromHex(p.Args["userInfoId"].(string))
					IDS, _ := primitive.ObjectIDFromHex(p.Args["sesId"].(string))
					id, e := create(client.Database("HHCustom").Collection("Users"),
						bson.M{"dob": p.Args["dob"], "username": p.Args["username"],
							"email": p.Args["email"], "photo": p.Args["photo"],
							"phone": p.Args["phone"], "role": p.Args["role"],
							"userInfoId": IDUI, "sesId": IDS})
					if e != nil {
						return nil, e
					}
					return bson.M{"_id": id, "dob": p.Args["dob"], "username": p.Args["username"],
						"email": p.Args["email"], "photo": p.Args["photo"],
						"phone": p.Args["phone"], "role": p.Args["role"],
						"userInfoId": IDUI, "sesId": IDS}, nil
				},
			},

			"updateMedCard": &graphql.Field{
				Type:        MedCardType,
				Description: "update MedCard",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"bloodGroup": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"ills": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"phobies": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ID, e := primitive.ObjectIDFromHex(p.Args["id"].(string))
					filter := bson.D{{Key: "_id", Value: ID}}
					if e != nil {
						return nil, e
					}
					data := bson.M{}
					if p.Args["bloodGroup"] != nil {
						data["bloodGroup"] = p.Args["bloodGroup"]
					}
					if p.Args["ills"] != nil {
						data["ills"] = p.Args["ills"]
					}
					if p.Args["phobies"] != nil {
						data["phobies"] = p.Args["phobies"]
					}
					send := bson.D{{Key: "$set", Value: data}}
					if len(data) != 0 {
						e = update(client.Database("HHCustom").Collection("MedCards"), filter, send)
						if e != nil {
							return nil, e
						}
						data["_id"] = ID
						return data, nil
					}
					return nil, nil
				},
			},
			"updateMsg": &graphql.Field{
				Type:        MsgType,
				Description: "update Msg",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"status": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Boolean),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ID, e := primitive.ObjectIDFromHex(p.Args["id"].(string))
					filter := bson.D{{Key: "_id", Value: ID}}
					if e != nil {
						return nil, e
					}
					data := bson.M{}
					if p.Args["status"] != nil {
						data["status"] = p.Args["status"]
					}
					send := bson.D{{Key: "$set", Value: data}}
					if len(data) != 0 {
						e = update(client.Database("HHCustom").Collection("Msgs"), filter, send)
						if e != nil {
							return nil, e
						}
						data["_id"] = ID
						return data, nil
					}
					return nil, nil
				},
			},
			"updateResume": &graphql.Field{
				Type:        ResumeType,
				Description: "update Resume",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"skills": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"whereWorks": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"aboutMe": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"date": &graphql.ArgumentConfig{
						Type: graphql.DateTime,
					},
					"link": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ID, e := primitive.ObjectIDFromHex(p.Args["id"].(string))
					filter := bson.D{{Key: "_id", Value: ID}}
					if e != nil {
						return nil, e
					}
					data := bson.M{}
					if p.Args["skills"] != nil {
						data["skills"] = p.Args["skills"]
					}
					if p.Args["whereWorks"] != nil {
						data["whereWorks"] = p.Args["whereWorks"]
					}
					if p.Args["aboutMe"] != nil {
						data["aboutMe"] = p.Args["aboutMe"]
					}
					if p.Args["date"] != nil {
						data["date"] = p.Args["date"]
					}
					if p.Args["link"] != nil {
						data["link"] = p.Args["link"]
					}
					send := bson.D{{Key: "$set", Value: data}}
					if len(data) != 0 {
						e = update(client.Database("HHCustom").Collection("Resumes"), filter, send)
						if e != nil {
							return nil, e
						}
						data["_id"] = ID
						return data, nil
					}
					return nil, nil
				},
			},
			"updateSession": &graphql.Field{
				Type:        SessionType,
				Description: "update Session",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"expire": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ID, e := primitive.ObjectIDFromHex(p.Args["id"].(string))
					filter := bson.D{{Key: "_id", Value: ID}}
					if e != nil {
						return nil, e
					}
					data := bson.M{}
					if p.Args["expire"] != nil {
						data["expire"] = p.Args["expire"]
					}
					send := bson.D{{Key: "$set", Value: data}}
					if len(data) != 0 {
						e = update(client.Database("HHCustom").Collection("Sessions"), filter, send)
						if e != nil {
							return nil, e
						}
						data["_id"] = ID
						return data, nil
					}
					return nil, nil
				},
			},
			"updateUserInfo": &graphql.Field{
				Type:        UserInfoType,
				Description: "update StudentInfo",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"medId": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
					"resumeId": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
					"sertificates": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"achievs": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ID, e := primitive.ObjectIDFromHex(p.Args["id"].(string))
					filter := bson.D{{Key: "_id", Value: ID}}
					if e != nil {
						return nil, e
					}
					data := bson.M{}
					if p.Args["medId"] != nil {
						data["medId"] = p.Args["medId"]
					}
					if p.Args["resumeId"] != nil {
						data["resumeId"] = p.Args["resumeId"]
					}
					if p.Args["sertificates"] != nil {
						data["sertificates"] = p.Args["sertificates"]
					}
					if p.Args["achievs"] != nil {
						data["achievs"] = p.Args["achievs"]
					}
					send := bson.D{{Key: "$set", Value: data}}
					if len(data) != 0 {
						e = update(client.Database("HHCustom").Collection("StudentInfos"), filter, send)
						if e != nil {
							return nil, e
						}
						data["_id"] = ID
						return data, nil
					}
					return nil, nil
				},
			},
			"updateUser": &graphql.Field{
				Type:        UserType,
				Description: "update Users",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"dob": &graphql.ArgumentConfig{
						Type: graphql.DateTime,
					},
					"username": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"photo": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"phone": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"role": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"userInfoId": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
					"sesId": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ID, e := primitive.ObjectIDFromHex(p.Args["id"].(string))
					filter := bson.D{{Key: "_id", Value: ID}}
					if e != nil {
						return nil, e
					}
					data := bson.M{}
					if p.Args["dob"] != nil {
						data["dob"] = p.Args["dob"]
					}
					if p.Args["username"] != nil {
						data["username"] = p.Args["username"]
					}
					if p.Args["email"] != nil {
						data["email"] = p.Args["email"]
					}
					if p.Args["photo"] != nil {
						data["photo"] = p.Args["photo"]
					}
					if p.Args["phone"] != nil {
						data["phone"] = p.Args["phone"]
					}
					if p.Args["role"] != nil {
						data["role"] = p.Args["role"]
					}
					if p.Args["userInfoId"] != nil {
						data["userInfoId"] = p.Args["userInfoId"]
					}
					if p.Args["sesId"] != nil {
						data["sesId"] = p.Args["sesId"]
					}
					send := bson.D{{Key: "$set", Value: data}}
					if len(data) != 0 {
						e = update(client.Database("HHCustom").Collection("Users"), filter, send)
						if e != nil {
							return nil, e
						}
						data["_id"] = ID
						return data, nil
					}
					return nil, nil
				},
			},

			"deleteMsg": &graphql.Field{
				Type:        MsgType,
				Description: "Delete messgage by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ID, e := primitive.ObjectIDFromHex(p.Args["id"].(string))
					filter := bson.D{{Key: "_id", Value: ID}}
					if e != nil {
						return nil, e
					}
					e = delete(client.Database("HHCustom").Collection("Msgs"), filter)
					if e != nil {
						return nil, e
					}
					return bson.M{"_id": ID}, nil
				},
			},
			"deleteSession": &graphql.Field{
				Type:        SessionType,
				Description: "Delete session by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ID, e := primitive.ObjectIDFromHex(p.Args["id"].(string))
					filter := bson.D{{Key: "_id", Value: ID}}
					if e != nil {
						return nil, e
					}
					e = delete(client.Database("HHCustom").Collection("Sessions"), filter)
					if e != nil {
						return nil, e
					}
					return bson.M{"_id": ID}, nil
				},
			},
			"deleteWork": &graphql.Field{
				Type:        WorkType,
				Description: "Delete work by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ID, e := primitive.ObjectIDFromHex(p.Args["id"].(string))
					filter := bson.D{{Key: "_id", Value: ID}}
					if e != nil {
						return nil, e
					}
					e = delete(client.Database("HHCustom").Collection("Works"), filter)
					if e != nil {
						return nil, e
					}
					return bson.M{"_id": ID}, nil
				},
			},
		},
	},
)

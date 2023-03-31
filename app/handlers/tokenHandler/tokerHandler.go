package tokenhandler

import (
	"app/database"
	"app/models/tokenModel"
	"app/utilities"
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createToken(c *fiber.Ctx) error {
	var token = tokenModel.TokenModel{}
	req := c.Body()
	json.Unmarshal(req, &token)
	coll := database.Instance.Db.Collection("tokens")
	_, err := coll.InsertOne(context.TODO(), token)
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	return c.Send(utilities.MsgJson(utilities.Success))
}

func getUsertokens(c *fiber.Ctx) error {

	var req map[string]interface{}
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		panic(err)
	}
	userId, err := primitive.ObjectIDFromHex(req["userid"].(string))
	filter := bson.D{{"userid", userId}}
	coll := database.Instance.Db.Collection("tokens")
	matchStage := bson.D{
		{"$match", filter},
	}

	subProjectStage := bson.D{
		{"$project", bson.D{
			{"name", 1},
			{"data", 1},
			{"email", 1},
			{"accound_status", bson.D{{"$not", "$account_deactivated"}}},
		}},
	}
	lookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "users"},
			{"pipeline", bson.A{subProjectStage}},
			{"localField", "userid"},
			{"foreignField", "_id"},
			{"as", "userdetails"},
		}}}
	unWindStage := bson.D{
		{"$unwind", "$userdetails"},
	}
	// projectStage := bson.D{{
	// 	"$project", bson.D{
	// 		{"adminid", 0},
	// 		{"userid", 0},
	// 		{"batchid", 0},
	// 		{"containerpassword", 0},
	// 	},
	// }}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, lookupStage, unWindStage})
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	var data []bson.M
	if err = cursor.All(context.TODO(), &data); err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	jsondata, err := json.Marshal(data)
	if data == nil {
		return c.Send(utilities.MsgJson(utilities.NoData))
	}
	return c.Send(jsondata)
}

func listTokens(c *fiber.Ctx) error {

	coll := database.Instance.Db.Collection("tokens")
	// matchStage := bson.D{
	// 	{"$match", filter},
	// }

	subProjectStage := bson.D{
		{"$project", bson.D{
			{"name", 1},
			{"data", 1},
			{"email", 1},
			{"accound_status", bson.D{{"$not", "$account_deactivated"}}},
		}},
	}
	lookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "users"},
			{"pipeline", bson.A{subProjectStage}},
			{"localField", "userid"},
			{"foreignField", "_id"},
			{"as", "userdetails"},
		}}}
	unWindStage := bson.D{
		{"$unwind", "$userdetails"},
	}
	// projectStage := bson.D{{
	// 	"$project", bson.D{
	// 		{"adminid", 0},
	// 		{"userid", 0},
	// 		{"batchid", 0},
	// 		{"containerpassword", 0},
	// 	},
	// }}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{lookupStage, unWindStage})
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	var data []bson.M
	if err = cursor.All(context.TODO(), &data); err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	jsondata, err := json.Marshal(data)
	if data == nil {
		return c.Send(utilities.MsgJson(utilities.NoData))
	}
	return c.Send(jsondata)
}
func Register(_route fiber.Router) {
	_route.Post("/create", createToken)
	_route.Post("/usertokens", getUsertokens)
	_route.Get("/list", listTokens)
}

package adminhandler

import (
	"app/database"
	"app/models/userModel"
	"app/utilities"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func signup(c *fiber.Ctx) error {
	var user = userModel.UserModel{}
	req := c.Body()
	json.Unmarshal(req, &user)
	coll := database.Instance.Db.Collection("users")
	_, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	return c.Send(utilities.MsgJson(utilities.Success))
}
func userList(c *fiber.Ctx) error {
	coll := database.Instance.Db.Collection("users")
	matchStage := bson.D{
		{"$match", bson.D{{"user_type", "user"}}}}
	projectStage := bson.D{{
		"$project", bson.D{
			{"name", 1},
			{"data", 1},
			{"email", 1},
			{"accound_status", bson.D{{"$not", "$account_deactivated"}}},
		},
	}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		fmt.Println(err)
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	var data []bson.M
	if err = cursor.All(context.TODO(), &data); err != nil {
		fmt.Println(err)
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	jsondata, err := json.Marshal(data)
	if data == nil {
		return c.Send(utilities.MsgJson(utilities.NoData))
	}
	return c.Send(jsondata)
}
func adminList(c *fiber.Ctx) error {
	coll := database.Instance.Db.Collection("users")
	matchStage := bson.D{
		{"$match", bson.D{{"user_type", "admin"}}}}
	projectStage := bson.D{{
		"$project", bson.D{
			{"name", 1},
			// {"data", 1},
			{"email", 1},
			{"accound_status", bson.D{{"$not", "$account_deactivated"}}},
		},
	}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		fmt.Println(err)
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	var data []bson.M
	if err = cursor.All(context.TODO(), &data); err != nil {
		fmt.Println(err)
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	jsondata, err := json.Marshal(data)
	if data == nil {
		return c.Send(utilities.MsgJson(utilities.NoData))
	}
	return c.Send(jsondata)
}
func Register(_route fiber.Router) {
	_route.Post("/create", signup)
	_route.Get("/user/list", userList)
	_route.Get("/list", adminList)
}

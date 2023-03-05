package imageHandler

import (
	"app/database"
	"app/models/imageModel"
	"app/utilities"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	active   = "active"
	deactive = "deactive"
)

func insertImage(c *fiber.Ctx) error {
	var image = imageModel.ImageModel{}
	req := c.Body()
	err := json.Unmarshal(req, &image)
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	image.ImageStatus = active
	// TODO add image status after pulling
	image.ImageId = "Tempid"
	image.ImageSize = "1G"
	coll := database.Instance.Db.Collection("images")
	_, err = coll.InsertOne(context.TODO(), image)
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	return c.Send(utilities.MsgJson(utilities.Success))
}
func list(c *fiber.Ctx) error {
	coll := database.Instance.Db.Collection("images")
	subProjectStage := bson.D{
		{"$project", bson.D{
			{"name", 1},
		}},
	}
	lookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "users"},
			{"pipeline", bson.A{subProjectStage}},
			{"localField", "adminId"},
			{"foreignField", "_id"},
			{"as", "adminUser"},
		}}}
	unWindStage := bson.D{
		{"$unwind", "$adminUser"},
	}
	projectStage := bson.D{{
		"$project", bson.D{
			{"adminId", 0},
			{"imagepull", 0},
		},
	}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{lookupStage, unWindStage, projectStage})
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
func get(c *fiber.Ctx) error {

	var req map[string]interface{}
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(req["id"])
	ImageId, err := primitive.ObjectIDFromHex(req["id"].(string))
	filter := bson.D{{"_id", ImageId}}
	coll := database.Instance.Db.Collection("images")
	matchStage := bson.D{
		{"$match", filter},
	}

	subProjectStage := bson.D{
		{"$project", bson.D{
			{"name", 1},
		}},
	}
	lookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "users"},
			{"pipeline", bson.A{subProjectStage}},
			{"localField", "adminId"},
			{"foreignField", "_id"},
			{"as", "adminUser"},
		}}}
	unWindStage := bson.D{
		{"$unwind", "$adminUser"},
	}
	projectStage := bson.D{{
		"$project", bson.D{
			{"adminId", 0},
			{"imagepull", 0},
		},
	}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, lookupStage, unWindStage, projectStage})
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	var data []bson.M
	if err = cursor.All(context.TODO(), &data); err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	jsondata, err := json.Marshal(data[0])
	if data == nil {
		return c.Send(utilities.MsgJson(utilities.NoData))
	}
	return c.Send(jsondata)
}
func deleteImage(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.D{{"_id", id}}
	coll := database.Instance.Db.Collection("images")
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	return c.Send(utilities.MsgJson(utilities.Success))
}
func Register(_route fiber.Router) {
	_route.Post("/create", insertImage)
	_route.Get("/list", list)
	_route.Post("/get", get)
	_route.Delete("/delete/:id", deleteImage)
}

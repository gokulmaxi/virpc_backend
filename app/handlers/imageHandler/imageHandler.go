package imageHandler

import (
	"app/database"
	"app/models/imageModel"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func insertImage(c *fiber.Ctx) error {
	var image = imageModel.ImageModel{}
	res := make(map[string]interface{})
	req := c.Body()
	err := json.Unmarshal(req, &image)
	if err != nil {
		panic(err)
	}
	fmt.Println(image)
	fmt.Println(image)
	coll := database.Instance.Db.Collection("images")
	_, err = coll.InsertOne(context.TODO(), image)
	if err != nil {
		res["message"] = "internal_error"
		data, _ := json.Marshal(res)
		return c.Send(data)
	}
	res["message"] = "done"
	data, _ := json.Marshal(res)
	return c.Send(data)
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
		},
	}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{lookupStage, unWindStage, projectStage})
	if err != nil {
		return err
	}
	var data []bson.M
	if err = cursor.All(context.TODO(), &data); err != nil {
		panic(err)
	}
	fmt.Println(data)
	jsondata, err := json.Marshal(data)
	return c.Send(jsondata)
}
func Register(_route fiber.Router) {
	_route.Post("/create", insertImage)
	_route.Get("/list", list)
}

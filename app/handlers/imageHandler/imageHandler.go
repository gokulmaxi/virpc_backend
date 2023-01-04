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
		panic(err)
	}
	return c.SendString("added image")
}
func list(c *fiber.Ctx) error {
	coll := database.Instance.Db.Collection("images")
	matchStage := bson.D{{"$match", bson.D{}}}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$imagename"},
			{"tags",
				bson.D{
					{"$addToSet", "$imagetag"},
				},
			},
		}}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage})
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

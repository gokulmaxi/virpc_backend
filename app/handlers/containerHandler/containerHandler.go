package containerHandler

import (
	"app/database"
	"app/models/containerModel"
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func insertContainer(c *fiber.Ctx) error {
	var container = containerModel.ContainerRequestModel{}
	req := c.Body()
	err := json.Unmarshal(req, &container)
	if err != nil {
		panic(err)
	}
	coll := database.Instance.Db.Collection("containers")
	_, err = coll.InsertOne(context.TODO(), container)
	if err != nil {
		return c.SendString("failed")
	}
	return c.SendString("Done")
}
func list(c *fiber.Ctx) error {
	results := []containerModel.ContainerRequestModel{}
	coll := database.Instance.Db.Collection("containers")
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return c.SendString("No containers found")
		}
		panic(err)
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	jsondata, err := json.Marshal(results)
	return c.Send(jsondata)
}
func Register(_route fiber.Router) {
	_route.Post("/create", insertContainer)
	_route.Get("/list", list)
}

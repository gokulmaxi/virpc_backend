package batchhandler

import (
	"app/database"
	"app/utilities"
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func listBatch(c *fiber.Ctx) error {

	coll := database.Instance.Db.Collection("batch")

	ImageSubProjectStage := bson.D{
		{"$project", bson.D{
			{"imagename", 1},
		}},
	}
	imageLookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "images"},
			{"pipeline", bson.A{ImageSubProjectStage}},
			{"localField", "imageid"},
			{"foreignField", "_id"},
			{"as", "imageData"},
		}}}
	imageUnwindStage := bson.D{
		{"$unwind", "$imageData"},
	}
	projectStage := bson.D{{
		"$project", bson.D{
			{"imageid", 0},
		},
	}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{imageLookupStage, imageUnwindStage, projectStage})
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	var data []bson.M
	if err = cursor.All(context.TODO(), &data); err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	jsondata, err := json.Marshal(data)
	return c.Send(jsondata)
}
func Register(_route fiber.Router) {
	_route.Get("/container/list", listBatch)
}

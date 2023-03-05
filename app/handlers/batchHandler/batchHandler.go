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
			{"as", "imagedata"},
		}}}
	imageUnwindStage := bson.D{
		{"$unwind", "$imagedata"},
	}
	adminSubProjectStage := bson.D{
		{"$project", bson.D{
			{"name", 1},
		}},
	}
	adminLookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "users"},
			{"pipeline", bson.A{adminSubProjectStage}},
			{"localField", "adminid"},
			{"foreignField", "_id"},
			{"as", "admindata"},
		}}}
	adminUnwindStage := bson.D{
		{"$unwind", "$admindata"},
	}
	projectStage := bson.D{{
		"$project", bson.D{
			{"imageid", 0},
		},
	}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{imageLookupStage, imageUnwindStage,
		adminLookupStage, adminUnwindStage, projectStage})
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
func list(c *fiber.Ctx) error {

	coll := database.Instance.Db.Collection("batch")

	projectStage := bson.D{{
		"$project", bson.D{
			{"batchname", 1},
			{"totaldays", 1},
			{"startdate", 1},
		},
	}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{projectStage})
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
	_route.Get("/list", list)
}

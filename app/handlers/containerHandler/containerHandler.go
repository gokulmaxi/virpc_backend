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
	res := make(map[string]interface{})
	err := json.Unmarshal(req, &container)
	if err != nil {
		panic(err)
	}
	// TODO create new user if not exist
	// TODO create new batch if not exist
	// TODO add container id while creating new containers
	coll := database.Instance.Db.Collection("containers")
	_, err = coll.InsertOne(context.TODO(), container)
	if err != nil {
		res["message"] = "internal_error"
		data, _ := json.Marshal(res)
		return c.Send(data)
	}
	res["message"] = "success"
	data, _ := json.Marshal(res)
	return c.Send(data)
}
func list(c *fiber.Ctx) error {
	coll := database.Instance.Db.Collection("containers")

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
			{"as", "adminUser"},
		}}}
	adminUnWindStage := bson.D{
		{"$unwind", "$adminUser"},
	}
	subProjectStage := bson.D{
		{"$project", bson.D{
			{"imagename", 1},
		}},
	}
	lookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "images"},
			{"pipeline", bson.A{subProjectStage}},
			{"localField", "imageId"},
			{"foreignField", "_id"},
			{"as", "image"},
		}}}
	unWindStage := bson.D{
		{"$unwind", "$image"},
	}
	projectStage := bson.D{{
		"$project", bson.D{
			{"adminid", 0},
			{"imageId", 0},
			{"userdetails", 0},
			{"add_features", 0},
		},
	}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{adminLookupStage, adminUnWindStage, lookupStage, unWindStage, projectStage})
	if err != nil {
		return err
	}
	var data []bson.M
	if err = cursor.All(context.TODO(), &data); err != nil {
		panic(err)
	}
	jsondata, err := json.Marshal(data)
	return c.Send(jsondata)
}
func Register(_route fiber.Router) {
	_route.Post("/create", insertContainer)
	_route.Get("/list", list)
}

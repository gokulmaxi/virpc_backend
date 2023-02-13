package containerHandler

import (
	"app/database"
	"app/models/batchModel"
	"app/models/containerModel"
	"app/models/userModel"
	"app/utilities"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func insertContainer(c *fiber.Ctx) error {
	var container = containerModel.ContainerRequestModel{}
	var req map[string]interface{}
	var results userModel.UserModel
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		panic(err)
	}
	//create new user if not exist
	userReq := req["userdetails"].(map[string]interface{})
	filter := bson.D{{"email", userReq["email"].(string)}}
	fmt.Println(filter)
	coll := database.Instance.Db.Collection("users")
	err = coll.FindOne(context.TODO(), filter).Decode(&results)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// create new user.
			var user = userModel.UserModel{}
			user.Email = userReq["email"].(string)
			user.Name = userReq["name"].(string)
			user.User_role = userModel.UserTypeNormal
			user.Data = userModel.NormalUser{}
			userId, err := coll.InsertOne(context.TODO(), user)
			if err != nil {
				return c.Send(utilities.MsgJson(utilities.Failure))
			}
			container.UserId = userId.InsertedID.(primitive.ObjectID)
		}
	} else {
		container.UserId = results.User_id
	}
	// TODO create new batch if not exist
	if req["batchid"] == nil {
		var batch = batchModel.BatchModel{}
		mapstructure.Decode(req["batch"], &batch)
		// TODO why obejct id is not parsed
		batchData := req["batch"].(map[string]interface{})
		batch.ImageId, err = primitive.ObjectIDFromHex(batchData["imageid"].(string))
		batchColl := database.Instance.Db.Collection("batch")
		col, err := batchColl.InsertOne(context.TODO(), batch)
		if err != nil {
			return c.Send(utilities.MsgJson(utilities.Failure))
		}
		container.BatchId = col.InsertedID.(primitive.ObjectID)
	} else {
		container.BatchId, err = primitive.ObjectIDFromHex(req["batchid"].(string))
	}
	// TODO add container id while creating new containers
	container.ContainerID = "asdf_asdf"
	container.ContainerPassword = req["containerpassword"].(string)
	collcontainer := database.Instance.Db.Collection("containers")
	_, err = collcontainer.InsertOne(context.TODO(), container)
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	return c.Send(utilities.MsgJson(utilities.Success))
}
func list(c *fiber.Ctx) error {
	coll := database.Instance.Db.Collection("containers")

	userSubProjectStage := bson.D{
		{"$project", bson.D{
			{"name", 1},
		}},
	}
	userLookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "users"},
			{"pipeline", bson.A{userSubProjectStage}},
			{"localField", "userid"},
			{"foreignField", "_id"},
			{"as", "userData"},
		}}}
	userUnWindStage := bson.D{
		{"$unwind", "$userData"},
	}

	batchSubLookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "images"},
			// {"pipeline", bson.A{userSubProjectStage}},
			{"localField", "imageid"},
			{"foreignField", "_id"},
			{"as", "imageData"},
		}}}
	batchLookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "batch"},
			{"pipeline", bson.A{batchSubLookupStage}},
			{"localField", "batchid"},
			{"foreignField", "_id"},
			{"as", "batchData"},
		}}}

	batchUnwindStage := bson.D{
		{"$unwind", "$batchData"},
	}
	imageUnwindStage := bson.D{
		{"$unwind", "$batchData.imageData"},
	}
	projectStage := bson.D{{
		"$project", bson.D{
			{"userData.name", 1},
			{"userData._id", 1},
			{"batchData.batchname", 1},
			{"batchData._id", 1},
			{"batchData.startdate", 1},
			{"batchData.imageData.imagename", 1},
			{"batchData.imageData._id", 1},
		},
	}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{userLookupStage, userUnWindStage,
		batchLookupStage, batchUnwindStage, imageUnwindStage, projectStage})
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
func listImage(c *fiber.Ctx) error {

	coll := database.Instance.Db.Collection("images")
	projectStage := bson.D{{
		"$project", bson.D{
			{"imagename", 1},
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
	if data == nil {
		return c.Send(utilities.MsgJson(utilities.NoData))
	}
	return c.Send(jsondata)
}
func Register(_route fiber.Router) {
	_route.Post("/create", insertContainer)
	_route.Get("/list", list)
	_route.Get("/imageList", listImage)
}

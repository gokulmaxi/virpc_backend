package containerHandler

import (
	"app/database"
	"app/models/batchModel"
	"app/models/containerModel"
	"app/models/imageModel"
	"app/models/userModel"
	"app/utilities"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/gofiber/fiber/v2"
	"github.com/goombaio/namegenerator"
	"github.com/mitchellh/mapstructure"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func startContainer(c *fiber.Ctx) error {
	var req map[string]interface{}
	var containerData containerModel.ContainerRequestModel
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		panic(err)
	}
	id, _ := primitive.ObjectIDFromHex(req["id"].(string))
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"status", containerModel.Running}}}}
	coll := database.Instance.Db.Collection("containers")
	err = coll.FindOne(context.TODO(), filter).Decode(&containerData)
	if err != nil {
		return c.Send(utilities.MsgJson(err.Error()))
	}
	if containerData.Status == containerModel.Running {
		return c.Send(utilities.MsgJson("Container is already running"))
	}
	if err := utilities.Docker.ContainerStart(context.TODO(), containerData.ContainerID, types.ContainerStartOptions{}); err != nil {
		return c.Send(utilities.MsgJson(err.Error()))
	}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return c.Send(utilities.MsgJson(utilities.Success))
}
func stopContainer(c *fiber.Ctx) error {
	var req map[string]interface{}
	var containerData containerModel.ContainerRequestModel
	timeout := time.Duration(10)
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		panic(err)
	}
	id, _ := primitive.ObjectIDFromHex(req["id"].(string))
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"status", containerModel.Stopped}}}}
	coll := database.Instance.Db.Collection("containers")
	err = coll.FindOne(context.TODO(), filter).Decode(&containerData)
	if err != nil {
		return c.Send(utilities.MsgJson(err.Error()))
	}
	if containerData.Status == containerModel.Stopped {
		return c.Send(utilities.MsgJson("Container is already stopped"))
	}
	if err := utilities.Docker.ContainerStop(context.TODO(), containerData.ContainerID, &timeout); err != nil {
		return c.Send(utilities.MsgJson(err.Error()))
	}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return c.Send(utilities.MsgJson(utilities.Success))
}
func insertContainer(c *fiber.Ctx) error {
	var containerData = containerModel.ContainerRequestModel{}
	var req map[string]interface{}
	var results userModel.UserModel
	var containerImageId primitive.ObjectID
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		panic(err)
	}
	containerData.ContainerPassword = req["containerpassword"].(string)
	containerData.AdminId, err = primitive.ObjectIDFromHex(req["adminId"].(string))
	collcontainer := database.Instance.Db.Collection("containers")
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
			containerData.UserId = userId.InsertedID.(primitive.ObjectID)
		}
	} else {
		containerData.UserId = results.User_id
	}
	// DONE create new batch if not exist
	if req["batchid"] == nil {
		var batch = batchModel.BatchModel{}
		mapstructure.Decode(req["batch"], &batch)
		// TODO why obejct id is not parsed
		batchData := req["batch"].(map[string]interface{})
		batch.ImageId, err = primitive.ObjectIDFromHex(batchData["imageid"].(string))
		containerImageId = batch.ImageId
		batch.AdminId, err = primitive.ObjectIDFromHex(req["adminId"].(string))
		batchColl := database.Instance.Db.Collection("batch")
		col, err := batchColl.InsertOne(context.TODO(), batch)
		if err != nil {
			return c.Send(utilities.MsgJson(utilities.Failure))
		}
		containerData.BatchId = col.InsertedID.(primitive.ObjectID)
	} else {
		var batchResult batchModel.BatchModel
		containerData.BatchId, err = primitive.ObjectIDFromHex(req["batchid"].(string))
		filter := bson.D{{"_id", containerData.BatchId}}
		coll := database.Instance.Db.Collection("batch")
		err = coll.FindOne(context.TODO(), filter).Decode(&batchResult)
		if err != nil {
			return c.Send(utilities.MsgJson(utilities.Failure))
		}
		containerImageId = batchResult.ImageId
	}
	// find image pull command
	var imageResult imageModel.ImageModel
	imageFilter := bson.D{{"_id", containerImageId}}
	imageColl := database.Instance.Db.Collection("images")
	err = imageColl.FindOne(context.TODO(), imageFilter).Decode(&imageResult)
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	fmt.Println("Creating container")
	// Find network id of backend
	// REVIEW will network id change every time if not get only once and use
	net, err := utilities.Docker.NetworkInspect(context.Background(), "vir-pc_backend", types.NetworkInspectOptions{})
	if err != nil {
		return c.Send(utilities.MsgJson(err.Error()))
	}
	containerConfig := &container.Config{Image: imageResult.ImagePull, Env: []string{"VNC_PW=" + containerData.ContainerPassword}, ExposedPorts: nat.PortSet{
		"6901/tcp": struct{}{},
	}}
	seed := time.Now().UTC().UnixNano()
	randomPort := strconv.Itoa(rand.Intn(2000) + 10000)
	containerData.ContainerPort = randomPort
	nameGenerator := namegenerator.NewNameGenerator(seed)
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"6901/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: randomPort,
				},
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
	}
	name := nameGenerator.Generate()
	//create container
	resp, err := utilities.Docker.ContainerCreate(context.TODO(), containerConfig, hostConfig, &network.NetworkingConfig{}, &v1.Platform{}, name)
	if err != nil {
		return c.Send(utilities.MsgJson(err.Error()))
	}
	// attach container to backend network
	networkConfig := &network.EndpointSettings{
		NetworkID: net.ID,
	}
	err = utilities.Docker.NetworkConnect(context.Background(), net.ID, resp.ID, networkConfig)
	if err != nil {
		return c.Send(utilities.MsgJson(err.Error()))
	}
	// Start the container
	if err := utilities.Docker.ContainerStart(context.TODO(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return c.Send(utilities.MsgJson(err.Error()))
	}
	containerData.ContainerID = resp.ID
	containerData.ContainerName = strings.Replace(name, "/", "", -1)
	containerData.Status = containerModel.Running
	dbContainerData, err := collcontainer.InsertOne(context.TODO(), containerData)
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	err = utilities.Redis.Set(context.Background(), dbContainerData.InsertedID.(primitive.ObjectID).Hex(), containerData.ContainerName, 0).Err()
	if err != nil {
		return c.Send(utilities.MsgJson(err.Error()))
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
			{"as", "adminData"},
		}}}
	adminUnWindStage := bson.D{
		{"$unwind", "$adminData"},
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
			{"adminData.name", 1},
			{"userData._id", 1},
			{"adminData._id", 1},
			{"batchData.batchname", 1},
			{"batchData._id", 1},
			{"batchData.startdate", 1},
			{"batchData.imageData.imagename", 1},
			{"batchData.imageData._id", 1},
			{"status", 1},
			{"containername", 1},
			{"containerid", 1},
		},
	}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{userLookupStage, userUnWindStage, adminLookupStage, adminUnWindStage,
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
func deleteContainer(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	var containerData containerModel.ContainerRequestModel
	timeout := time.Duration(10)
	filter := bson.D{{"_id", id}}
	coll := database.Instance.Db.Collection("containers")
	err = coll.FindOne(context.TODO(), filter).Decode(&containerData)
	if err != nil {
		return c.Send(utilities.MsgJson(err.Error()))
	}
	if containerData.Status == containerModel.Running {
		if err := utilities.Docker.ContainerStop(context.TODO(), containerData.ContainerID, &timeout); err != nil {
			return c.Send(utilities.MsgJson(err.Error()))
		}
	}
	if err := utilities.Docker.ContainerRemove(context.TODO(), containerData.ContainerID, types.ContainerRemoveOptions{}); err != nil {
		return c.Send(utilities.MsgJson(err.Error()))
	}
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return c.Send(utilities.MsgJson(utilities.Failure))
	}
	return c.Send(utilities.MsgJson(utilities.Success))
}
func Register(_route fiber.Router) {
	_route.Post("/create", insertContainer)
	_route.Get("/list", list)
	_route.Get("/imageList", listImage)
	_route.Post("/stop", stopContainer)
	_route.Post("/start", startContainer)
	_route.Delete("/delete/:id", deleteContainer)
}

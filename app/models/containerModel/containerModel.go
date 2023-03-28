package containerModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Running = "running"
	Stopped = "stopped"
)

type UserDetails struct {
	Email             string
	Name              string
	ContainerPassword string
}
type ContainerRequestModel struct {
	UserId            primitive.ObjectID
	BatchId           primitive.ObjectID
	AdminId           primitive.ObjectID
	ContainerImage    string
	ContainerPassword string
	ContainerID       string
	ContainerName     string
	Status            string
}

// {
// 	"batchname":"linux training - I",
// 	"batchdescription":"this is new batch for linux training one",
// 	"imageid":"39193492349",
// 	"startdate":"02/05/2020",
// 	"enddate":"02/05/2020",
// 	"totaldays":55,
// 	"cpulimit":5,
// 	"addfeatures" : ["internet_access", "root_access", "gpu_support"],
// 	"userdetails" :
// 	{
// 				"email" : "kishore.ct19@bits",
// 				"name"  : "Kishhh",
// 				"containerpassword": "asdasd"
// 	},
// 	"adminid" : "asdasdasdasd"
// }

package containerModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDetails struct {
	Email             string
	Name              string
	ContainerPassword string
}
type ContainerRequestModel struct {
	ImageId          primitive.ObjectID `bson:"imageId"`
	BatchName        string
	BatchDescription string
	StartDate        string
	EndDate          string
	CpuLimit         int
	Totaldays        int
	Add_features     []string
	UserDetails      UserDetails
	AdminId          primitive.ObjectID
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

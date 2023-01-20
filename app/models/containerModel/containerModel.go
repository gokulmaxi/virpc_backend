package containerModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContainerRequestModel struct {
	ImageId          primitive.ObjectID `bson:"imageId"`
	BatchName        string
	BatchDescription string
	StartDate        string
	EndDate          string
	CpuLimit         int
	Add_features     []string
}

// {
// 	"image_name": "python",
// 	"image_tag": "minimal",
// 	"image_visibility": true,
// 	"image_description": "ubuntu image",
// 	"expose_ports": 5000,
// 	"require_gpu": false,
// 	"cpu_limit": 2,
// 	"approval_status": false,
// 	"env_values": [
// 		{
// 			"HOME_DIR": "/root/home"
// 		}
// 	]
// }

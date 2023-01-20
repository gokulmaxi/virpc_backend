package containerModel

type containerModel struct {
	ImageName        string `bson:"imagename"`
	BaseImage        string
	ImageVersion     string
	ImagePull        string
	ImageDescription string
	RequireGpu       bool
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

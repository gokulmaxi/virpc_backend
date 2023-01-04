package imageModel

type ImageModel struct {
	ImageName        string `bson:"imagename"`
	ImageTag         string
	ImageVisibiltity string
	ImageDescription string
	ExposePorts      int
	RequireGpu       bool
	CpuLimit         int
	AprovalStatus    bool
	Env              map[string]string
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

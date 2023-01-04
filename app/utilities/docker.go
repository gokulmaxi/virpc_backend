package utilities

import (
	"github.com/docker/docker/client"
)

var Docker client.Client

func InitDocker() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	Docker = *cli
}

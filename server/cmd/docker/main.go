package main

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func main() {
	// client.NewEnvClient is deprecated: use [NewClientWithOpts] passing the [FromEnv] option.
	c, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion("1.41"), // 报错指定版本号，error: Error response from daemon: client version 1.41 is too new. Maximum supported API version is 1.39
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	resp, err := c.ContainerCreate(
		ctx,
		&container.Config{
			Image:        "mongo:latest",
			ExposedPorts: nat.PortSet{"27017/tcp": {}},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"27017/tcp": []nat.PortBinding{
					{HostIP: "127.0.0.1", HostPort: "0"}, // 0自动选择端口
				},
			},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		panic(err)
	}

	err = c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("container started")
	time.Sleep(5 * time.Second)

	inspRes, err := c.ContainerInspect(ctx, resp.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("listening at %+v\n", inspRes.NetworkSettings.Ports["27017/tcp"][0])

	fmt.Println("killing container")
	err = c.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
		Force: true, // 强制删除
	})
	if err != nil {
		panic(err)
	}
}

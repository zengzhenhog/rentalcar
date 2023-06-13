package mongotesting

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	image           = "mongo"
	containerPort   = "27017/tcp"
	defaultMongoURI = "mongodb://admin:123456@localhost:27017"
)

var mongoURI string

func RunWithMongoInDocker(m *testing.M) int {
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
			Image:        image,
			ExposedPorts: nat.PortSet{containerPort: {}},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				containerPort: []nat.PortBinding{
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
	containerID := resp.ID
	defer func() {
		err = c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force: true, // 强制删除
		})
		if err != nil {
			panic(err)
		}
	}()

	err = c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	inspRes, err := c.ContainerInspect(ctx, containerID)
	if err != nil {
		panic(err)
	}
	hostPort := inspRes.NetworkSettings.Ports[containerPort][0]
	mongoURI = fmt.Sprintf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)

	return m.Run()
}

func NewClient(c context.Context) (*mongo.Client, error) {
	if mongoURI == "" {
		return nil, fmt.Errorf("mongo uri not set. Please run RunWithMongoInDocker in TestMain")
	}
	return mongo.Connect(c, options.Client().ApplyURI(mongoURI))
}

func NewDefaultClient(c context.Context) (*mongo.Client, error) {
	return mongo.Connect(c, options.Client().ApplyURI(defaultMongoURI))
}

func SetupIndexes(c context.Context, d *mongo.Database) error {
	_, err := d.Collection("account").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{
			{Key: "open_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return err
	}

	_, err = d.Collection("trip").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{ // bson.D有序
			{Key: "trip.accountid", Value: 1},
			{Key: "trip.status", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{ // bson.M无序
			"trip.status": 1,
		}),
	})

	return err
}

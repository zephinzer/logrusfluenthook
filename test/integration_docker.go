package integrationtest

import (
	"bytes"
	"context"
	"log"
	"io"
	"strings"
	"time"

	dockerClient "github.com/docker/docker/client"
	dockerTypes "github.com/docker/docker/api/types"
)

func createDockerClient() *dockerClient.Client {
	client, err := dockerClient.NewEnvClient()
	if err != nil {
		panic(err)
	}
	return client
}

// getContainerLogs retrieves the logs of a container as a string
func getContainerLogs(
	docker *dockerClient.Client,
	containerID string,
) string {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	var containerLogs bytes.Buffer
	logsReader, err := docker.ContainerLogs(
		ctx,
		containerID,
		dockerTypes.ContainerLogsOptions{ShowStdout: true},
	)
	if err != nil {
		panic(err)
	} else if _, err = io.Copy(&containerLogs, logsReader); err != nil {
		panic(err)
	}

	return strings.Trim(containerLogs.String(), "\n")
}

func pullImage(
	docker *dockerClient.Client,
	imageUrl string,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()
	pullLogs, err := docker.ImagePull(
		ctx,
		imageUrl,
		dockerTypes.ImagePullOptions{},
	)
	if err != nil {
		panic(err)
	} else {
		log.Println(pullLogs)
	}
}

func startContainer(
	docker *dockerClient.Client,
	containerID string,
) {
	if err := docker.ContainerStart(
		context.Background(),
		containerID,
		dockerTypes.ContainerStartOptions{},
	); err != nil {
		panic(err)
	}
}

func stopContainer(
	docker *dockerClient.Client,
	containerID string,
) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()
	err := docker.ContainerStop(ctx, containerID, &dockerOpsTimeout)
	if err != nil {
		panic(err)
	}
}
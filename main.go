package main

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"golang.org/x/net/context"
)

func main() {
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, nil)
	if err != nil {
		panic(err)

	}

	ctx := context.Background()

	_, err = cli.ImagePull(ctx, "alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)

	}

	containerConfig := &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID); err != nil {
		panic(err)
	}

	statusCode, err := cli.ContainerWait(ctx, resp.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println(statusCode)

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}

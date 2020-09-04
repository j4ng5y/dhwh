package dkr

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RestartContainer(ctx context.Context, containerName string) error {
	to := time.Duration(30) * time.Second

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	if err := cli.ContainerRestart(ctx, containerName, &to); err != nil {
		return err
	}

	return nil
}

func PullNewContainer(ctx context.Context, name string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	out, err := cli.ImagePull(ctx, name, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	defer out.Close()
	return nil
}

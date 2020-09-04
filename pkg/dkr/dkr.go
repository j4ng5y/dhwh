// dkr contains all of the Docker specific logic

package dkr

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// RestartContainer is a function that will programmatically restart a docker
// container based on the name of the running container.
//
// Arguments:
//     ctx (context.Context):  The context in which to perform the action.
//     containerName (string): The name of the running container to restart.
//
// Returns:
//     (error): An error if one exists, nil otherwise.
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

// PullNewContainerImage is a function that will programmatically pull a
// new version of a container image if one exists.
//
// Arguments:
//     ctx (context.Context):  The context in which to perform the action.
//     name (string):          The name of the container image to pull.
//
// Returns:
//     (error): An error if one exists, nil otherwise.
func PullNewContainerImage(ctx context.Context, name string) error {
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

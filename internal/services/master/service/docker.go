package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/boreas/internal/pkg/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type DockerService struct {
	client   *client.Client
	registry string
}

func NewDockerService(registry string) (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &DockerService{
		client:   cli,
		registry: registry,
	}, nil
}

func (d *DockerService) BuildImage(ctx context.Context, appPath, appName, tag string, buildConfig *models.BuildConfig) (string, error) {
	if buildConfig == nil {
		buildConfig = &models.BuildConfig{
			Dockerfile: "Dockerfile",
			Context:    ".",
		}
	}

	contextPath := filepath.Join(appPath, buildConfig.Context)
	dockerfilePath := filepath.Join(appPath, buildConfig.Dockerfile)

	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("dockerfile not found: %s", dockerfilePath)
	}

	tar, err := archive.TarWithOptions(contextPath, &archive.TarOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create tar: %w", err)
	}
	defer tar.Close()

	imageName := fmt.Sprintf("%s/%s:%s", d.registry, appName, tag)
	buildOptions := types.ImageBuildOptions{
		Dockerfile: buildConfig.Dockerfile,
		Tags:       []string{imageName},
		Remove:     true,
		BuildArgs:  buildConfig.BuildArgs,
	}

	resp, err := d.client.ImageBuild(ctx, tar, buildOptions)
	if err != nil {
		return "", fmt.Errorf("failed to build image: %w", err)
	}
	defer resp.Body.Close()

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read build output: %w", err)
	}

	log.Printf("Build output: %s", string(output))

	return imageName, nil
}

func (d *DockerService) PushImage(ctx context.Context, imageName string) error {
	authConfig := registry.AuthConfig{}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return fmt.Errorf("failed to encode auth config: %w", err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	resp, err := d.client.ImagePush(ctx, imageName, image.PushOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}
	defer resp.Close()

	output, err := io.ReadAll(resp)
	if err != nil {
		return fmt.Errorf("failed to read push output: %w", err)
	}

	log.Printf("Push output: %s", string(output))

	return nil
}

func (d *DockerService) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}

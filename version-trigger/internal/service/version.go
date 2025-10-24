package service

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
)

type Config struct {
	WebhookSecret  string
	WorkDir        string
	DockerRegistry string
	ManagementAPI  string
}

type VersionService struct {
	config     *Config
	git        *GitService
	docker     *DockerService
	management *ManagementClient
}

func NewVersionService(config *Config) *VersionService {
	dockerService, err := NewDockerService(config.DockerRegistry)
	if err != nil {
		log.Printf("Warning: Failed to create docker service: %v", err)
	}

	return &VersionService{
		config:     config,
		git:        NewGitService(config.WorkDir),
		docker:     dockerService,
		management: NewManagementClient(config.ManagementAPI),
	}
}

type TagEvent struct {
	TagName    string
	Repository string
	Commit     string
	Pusher     string
}

type ProcessResult struct {
	Message        string     `json:"message"`
	VersionCreated bool       `json:"version_created"`
	VersionID      string     `json:"version_id,omitempty"`
	AppsBuilt      []AppBuild `json:"apps_built,omitempty"`
	Errors         []string   `json:"errors,omitempty"`
}

func (v *VersionService) ProcessTagEvent(ctx context.Context, event *TagEvent) (*ProcessResult, error) {
	result := &ProcessResult{
		Message: "Processing tag event",
	}

	log.Printf("Step 1: Cloning repository and checking out tag %s", event.TagName)
	repoPath, err := v.git.CloneOrPull(ctx, event.Repository, event.TagName)
	if err != nil {
		return nil, fmt.Errorf("failed to clone/pull repository: %w", err)
	}

	log.Printf("Step 2: Getting previous tag")
	previousTag, err := v.git.GetPreviousTag(ctx, repoPath, event.TagName)
	if err != nil {
		log.Printf("Warning: Failed to get previous tag: %v", err)
		previousTag = ""
	}

	log.Printf("Step 3: Getting changed apps (diff from %s to %s)", previousTag, event.TagName)
	changedApps, err := v.git.GetChangedApps(ctx, repoPath, previousTag, event.TagName)
	if err != nil {
		return nil, fmt.Errorf("failed to get changed apps: %w", err)
	}

	if len(changedApps) == 0 {
		result.Message = "No apps changed, skipping build"
		return result, nil
	}

	log.Printf("Found %d changed apps: %v", len(changedApps), changedApps)

	var builtApps []AppBuild
	var buildErrors []string

	for _, appName := range changedApps {
		log.Printf("Step 4.%d: Processing app %s", len(builtApps)+1, appName)

		app, err := v.management.GetApplication(ctx, appName)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to get app config for %s: %v", appName, err)
			log.Printf("Warning: %s", errMsg)
			buildErrors = append(buildErrors, errMsg)
			continue
		}

		if v.docker == nil {
			errMsg := fmt.Sprintf("Docker service not available for app %s", appName)
			log.Printf("Warning: %s", errMsg)
			buildErrors = append(buildErrors, errMsg)
			continue
		}

		appPath := filepath.Join(repoPath, appName)
		log.Printf("Building docker image for %s at %s", appName, appPath)

		imageName, err := v.docker.BuildImage(ctx, appPath, appName, event.TagName, app.Config.BuildConfig)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to build image for %s: %v", appName, err)
			log.Printf("Error: %s", errMsg)
			buildErrors = append(buildErrors, errMsg)
			continue
		}

		log.Printf("Pushing image %s", imageName)
		err = v.docker.PushImage(ctx, imageName)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to push image for %s: %v", appName, err)
			log.Printf("Error: %s", errMsg)
			buildErrors = append(buildErrors, errMsg)
			continue
		}

		builtApps = append(builtApps, AppBuild{
			AppID:       app.ID,
			AppName:     appName,
			DockerImage: imageName,
		})

		log.Printf("Successfully built and pushed %s", imageName)
	}

	if len(buildErrors) > 0 {
		result.Errors = buildErrors
	}

	if len(builtApps) == 0 {
		return result, fmt.Errorf("no apps were successfully built")
	}

	log.Printf("Step 5: Creating version in management system")
	versionReq := &CreateVersionRequest{
		Name:       event.TagName,
		GitTag:     event.TagName,
		GitCommit:  event.Commit,
		Repository: event.Repository,
		Metadata: map[string]string{
			"pusher": event.Pusher,
		},
		Apps: builtApps,
	}

	version, err := v.management.CreateVersion(ctx, versionReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	result.VersionCreated = true
	result.VersionID = version.ID
	result.AppsBuilt = builtApps
	result.Message = fmt.Sprintf("Successfully created version %s with %d apps", version.ID, len(builtApps))

	log.Printf("Process completed: %s", result.Message)

	return result, nil
}

func (v *VersionService) Close() error {
	if v.docker != nil {
		return v.docker.Close()
	}
	return nil
}

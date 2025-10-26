package service

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
)

type TriggerConfig struct {
	WebhookSecret  string
	WorkDir        string
	DockerRegistry string

	Apps    interfaces.ApplicationService
	Version interfaces.VersionService
}

type TriggerService struct {
	config *TriggerConfig
	git    *GitService
	docker *DockerService

	interfaces.ApplicationService
	interfaces.VersionService
}

func NewTriggerService(config *TriggerConfig) *TriggerService {
	dockerService, err := NewDockerService(config.DockerRegistry)
	if err != nil {
		log.Printf("Warning: Failed to create docker service: %v", err)
	}

	return &TriggerService{
		config:             config,
		git:                NewGitService(config.WorkDir),
		docker:             dockerService,
		ApplicationService: config.Apps,
		VersionService:     config.Version,
	}
}

type TagEvent struct {
	TagName    string
	Repository string
	Commit     string
	Pusher     string
}

type ProcessResult struct {
	Message        string            `json:"message"`
	TriggerCreated bool              `json:"Trigger_created"`
	TriggerID      string            `json:"Trigger_id,omitempty"`
	AppsBuilt      []models.AppBuild `json:"apps_built,omitempty"`
	Errors         []string          `json:"errors,omitempty"`
}

func (v *TriggerService) ProcessTagEvent(ctx context.Context, event *TagEvent) (*ProcessResult, error) {
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

	var builtApps []models.AppBuild
	var buildErrors []string

	for _, appName := range changedApps {
		log.Printf("Step 4.%d: Processing app %s", len(builtApps)+1, appName)

		app, err := v.GetApplication(ctx, appName)
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

		imageName, err := v.docker.BuildImage(ctx, appPath, appName, event.TagName, app.GetBuildConfig())
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

		builtApps = append(builtApps, models.AppBuild{
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

	log.Printf("Step 5: Creating Trigger in management system")
	TriggerReq := &models.CreateVersionRequest{
		GitTag:     event.TagName,
		GitCommit:  event.Commit,
		Repository: event.Repository,
		AppBuilds:  builtApps,
	}

	Trigger, err := v.CreateVersion(ctx, TriggerReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create Trigger: %w", err)
	}

	result.TriggerCreated = true
	result.TriggerID = Trigger.ID
	result.AppsBuilt = builtApps
	result.Message = fmt.Sprintf("Successfully created Trigger %s with %d apps", Trigger.ID, len(builtApps))

	log.Printf("Process completed: %s", result.Message)

	return result, nil
}

func (v *TriggerService) Close() error {
	if v.docker != nil {
		return v.docker.Close()
	}
	return nil
}

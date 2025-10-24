package service

import (
	"github.com/XR0118/1024Hackathon/backend/model"
)

type WebhookService interface {
	HandleGitEvent(event *GitEvent) (*WebhookResponse, error)
	TriggerDeployment(versionID string, environments []string) ([]*model.Deployment, error)
}

type GitEvent struct {
	Event         string `json:"event"`
	Ref           string `json:"ref"`
	Repository    string `json:"repository"`
	Commit        string `json:"commit"`
	Author        string `json:"author"`
	CommitMessage string `json:"commit_message"`
	Tag           string `json:"tag"`
}

type WebhookResponse struct {
	VersionID   string   `json:"version_id"`
	Deployments []string `json:"deployments"`
}

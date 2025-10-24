package service

import (
	"github.com/XR0118/1024Hackathon/backend/model"
)

type VersionService interface {
	Create(version *model.Version) error
	GetByID(id string) (*model.Version, error)
	List(repository string, status string, page, size int) ([]*model.Version, int, error)
	UpdateStatus(id string, status string) error
	Compare(fromID, toID string) (*VersionComparison, error)
}

type VersionComparison struct {
	From *model.Version `json:"from"`
	To   *model.Version `json:"to"`
	Diff *VersionDiff   `json:"diff"`
}

type VersionDiff struct {
	Commits []string `json:"commits"`
	Changes []string `json:"changes"`
}

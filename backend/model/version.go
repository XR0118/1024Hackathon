package model

import "time"

type Version struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	GitTag        string            `json:"git_tag"`
	GitCommit     string            `json:"git_commit"`
	Repository    string            `json:"repository"`
	IsRevert      bool              `json:"is_revert"`
	ParentVersion string            `json:"parent_version,omitempty"`
	CreatedBy     string            `json:"created_by"`
	CreatedAt     time.Time         `json:"created_at"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	Status        string            `json:"status"`
}

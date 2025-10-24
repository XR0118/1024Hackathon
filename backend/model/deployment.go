package model

import "time"

type Deployment struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	VersionID    string         `json:"version_id"`
	Applications []string       `json:"applications"`
	TargetEnvs   []string       `json:"target_envs"`
	Strategy     DeployStrategy `json:"strategy"`
	Status       string         `json:"status"`
	Progress     Progress       `json:"progress"`
	Approvals    []Approval     `json:"approvals"`
	CreatedBy    string         `json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	StartedAt    *time.Time     `json:"started_at,omitempty"`
	CompletedAt  *time.Time     `json:"completed_at,omitempty"`
}

type DeployStrategy struct {
	Type           string  `json:"type"`
	BatchSize      int     `json:"batch_size"`
	BatchInterval  int     `json:"batch_interval"`
	CanaryRatio    float64 `json:"canary_ratio"`
	AutoRollback   bool    `json:"auto_rollback"`
	ManualApproval bool    `json:"manual_approval"`
}

type Progress struct {
	Total        int          `json:"total"`
	Completed    int          `json:"completed"`
	Failed       int          `json:"failed"`
	CurrentBatch int          `json:"current_batch"`
	Details      []TaskDetail `json:"details"`
}

type TaskDetail struct {
	AppID       string     `json:"app_id"`
	EnvID       string     `json:"env_id"`
	Status      string     `json:"status"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	ErrorMsg    string     `json:"error_msg,omitempty"`
}

type Approval struct {
	ApproverID string    `json:"approver_id"`
	Action     string    `json:"action"`
	Comment    string    `json:"comment"`
	Timestamp  time.Time `json:"timestamp"`
}

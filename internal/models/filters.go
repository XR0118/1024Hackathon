package models

// VersionFilter 版本过滤器
type VersionFilter struct {
	Repository string
	Page       int
	PageSize   int
}

// ApplicationFilter 应用过滤器
type ApplicationFilter struct {
	Repository string
	Type       string
	Page       int
	PageSize   int
}

// EnvironmentFilter 环境过滤器
type EnvironmentFilter struct {
	Type     string
	IsActive *bool
	Page     int
	PageSize int
}

// DeploymentFilter 部署过滤器
type DeploymentFilter struct {
	Status        DeploymentStatus
	EnvironmentID string
	VersionID     string
	Page          int
	PageSize      int
}

// TaskFilter 任务过滤器
type TaskFilter struct {
	DeploymentID string
	Status       TaskStatus
	Type         string
	Page         int
	PageSize     int
}

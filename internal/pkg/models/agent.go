package models

import (
	"time"
)

// AgentAppStatus 应用在物理机上的状态
type AgentAppStatus struct {
	App      string            `json:"app"`
	Version  string            `json:"version"`
	Replicas int               `json:"replicas"`
	Healthy  HealthStatus      `json:"healthy"`
	Config   map[string]string `json:"config,omitempty"`
	Status   string            `json:"status"` // running, stopped, error
	Updated  time.Time         `json:"updated"`
}

// HealthStatus 健康状态
type HealthStatus struct {
	Level int    `json:"level"` // 0-100, 0表示不健康，100表示完全健康
	Msg   string `json:"msg,omitempty"`
}

// ApplyRequest 应用部署请求
type ApplyRequest struct {
	App     string                 `json:"app" binding:"required"`
	Version string                 `json:"version" binding:"required"`
	Pkg     map[string]interface{} `json:"pkg" binding:"required"`
}

// ApplyResponse 应用部署响应
type ApplyResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	App     string `json:"app"`
	Version string `json:"version"`
}

// StatusResponse 状态查询响应
type StatusResponse struct {
	Apps []AgentAppStatus `json:"apps"`
}

// AppStatusResponse 单个应用状态响应
type AppStatusResponse struct {
	App     string       `json:"app"`
	Version string       `json:"version"`
	Healthy HealthStatus `json:"healthy"`
}

// AgentConfig Agent配置
type AgentConfig struct {
	ID       string            `json:"id"`
	Hostname string            `json:"hostname"`
	IP       string            `json:"ip"`
	Config   map[string]string `json:"config"`
	Status   string            `json:"status"` // active, inactive
	Updated  time.Time         `json:"updated"`
}

// DeploymentPackage 部署包信息
type DeploymentPackage struct {
	Type        string            `json:"type"` // docker, binary, script
	Replicas    int               `json:"replicas,omitempty"`
	Image       string            `json:"image,omitempty"`
	Command     []string          `json:"command,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Volumes     []VolumeMount     `json:"volumes,omitempty"`
	Ports       []PortMapping     `json:"ports,omitempty"`
	Resources   ResourceLimits    `json:"resources,omitempty"`
}

// VolumeMount 卷挂载
type VolumeMount struct {
	HostPath      string `json:"host_path"`
	ContainerPath string `json:"container_path"`
	ReadOnly      bool   `json:"read_only"`
}

// PortMapping 端口映射
type PortMapping struct {
	HostPort      int    `json:"host_port"`
	ContainerPort int    `json:"container_port"`
	Protocol      string `json:"protocol"` // tcp, udp
}

// ResourceLimits 资源限制
type ResourceLimits struct {
	CPULimit    string `json:"cpu_limit,omitempty"`
	MemoryLimit string `json:"memory_limit,omitempty"`
	CPUSet      string `json:"cpu_set,omitempty"`
}

// ContainerInfo 容器信息
type ContainerInfo struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	Status  string            `json:"status"`
	Created time.Time         `json:"created"`
	Config  map[string]string `json:"config,omitempty"`
}

// ProcessInfo 进程信息
type ProcessInfo struct {
	PID     int               `json:"pid"`
	Name    string            `json:"name"`
	Status  string            `json:"status"`
	Created time.Time         `json:"created"`
	Config  map[string]string `json:"config,omitempty"`
}

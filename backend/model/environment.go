package model

import "time"

type TargetEnvironment struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Region    string    `json:"region"`
	Config    EnvConfig `json:"config"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type EnvConfig struct {
	K8SConfig      *K8SConfig      `json:"k8s_config,omitempty"`
	PhysicalConfig *PhysicalConfig `json:"physical_config,omitempty"`
}

type K8SConfig struct {
	KubeConfig  string `json:"kube_config"`
	Namespace   string `json:"namespace"`
	ClusterName string `json:"cluster_name"`
}

type PhysicalConfig struct {
	Hosts     []Host    `json:"hosts"`
	SSHConfig SSHConfig `json:"ssh_config"`
}

type Host struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Role     string `json:"role"`
}

type SSHConfig struct {
	User    string `json:"user"`
	Port    int    `json:"port"`
	KeyPath string `json:"key_path"`
}

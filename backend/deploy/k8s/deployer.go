package k8s

import (
	"github.com/XR0118/1024Hackathon/backend/model"
)

type K8SDeployer interface {
	Deploy(app *model.Application, version *model.Version, env *model.TargetEnvironment) error
	CreateDeployment(app *model.Application, imageName string, env *model.TargetEnvironment) error
	UpdateDeployment(app *model.Application, imageName string, env *model.TargetEnvironment) error
	GetDeploymentStatus(app *model.Application, env *model.TargetEnvironment) (*DeploymentStatus, error)
	HealthCheck(app *model.Application, env *model.TargetEnvironment) bool
	Rollback(app *model.Application, env *model.TargetEnvironment, targetVersion string) error
}

type K8SClient interface {
	Connect(config *model.K8SConfig) error
	CreateDeployment(namespace, name string, spec *DeploymentSpec) error
	UpdateDeployment(namespace, name string, spec *DeploymentSpec) error
	GetDeployment(namespace, name string) (*K8SDeployment, error)
	GetPods(namespace, appName string) ([]*Pod, error)
	DeleteDeployment(namespace, name string) error
	ScaleDeployment(namespace, name string, replicas int) error
}

type DeploymentSpec struct {
	Name      string
	Namespace string
	Replicas  int32
	Image     string
	Ports     []ContainerPort
	Env       map[string]string
	Resources *model.Resources
}

type K8SDeployment struct {
	Name      string
	Namespace string
	Replicas  int32
	Ready     int32
	Image     string
	Status    string
}

type Pod struct {
	Name      string
	IP        string
	Status    string
	Ready     bool
	Namespace string
}

type ContainerPort struct {
	Name          string
	ContainerPort int32
	Protocol      string
}

type DeploymentStatus struct {
	Phase             string
	AvailableReplicas int32
	ReadyReplicas     int32
	UpdatedReplicas   int32
	Conditions        []DeploymentCondition
}

type DeploymentCondition struct {
	Type    string
	Status  string
	Reason  string
	Message string
}

type ResourceManager interface {
	CreateService(namespace, name string, spec *ServiceSpec) error
	UpdateService(namespace, name string, spec *ServiceSpec) error
	CreateConfigMap(namespace, name string, data map[string]string) error
	UpdateConfigMap(namespace, name string, data map[string]string) error
}

type ServiceSpec struct {
	Name      string
	Namespace string
	Selector  map[string]string
	Ports     []ServicePort
	Type      string
}

type ServicePort struct {
	Name       string
	Port       int32
	TargetPort int32
	Protocol   string
}

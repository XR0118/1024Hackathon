package service

import (
	"context"
	"fmt"
	"time"

	"github.com/boreas/internal/pkg/models"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type OperatorK8sService struct {
	clientset *kubernetes.Clientset
	namespace string
	timeout   time.Duration
}

func NewOperatorK8sService(kubeconfig, namespace string, timeout int) (*OperatorK8sService, error) {
	var config *rest.Config
	var err error

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to build kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	if namespace == "" {
		namespace = "default"
	}

	if timeout <= 0 {
		timeout = 30
	}

	return &OperatorK8sService{
		clientset: clientset,
		namespace: namespace,
		timeout:   time.Duration(timeout) * time.Second,
	}, nil
}

func (s *OperatorK8sService) CheckK8sConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	_, err := s.clientset.CoreV1().Namespaces().Get(ctx, s.namespace, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to connect to kubernetes: %w", err)
	}

	return nil
}

func (s *OperatorK8sService) Apply(req *models.ApplyRequest) (*models.ApplyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	pkg := req.Package.Type
	if pkg == "" {
		return nil, fmt.Errorf("package type is required")
	}

	switch pkg {
	case "docker":
		return s.applyDockerDeployment(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported package type: %s", pkg)
	}
}

func (s *OperatorK8sService) GetStatus(app string) (*models.AppStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	deployment, err := s.clientset.AppsV1().Deployments(s.namespace).Get(ctx, app, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	// 计算健康度：加权平均（假设每个就绪Pod健康度为100，未就绪为0）
	// 健康度 = (就绪副本数 / 总副本数) * 100
	healthy := models.HealthInfo{
		Level: 0,
		Msg:   "No replicas ready",
	}

	if deployment.Status.Replicas > 0 {
		if deployment.Status.ReadyReplicas == deployment.Status.Replicas {
			healthy.Level = 100
			healthy.Msg = "All replicas ready"
		} else if deployment.Status.ReadyReplicas > 0 {
			healthy.Level = int(float64(deployment.Status.ReadyReplicas) / float64(deployment.Status.Replicas) * 100)
			healthy.Msg = fmt.Sprintf("%d/%d replicas ready", deployment.Status.ReadyReplicas, deployment.Status.Replicas)
		} else {
			healthy.Msg = "No replicas ready"
		}
	} else {
		healthy.Msg = "No replicas configured"
	}

	version := ""
	if deployment.Spec.Template.Spec.Containers != nil && len(deployment.Spec.Template.Spec.Containers) > 0 {
		version = deployment.Spec.Template.Spec.Containers[0].Image
	}

	return &models.AppStatusResponse{
		App:     app,
		Version: version,
		Healthy: healthy,
	}, nil
}

func (s *OperatorK8sService) applyDockerDeployment(ctx context.Context, req *models.ApplyRequest) (*models.ApplyResponse, error) {
	image := req.Package.Image
	if image == "" {
		return nil, fmt.Errorf("image is required for docker deployment")
	}

	replicas := int32(1)
	if req.Package.Replicas > 0 {
		replicas = int32(req.Package.Replicas)
	}

	deploymentSpec := s.buildDeploymentSpec(req.App, req.Version, image, replicas, req.Package)

	_, err := s.clientset.AppsV1().Deployments(s.namespace).Get(ctx, req.App, metav1.GetOptions{})
	if err != nil {
		_, err = s.clientset.AppsV1().Deployments(s.namespace).Create(ctx, deploymentSpec, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create deployment: %w", err)
		}
		return &models.ApplyResponse{
			Success: true,
			Message: fmt.Sprintf("Deployment %s created successfully", req.App),
			App:     req.App,
			Version: req.Version,
		}, nil
	}

	_, err = s.clientset.AppsV1().Deployments(s.namespace).Update(ctx, deploymentSpec, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update deployment: %w", err)
	}

	return &models.ApplyResponse{
		Success: true,
		Message: fmt.Sprintf("Deployment %s updated successfully", req.App),
		App:     req.App,
		Version: req.Version,
	}, nil
}

func (s *OperatorK8sService) buildDeploymentSpec(name, version, image string, replicas int32, pkg models.DeploymentPackage) *appsv1.Deployment {
	labels := map[string]string{
		"app":     name,
		"version": version,
	}

	container := corev1.Container{
		Name:  name,
		Image: image,
	}

	// 处理端口映射
	for _, p := range pkg.Ports {
		container.Ports = append(container.Ports, corev1.ContainerPort{
			ContainerPort: int32(p.ContainerPort),
		})
	}

	// 处理环境变量
	for k, v := range pkg.Environment {
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: s.namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{container},
				},
			},
		},
	}
}

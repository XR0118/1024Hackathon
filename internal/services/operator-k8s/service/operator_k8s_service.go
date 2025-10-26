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

	pkg, ok := req.Pkg["type"].(string)
	if !ok {
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

	healthy := models.HealthStatus{
		Level: 0,
		Msg:   "Not healthy",
	}

	if deployment.Status.ReadyReplicas == deployment.Status.Replicas && deployment.Status.Replicas > 0 {
		healthy.Level = 100
		healthy.Msg = "Healthy"
	} else if deployment.Status.ReadyReplicas > 0 {
		healthy.Level = int(float64(deployment.Status.ReadyReplicas) / float64(deployment.Status.Replicas) * 100)
		healthy.Msg = fmt.Sprintf("%d/%d replicas ready", deployment.Status.ReadyReplicas, deployment.Status.Replicas)
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
	image, ok := req.Pkg["image"].(string)
	if !ok {
		return nil, fmt.Errorf("image is required for docker deployment")
	}

	replicas := int32(1)
	if r, ok := req.Pkg["replicas"].(float64); ok {
		replicas = int32(r)
	}

	deploymentSpec := s.buildDeploymentSpec(req.App, req.Version, image, replicas, req.Pkg)

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

func (s *OperatorK8sService) buildDeploymentSpec(name, version, image string, replicas int32, pkg map[string]interface{}) *appsv1.Deployment {
	labels := map[string]string{
		"app":     name,
		"version": version,
	}

	container := corev1.Container{
		Name:  name,
		Image: image,
	}

	if ports, ok := pkg["ports"].([]interface{}); ok {
		for _, p := range ports {
			if portMap, ok := p.(map[string]interface{}); ok {
				if containerPort, ok := portMap["container_port"].(float64); ok {
					container.Ports = append(container.Ports, corev1.ContainerPort{
						ContainerPort: int32(containerPort),
					})
				}
			}
		}
	}

	if env, ok := pkg["environment"].(map[string]interface{}); ok {
		for k, v := range env {
			if vStr, ok := v.(string); ok {
				container.Env = append(container.Env, corev1.EnvVar{
					Name:  k,
					Value: vStr,
				})
			}
		}
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

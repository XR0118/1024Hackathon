package handler

import (
	"net/http"

	"github.com/boreas/internal/pkg/utils"
	"github.com/boreas/internal/services/operator-k8s/service"
	"github.com/gin-gonic/gin"
)

type OperatorK8sHandler struct {
	operatorService *service.OperatorK8sService
}

func NewOperatorK8sHandler(operatorService *service.OperatorK8sService) *OperatorK8sHandler {
	return &OperatorK8sHandler{
		operatorService: operatorService,
	}
}

func (h *OperatorK8sHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/health", h.HealthCheck)
		api.GET("/ready", h.ReadyCheck)

		// 部署相关API
		deploy := api.Group("/deploy")
		{
			deploy.POST("/:id/execute", h.ExecuteDeployment)
			deploy.GET("/:id/status", h.GetDeploymentStatus)
			deploy.GET("/:id/logs", h.GetDeploymentLogs)
			deploy.POST("/:id/cancel", h.CancelDeployment)
		}
	}
}

// HealthCheck 健康检查
func (h *OperatorK8sHandler) HealthCheck(c *gin.Context) {
	utils.Success(c, gin.H{
		"status":  "healthy",
		"service": "operator-k8s",
	})
}

// ReadyCheck 就绪检查
func (h *OperatorK8sHandler) ReadyCheck(c *gin.Context) {
	// 检查Kubernetes连接状态
	if err := h.operatorService.CheckK8sConnection(); err != nil {
		utils.Error(c, http.StatusServiceUnavailable, "Service not ready: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"status":  "ready",
		"service": "operator-k8s",
	})
}

// ExecuteDeployment 执行Kubernetes部署
func (h *OperatorK8sHandler) ExecuteDeployment(c *gin.Context) {
	deploymentID := c.Param("id")

	// 执行部署
	result, err := h.operatorService.ExecuteDeployment(deploymentID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to execute deployment: "+err.Error())
		return
	}

	utils.Success(c, result)
}

// GetDeploymentStatus 获取部署状态
func (h *OperatorK8sHandler) GetDeploymentStatus(c *gin.Context) {
	deploymentID := c.Param("id")

	status, err := h.operatorService.GetDeploymentStatus(deploymentID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get deployment status: "+err.Error())
		return
	}

	utils.Success(c, status)
}

// GetDeploymentLogs 获取部署日志
func (h *OperatorK8sHandler) GetDeploymentLogs(c *gin.Context) {
	deploymentID := c.Param("id")

	logs, err := h.operatorService.GetDeploymentLogs(deploymentID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get deployment logs: "+err.Error())
		return
	}

	utils.Success(c, logs)
}

// CancelDeployment 取消部署
func (h *OperatorK8sHandler) CancelDeployment(c *gin.Context) {
	deploymentID := c.Param("id")

	err := h.operatorService.CancelDeployment(deploymentID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to cancel deployment: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"message": "Deployment cancelled successfully",
	})
}

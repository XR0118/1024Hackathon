package handler

import (
	"net/http"

	"github.com/boreas/internal/pkg/utils"
	"github.com/boreas/internal/services/operator-pm/service"
	"github.com/gin-gonic/gin"
)

type OperatorPMHandler struct {
	operatorService *service.OperatorPMService
}

func NewOperatorPMHandler(operatorService *service.OperatorPMService) *OperatorPMHandler {
	return &OperatorPMHandler{
		operatorService: operatorService,
	}
}

func (h *OperatorPMHandler) RegisterRoutes(r *gin.Engine) {
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
func (h *OperatorPMHandler) HealthCheck(c *gin.Context) {
	utils.Success(c, gin.H{
		"status":  "healthy",
		"service": "operator-pm",
	})
}

// ReadyCheck 就绪检查
func (h *OperatorPMHandler) ReadyCheck(c *gin.Context) {
	// 检查物理机连接状态
	if err := h.operatorService.CheckPMConnection(); err != nil {
		utils.Error(c, http.StatusServiceUnavailable, "Service not ready: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"status":  "ready",
		"service": "operator-pm",
	})
}

// ExecuteDeployment 执行物理机部署
func (h *OperatorPMHandler) ExecuteDeployment(c *gin.Context) {
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
func (h *OperatorPMHandler) GetDeploymentStatus(c *gin.Context) {
	deploymentID := c.Param("id")

	status, err := h.operatorService.GetDeploymentStatus(deploymentID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get deployment status: "+err.Error())
		return
	}

	utils.Success(c, status)
}

// GetDeploymentLogs 获取部署日志
func (h *OperatorPMHandler) GetDeploymentLogs(c *gin.Context) {
	deploymentID := c.Param("id")

	logs, err := h.operatorService.GetDeploymentLogs(deploymentID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get deployment logs: "+err.Error())
		return
	}

	utils.Success(c, logs)
}

// CancelDeployment 取消部署
func (h *OperatorPMHandler) CancelDeployment(c *gin.Context) {
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

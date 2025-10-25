package handler

import (
	"net/http"

	"github.com/boreas/internal/pkg/models"
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
	api := r.Group("/v1")
	{
		// 健康检查
		api.GET("/health", h.HealthCheck)
		api.GET("/ready", h.ReadyCheck)

		// 核心API
		api.POST("/apply", h.ApplyDeployment)
		api.GET("/status/:app", h.GetApplicationStatus)
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

// ApplyDeployment 应用部署 - 核心API
func (h *OperatorPMHandler) ApplyDeployment(c *gin.Context) {
	var req models.ApplyDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// 执行部署
	result, err := h.operatorService.ApplyDeployment(&req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to apply deployment: "+err.Error())
		return
	}

	utils.Success(c, result)
}

// GetApplicationStatus 获取应用状态 - 核心API
func (h *OperatorPMHandler) GetApplicationStatus(c *gin.Context) {
	appName := c.Param("app")
	if appName == "" {
		utils.Error(c, http.StatusBadRequest, "Application name is required")
		return
	}

	// 获取应用状态
	status, err := h.operatorService.GetApplicationStatus(appName)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get application status: "+err.Error())
		return
	}

	utils.Success(c, status)
}

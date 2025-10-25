package handler

import (
	"net/http"

	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/pkg/utils"
	"github.com/boreas/internal/services/operator-pm-agent/service"
	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	agentService *service.AgentService
}

func NewAgentHandler(agentService *service.AgentService) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
	}
}

func (h *AgentHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		// 应用部署接口
		v1.POST("/apply", h.Apply)

		// 状态查询接口
		v1.GET("/status", h.GetAllStatus)
		v1.GET("/status/:app", h.GetAppStatus)

		// 健康检查接口
		v1.GET("/health", h.HealthCheck)
	}
}

// Apply 应用部署
func (h *AgentHandler) Apply(c *gin.Context) {
	var req models.ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// 执行应用部署
	result, err := h.agentService.ApplyApp(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to apply app: "+err.Error())
		return
	}

	utils.Success(c, result)
}

// GetAllStatus 获取所有应用状态
func (h *AgentHandler) GetAllStatus(c *gin.Context) {
	status, err := h.agentService.GetAllAppStatus()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get app status: "+err.Error())
		return
	}

	utils.Success(c, status)
}

// GetAppStatus 获取指定应用状态
func (h *AgentHandler) GetAppStatus(c *gin.Context) {
	appName := c.Param("app")
	if appName == "" {
		utils.Error(c, http.StatusBadRequest, "App name is required")
		return
	}

	status, err := h.agentService.GetAppStatus(appName)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get app status: "+err.Error())
		return
	}

	utils.Success(c, status)
}

// HealthCheck 健康检查
func (h *AgentHandler) HealthCheck(c *gin.Context) {
	utils.Success(c, gin.H{
		"status":  "healthy",
		"service": "operator-pm-agent",
	})
}

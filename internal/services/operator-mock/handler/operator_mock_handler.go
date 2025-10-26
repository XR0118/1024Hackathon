package handler

import (
	"net/http"

	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/pkg/utils"
	"github.com/boreas/internal/services/operator-mock/service"
	"github.com/gin-gonic/gin"
)

type OperatorMockHandler struct {
	mockService *service.MockDeploymentClient
}

func NewOperatorMockHandler(mockService *service.MockDeploymentClient) *OperatorMockHandler {
	return &OperatorMockHandler{
		mockService: mockService,
	}
}

func (h *OperatorMockHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/apply", h.ApplyDeployment)
		api.GET("/status", h.GetStatus)
	}
}

func (h *OperatorMockHandler) ApplyDeployment(c *gin.Context) {
	var req models.ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	pkg, ok := req.Pkg["deployment"].(map[string]interface{})
	if !ok {
		utils.Error(c, http.StatusBadRequest, "Invalid deployment package format")
		return
	}

	deployPkg := models.DeploymentPackage{}
	if replicas, ok := pkg["replicas"].(float64); ok {
		deployPkg.Replicas = int(replicas)
	}
	if env, ok := pkg["environment"].(map[string]interface{}); ok {
		deployPkg.Environment = make(map[string]string)
		for k, v := range env {
			if strVal, ok := v.(string); ok {
				deployPkg.Environment[k] = strVal
			}
		}
	}

	result, err := h.mockService.Apply(c.Request.Context(), req.App, req.Version, deployPkg)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to apply deployment: "+err.Error())
		return
	}

	utils.Success(c, result)
}

func (h *OperatorMockHandler) GetStatus(c *gin.Context) {
	appName := c.Query("app")

	if appName != "" {
		statuses, err := h.mockService.AppStatus(c.Request.Context(), appName)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "Failed to get application status: "+err.Error())
			return
		}

		utils.Success(c, models.StatusResponse{
			Apps: statuses,
		})
		return
	}

	h.mockService.RLock()
	defer h.mockService.RUnlock()

	var allApps []models.AgentAppStatus
	for _, app := range h.mockService.Instances {
		allApps = append(allApps, app)
	}

	utils.Success(c, models.StatusResponse{
		Apps: allApps,
	})
}

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
	api := r.Group("/v1")
	{
		api.POST("/apply", h.ApplyDeployment)
		api.GET("/status", h.GetStatus)
	}
}

func (h *OperatorMockHandler) ApplyDeployment(c *gin.Context) {
	var req struct {
		App     string                   `json:"app" binding:"required"`
		Version string                   `json:"version" binding:"required"`
		Pkg     models.DeploymentPackage `json:"pkg" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	result, err := h.mockService.Apply(c.Request.Context(), req.App, req.Version, req.Pkg)
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

	allStatuses, err := h.mockService.AllStatus(c.Request.Context())
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get all status: "+err.Error())
		return
	}

	utils.Success(c, models.StatusResponse{
		Apps: allStatuses,
	})
}

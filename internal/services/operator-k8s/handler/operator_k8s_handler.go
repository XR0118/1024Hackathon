package handler

import (
	"net/http"

	"github.com/boreas/internal/pkg/models"
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
	v1 := r.Group("/v1")
	{
		v1.GET("/health", h.HealthCheck)
		v1.GET("/ready", h.ReadyCheck)
		v1.POST("/apply", h.Apply)
		v1.GET("/status/:app", h.GetStatus)
	}
}

func (h *OperatorK8sHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "operator-k8s",
	})
}

func (h *OperatorK8sHandler) ReadyCheck(c *gin.Context) {
	if err := h.operatorService.CheckK8sConnection(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not ready",
			"service": "operator-k8s",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"service": "operator-k8s",
	})
}

func (h *OperatorK8sHandler) Apply(c *gin.Context) {
	var req models.ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ApplyResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	resp, err := h.operatorService.Apply(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApplyResponse{
			Success: false,
			Message: "Failed to apply: " + err.Error(),
			App:     req.App,
			Version: req.Version,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OperatorK8sHandler) GetStatus(c *gin.Context) {
	app := c.Param("app")

	status, err := h.operatorService.GetStatus(app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

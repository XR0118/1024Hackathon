package handler

import (
	"net/http"
	"strconv"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

type deploymentHandler struct {
	deploymentService interfaces.DeploymentService
}

// NewDeploymentHandler 创建部署处理器
func NewDeploymentHandler(deploymentService interfaces.DeploymentService) *deploymentHandler {
	return &deploymentHandler{
		deploymentService: deploymentService,
	}
}

// CreateDeployment 创建部署
func (h *deploymentHandler) CreateDeployment(c *gin.Context) {
	var req models.CreateDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	deployment, err := h.deploymentService.CreateDeployment(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "DEPLOYMENT_CREATION_FAILED", err.Error(), nil)
		return
	}

	utils.Created(c, deployment)
}

// GetDeploymentList 获取部署列表
func (h *deploymentHandler) GetDeploymentList(c *gin.Context) {
	var req models.ListDeploymentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	// 解析分页参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			req.PageSize = pageSize
		}
	}

	response, err := h.deploymentService.GetDeploymentList(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "DEPLOYMENT_LIST_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, response)
}

// GetDeployment 获取部署详情
func (h *deploymentHandler) GetDeployment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "deployment id is required")
		return
	}

	deployment, err := h.deploymentService.GetDeployment(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "DEPLOYMENT_NOT_FOUND", err.Error(), nil)
		return
	}

	utils.Success(c, deployment)
}

// StartDeployment 开始部署
func (h *deploymentHandler) StartDeployment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "deployment id is required")
		return
	}

	deployment, err := h.deploymentService.StartDeployment(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "DEPLOYMENT_START_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, deployment)
}

// RollbackDeployment 回滚部署
func (h *deploymentHandler) RollbackDeployment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "deployment id is required")
		return
	}

	var req models.RollbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	deployment, err := h.deploymentService.RollbackDeployment(c.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "DEPLOYMENT_ROLLBACK_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, deployment)
}

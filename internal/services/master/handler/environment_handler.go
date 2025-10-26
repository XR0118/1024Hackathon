package handler

import (
	"net/http"
	"strconv"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

type environmentHandler struct {
	environmentService interfaces.EnvironmentService
}

// NewEnvironmentHandler 创建环境处理器
func NewEnvironmentHandler(environmentService interfaces.EnvironmentService) *environmentHandler {
	return &environmentHandler{
		environmentService: environmentService,
	}
}

// CreateEnvironment 创建环境
func (h *environmentHandler) CreateEnvironment(c *gin.Context) {
	var req models.CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	environment, err := h.environmentService.CreateEnvironment(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "ENVIRONMENT_CREATION_FAILED", err.Error(), nil)
		return
	}

	utils.Created(c, environment)
}

// GetEnvironmentList 获取环境列表
func (h *environmentHandler) GetEnvironmentList(c *gin.Context) {
	var req models.ListEnvironmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	// 解析分页参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	} else {
		req.Page = 1
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			req.PageSize = pageSize
		}
	} else {
		req.PageSize = 100
	}

	// 解析布尔参数
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			req.IsActive = &isActive
		}
	}

	response, err := h.environmentService.GetEnvironmentList(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "ENVIRONMENT_LIST_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, response)
}

// GetEnvironment 获取环境详情
func (h *environmentHandler) GetEnvironment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "environment id is required")
		return
	}

	environment, err := h.environmentService.GetEnvironment(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "ENVIRONMENT_NOT_FOUND", err.Error(), nil)
		return
	}

	utils.Success(c, environment)
}

// UpdateEnvironment 更新环境
func (h *environmentHandler) UpdateEnvironment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "environment id is required")
		return
	}

	var req models.UpdateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	environment, err := h.environmentService.UpdateEnvironment(c.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "ENVIRONMENT_UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, environment)
}

// DeleteEnvironment 删除环境
func (h *environmentHandler) DeleteEnvironment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "environment id is required")
		return
	}

	err := h.environmentService.DeleteEnvironment(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "ENVIRONMENT_DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, gin.H{"message": "Environment deleted successfully"})
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

type applicationHandler struct {
	applicationService interfaces.ApplicationService
}

// NewApplicationHandler 创建应用处理器
func NewApplicationHandler(applicationService interfaces.ApplicationService) *applicationHandler {
	return &applicationHandler{
		applicationService: applicationService,
	}
}

// CreateApplication 创建应用
func (h *applicationHandler) CreateApplication(c *gin.Context) {
	var req models.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	application, err := h.applicationService.CreateApplication(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "APPLICATION_CREATION_FAILED", err.Error(), nil)
		return
	}

	utils.Created(c, application)
}

// GetApplicationList 获取应用列表
func (h *applicationHandler) GetApplicationList(c *gin.Context) {
	var req models.ListApplicationsRequest
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

	response, err := h.applicationService.GetApplicationList(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "APPLICATION_LIST_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, response)
}

// GetApplication 获取应用详情（按应用名称查询）
func (h *applicationHandler) GetApplication(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		utils.BadRequest(c, "application name is required")
		return
	}

	// 按应用名称查询
	application, err := h.applicationService.GetApplicationByName(c.Request.Context(), name)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "APPLICATION_NOT_FOUND", err.Error(), nil)
		return
	}

	utils.Success(c, application)
}

// UpdateApplication 更新应用
func (h *applicationHandler) UpdateApplication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "application id is required")
		return
	}

	var req models.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	application, err := h.applicationService.UpdateApplication(c.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "APPLICATION_UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, application)
}

// DeleteApplication 删除应用
func (h *applicationHandler) DeleteApplication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "application id is required")
		return
	}

	err := h.applicationService.DeleteApplication(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "APPLICATION_DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, gin.H{"message": "Application deleted successfully"})
}

// GetApplicationVersions 获取应用版本详情（按环境组织）
func (h *applicationHandler) GetApplicationVersions(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		utils.BadRequest(c, "application name is required")
		return
	}

	response, err := h.applicationService.GetApplicationVersionsDetail(c.Request.Context(), name)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "GET_VERSIONS_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, response)
}

// GetApplicationVersionsSummary 获取应用版本概要
func (h *applicationHandler) GetApplicationVersionsSummary(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		utils.BadRequest(c, "application name is required")
		return
	}

	response, err := h.applicationService.GetApplicationVersionsSummary(c.Request.Context(), name)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "GET_VERSIONS_SUMMARY_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, response)
}

// GetApplicationVersionCoverage 获取应用指定版本的覆盖率（累积覆盖率）
func (h *applicationHandler) GetApplicationVersionCoverage(c *gin.Context) {
	name := c.Param("name")
	version := c.Param("version")

	if name == "" {
		utils.BadRequest(c, "application name is required")
		return
	}
	if version == "" {
		utils.BadRequest(c, "version is required")
		return
	}

	response, err := h.applicationService.GetApplicationVersionCoverage(c.Request.Context(), name, version)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "GET_VERSION_COVERAGE_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, response)
}

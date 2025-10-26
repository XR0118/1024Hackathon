package handler

import (
	"net/http"
	"strconv"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/pkg/utils"
	"github.com/boreas/internal/services/master/mock"
	"github.com/gin-gonic/gin"
)

type applicationHandler struct {
	applicationService interfaces.ApplicationService
	versionService     interfaces.VersionService
}

// NewApplicationHandler 创建应用处理器
func NewApplicationHandler(applicationService interfaces.ApplicationService, versionService interfaces.VersionService) *applicationHandler {
	return &applicationHandler{
		applicationService: applicationService,
		versionService:     versionService,
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

// GetApplication 获取应用详情
func (h *applicationHandler) GetApplication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "application id is required")
		return
	}

	application, err := h.applicationService.GetApplication(c.Request.Context(), id)
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

// GetApplicationVersions 获取应用的版本信息（从 Operator 查询）
// 使用应用名称作为查询参数
func (h *applicationHandler) GetApplicationVersions(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		utils.BadRequest(c, "application name is required")
		return
	}

	app, err := h.applicationService.GetApplicationByName(c, name)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "APPLICATION_NOT_FOUND", err.Error(), nil)
		return
	}

	appStatus, err := mock.MockAgent.AppStatus(c, name)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "APPLICATION_STATUS_FAILED", err.Error(), nil)
		return
	}

	var versions = []models.ApplicationVersionInfo{}
	var total int
	var replicas = []int{}
	// 从 appStatus 中提取版本信息
	for _, status := range appStatus {
		tmpV, _ := h.versionService.GetVersion(c, status.Version)
		if tmpV == nil {
			continue
		}
		versions = append(versions, models.ApplicationVersionInfo{
			Version:       status.Version,
			Status:        tmpV.Status,
			Health:        status.Healthy.Level,
			LastUpdatedAt: status.Updated.Format("2006-01-02 15:04:05"),
		})
		total += status.Replicas
		replicas = append(replicas, status.Replicas)
	}

	for i := range versions {
		versions[i].Coverage = int(float64(replicas[i]) / float64(total) * 100)
	}

	response := models.ApplicationVersionsResponse{
		ApplicationID: app.ID, // 后续从查询结果获取
		Name:          name,
		Versions:      versions,
	}

	utils.Success(c, response)
}

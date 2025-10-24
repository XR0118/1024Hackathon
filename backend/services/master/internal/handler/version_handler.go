package handler

import (
	"net/http"
	"strconv"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/models"
	"github.com/boreas/internal/utils"
	"github.com/gin-gonic/gin"
)

type versionHandler struct {
	versionService interfaces.VersionService
}

// NewVersionHandler 创建版本处理器
func NewVersionHandler(versionService interfaces.VersionService) *versionHandler {
	return &versionHandler{
		versionService: versionService,
	}
}

// CreateVersion 创建版本
func (h *versionHandler) CreateVersion(c *gin.Context) {
	var req models.CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	version, err := h.versionService.CreateVersion(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "VERSION_CREATION_FAILED", err.Error(), nil)
		return
	}

	utils.Created(c, version)
}

// GetVersionList 获取版本列表
func (h *versionHandler) GetVersionList(c *gin.Context) {
	var req models.ListVersionsRequest
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

	response, err := h.versionService.GetVersionList(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "VERSION_LIST_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, response)
}

// GetVersion 获取版本详情
func (h *versionHandler) GetVersion(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "version id is required")
		return
	}

	version, err := h.versionService.GetVersion(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "VERSION_NOT_FOUND", err.Error(), nil)
		return
	}

	utils.Success(c, version)
}

// DeleteVersion 删除版本
func (h *versionHandler) DeleteVersion(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "version id is required")
		return
	}

	err := h.versionService.DeleteVersion(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "VERSION_DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, gin.H{"message": "Version deleted successfully"})
}

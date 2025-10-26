package handler

import (
	"net/http"
	"strconv"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/pkg/utils"
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
	version := c.Param("version")
	if version == "" {
		utils.BadRequest(c, "version is required")
		return
	}

	versionObj, err := h.versionService.GetVersion(c.Request.Context(), version)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "VERSION_NOT_FOUND", err.Error(), nil)
		return
	}

	utils.Success(c, versionObj)
}

// DeleteVersion 删除版本
func (h *versionHandler) DeleteVersion(c *gin.Context) {
	version := c.Param("version")
	if version == "" {
		utils.BadRequest(c, "version is required")
		return
	}

	err := h.versionService.DeleteVersion(c.Request.Context(), version)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "VERSION_DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, gin.H{"message": "Version deleted successfully"})
}

// RollbackVersion 回滚版本
func (h *versionHandler) RollbackVersion(c *gin.Context) {
	version := c.Param("version")
	if version == "" {
		utils.BadRequest(c, "version is required")
		return
	}

	var req models.RollbackVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	// TODO: 实现版本回滚逻辑
	// 1. 获取目标版本信息
	// 2. 创建新的部署任务，将应用回滚到目标版本
	// 3. 更新版本状态为 revert
	// 4. 记录回滚原因和操作者

	// 暂时返回成功响应
	utils.Success(c, gin.H{
		"message": "Rollback initiated",
		"version": version,
		"reason":  req.Reason,
	})
}

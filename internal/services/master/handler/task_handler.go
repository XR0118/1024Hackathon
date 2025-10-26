package handler

import (
	"net/http"
	"strconv"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

type taskHandler struct {
	taskService interfaces.TaskService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(taskService interfaces.TaskService) *taskHandler {
	return &taskHandler{
		taskService: taskService,
	}
}

// GetTaskList 获取任务列表
func (h *taskHandler) GetTaskList(c *gin.Context) {
	var req models.ListTasksRequest
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

	response, err := h.taskService.GetTaskList(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "TASK_LIST_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, response)
}

// GetTask 获取任务详情
func (h *taskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "task id is required")
		return
	}

	task, err := h.taskService.GetTask(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "TASK_NOT_FOUND", err.Error(), nil)
		return
	}

	utils.Success(c, task)
}

// RetryTask 重试任务
func (h *taskHandler) RetryTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "task id is required")
		return
	}

	task, err := h.taskService.RetryTask(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "TASK_RETRY_FAILED", err.Error(), nil)
		return
	}

	utils.Success(c, task)
}

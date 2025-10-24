package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/boreas/internal/models"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithDetails 带详情的错误响应
func ErrorWithDetails(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    details,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden 403 错误
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

// NotFound 404 错误
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// Conflict 409 错误
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message)
}

// InternalServerError 500 错误
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

// ServiceUnavailable 503 错误
func ServiceUnavailable(c *gin.Context, message string) {
	Error(c, http.StatusServiceUnavailable, message)
}

// ValidationError 验证错误
func ValidationError(c *gin.Context, err error) {
	BadRequest(c, "Validation failed: "+err.Error())
}

// DatabaseError 数据库错误
func DatabaseError(c *gin.Context, err error) {
	InternalServerError(c, "Database operation failed: "+err.Error())
}

// ErrorResponse 错误响应（使用标准错误格式）
func ErrorResponse(c *gin.Context, code int, errorCode, message string, details map[string]interface{}) {
	c.JSON(code, models.ErrorResponse{
		Error: models.ErrorDetail{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
	})
}

// PaginationResponse 分页响应
func PaginationResponse(c *gin.Context, data interface{}, total, page, pageSize int) {
	Success(c, gin.H{
		"data":      data,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

package service

// workflowController 结构体，用于管理 workflow 相关操作
type workflowController struct {
	taskService *taskService // 注入 taskService 以处理任务相关逻辑
}

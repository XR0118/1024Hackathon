package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/pkg/utils"
	"github.com/google/uuid"
)

type applicationService struct {
	appRepo         interfaces.ApplicationRepository
	versionRepo     interfaces.VersionRepository
	deploymentRepo  interfaces.DeploymentRepository
	operatorManager interfaces.OperatorManager
}

// NewApplicationService 创建应用服务
func NewApplicationService(
	appRepo interfaces.ApplicationRepository,
	versionRepo interfaces.VersionRepository,
	deploymentRepo interfaces.DeploymentRepository,
	operatorManager interfaces.OperatorManager,
) interfaces.ApplicationService {
	return &applicationService{
		appRepo:         appRepo,
		versionRepo:     versionRepo,
		deploymentRepo:  deploymentRepo,
		operatorManager: operatorManager,
	}
}

func (s *applicationService) CreateApplication(ctx context.Context, req *models.CreateApplicationRequest) (*models.Application, error) {
	// 验证请求
	if err := utils.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 检查应用名称是否已存在
	existingApps, _, err := s.appRepo.List(ctx, &models.ApplicationFilter{
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check existing applications: %w", err)
	}

	for _, app := range existingApps {
		if app.Name == req.Name {
			return nil, fmt.Errorf("application with name %s already exists", req.Name)
		}
	}

	cfg, _ := json.Marshal(req.Config)
	// 创建应用
	app := &models.Application{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Repository:  req.Repository,
		Type:        req.Type,
		Config:      cfg,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.appRepo.Create(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return app, nil
}

func (s *applicationService) GetApplicationList(ctx context.Context, req *models.ListApplicationsRequest) (*models.ApplicationListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	filter := &models.ApplicationFilter{
		Repository: req.Repository,
		Type:       req.Type,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	applications, total, err := s.appRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}

	return &models.ApplicationListResponse{
		Applications: applications,
		Total:        total,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}, nil
}

func (s *applicationService) GetApplication(ctx context.Context, id string) (*models.Application, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	return app, nil
}

func (s *applicationService) GetApplicationByName(ctx context.Context, name string) (*models.Application, error) {
	app, err := s.appRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	return app, nil
}

func (s *applicationService) UpdateApplication(ctx context.Context, id string, req *models.UpdateApplicationRequest) (*models.Application, error) {
	// 验证请求
	if err := utils.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有应用
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// 更新字段
	if req.Name != "" {
		app.Name = req.Name
	}
	if req.Description != "" {
		app.Description = req.Description
	}
	if req.Type != "" {
		app.Type = req.Type
	}
	if req.Config != nil {
		bs, _ := json.Marshal(req.Config)
		app.Config = bs
	}
	app.UpdatedAt = time.Now()

	// 保存更新
	if err := s.appRepo.Update(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}

	return app, nil
}

func (s *applicationService) DeleteApplication(ctx context.Context, id string) error {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// 删除应用
	if err := s.appRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	return nil
}

// GetApplicationVersionsSummary 获取应用版本概要（包含跨环境汇总的运行时信息）
func (s *applicationService) GetApplicationVersionsSummary(ctx context.Context, appName string) (*models.ApplicationVersionsSummaryResponse, error) {
	// 1. 获取应用信息（包括关联的环境）
	app, err := s.appRepo.GetByName(ctx, appName)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// 2. 查询该应用的所有版本（通过 app_builds 包含该应用的版本）
	versions, _, err := s.versionRepo.List(ctx, &models.VersionFilter{
		Repository: app.Repository,
		Page:       1,
		PageSize:   100, // 获取最近100个版本
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	// 3. 创建版本映射（用于快速查找版本静态信息）
	versionMap := make(map[string]*models.Version)
	for _, v := range versions {
		// 检查该版本是否包含当前应用
		appBuilds := v.GetAppBuilds()
		hasApp := false
		for _, build := range appBuilds {
			if build.AppName == appName {
				hasApp = true
				break
			}
		}
		if hasApp {
			versionMap[v.Version] = v
		}
	}

	// 4. 从各个环境查询应用状态，按版本汇总
	type VersionStats struct {
		TotalInstances   int
		HealthyInstances int
		EnvironmentCount int
		HealthSum        int // 用于计算平均健康度
		HealthCount      int // 健康度样本数
		LastDeployedAt   time.Time
	}

	versionStats := make(map[string]*VersionStats)

	for _, env := range app.Environments {
		// 通过 operator manager 获取应用在该环境的状态
		appStatus, err := s.operatorManager.GetApplicationStatus(ctx, env.ID, appName)
		if err != nil {
			// 如果查询失败，跳过该环境
			continue
		}

		// 汇总该环境的版本信息
		for _, versionStatus := range appStatus.Versions {
			stats, exists := versionStats[versionStatus.Version]
			if !exists {
				stats = &VersionStats{}
				versionStats[versionStatus.Version] = stats
			}

			// 统计实例数
			instanceCount := len(versionStatus.Nodes)
			stats.TotalInstances += instanceCount
			stats.EnvironmentCount++

			// 统计健康实例数（健康度 >= 80 认为健康）
			for _, node := range versionStatus.Nodes {
				if node.Healthy.Level >= 80 {
					stats.HealthyInstances++
				}
				stats.HealthSum += node.Healthy.Level
				stats.HealthCount++
			}

			// 记录最后部署时间（简化处理，使用当前时间）
			if stats.LastDeployedAt.IsZero() {
				stats.LastDeployedAt = time.Now()
			}
		}
	}

	// 5. 构建版本概要列表（只包含核心运行时指标）
	summaries := make([]models.VersionSummary, 0, len(versionMap))
	for versionNum, v := range versionMap {
		stats := versionStats[versionNum]

		summary := models.VersionSummary{
			Version: v.Version,
			Status:  v.Status,
		}

		// 计算运行时指标
		if stats != nil {
			// 计算健康度百分比
			if stats.HealthCount > 0 {
				summary.HealthPercent = float64(stats.HealthSum) / float64(stats.HealthCount)
			}

			// 计算覆盖度百分比（部署环境数 / 总环境数）
			if len(app.Environments) > 0 {
				summary.CoveragePercent = float64(stats.EnvironmentCount) / float64(len(app.Environments)) * 100
			}
		}

		summaries = append(summaries, summary)
	}

	return &models.ApplicationVersionsSummaryResponse{
		ApplicationID:   app.ID,
		ApplicationName: app.Name,
		Versions:        summaries,
	}, nil
}

// GetApplicationVersionsDetail 获取应用版本详情（按环境组织）
func (s *applicationService) GetApplicationVersionsDetail(ctx context.Context, appName string) (*models.ApplicationVersionsDetailResponse, error) {
	// 1. 获取应用信息（包括关联的环境）
	app, err := s.appRepo.GetByName(ctx, appName)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// 2. 对每个环境查询版本信息
	environments := make([]models.EnvironmentVersions, 0, len(app.Environments))
	for _, env := range app.Environments {
		// 通过 operator manager 获取应用在该环境的状态
		appStatus, err := s.operatorManager.GetApplicationStatus(ctx, env.ID, appName)
		if err != nil {
			// 如果查询失败，记录错误但继续处理其他环境
			// 返回空的版本列表
			environments = append(environments, models.EnvironmentVersions{
				Environment: env,
				Versions:    []models.EnvironmentVersionDetail{},
			})
			continue
		}

		// 转换 VersionStatus 为 EnvironmentVersionDetail
		versions := make([]models.EnvironmentVersionDetail, 0, len(appStatus.Versions))
		for _, versionStatus := range appStatus.Versions {
			// 转换实例信息
			instances := make([]models.VersionInstance, 0, len(versionStatus.Nodes))
			for _, nodeStatus := range versionStatus.Nodes {
				instances = append(instances, models.VersionInstance{
					NodeName:      nodeStatus.Node,
					Health:        nodeStatus.Healthy.Level,
					Status:        "running", // 默认状态
					LastUpdatedAt: time.Now(),
				})
			}

			// 查询版本信息以获取 git tag 和 commit
			version, err := s.versionRepo.GetByVersion(ctx, versionStatus.Version)
			var gitTag, gitCommit, versionStatusStr string
			if err == nil {
				gitTag = version.GitTag
				gitCommit = version.GitCommit
				versionStatusStr = version.Status
			}

			versions = append(versions, models.EnvironmentVersionDetail{
				Version:       versionStatus.Version,
				Status:        versionStatusStr,
				GitTag:        gitTag,
				GitCommit:     gitCommit,
				Instances:     instances,
				Health:        versionStatus.Healthy.Level,
				Coverage:      int(versionStatus.Percent), // 覆盖率从 operator 返回
				LastUpdatedAt: time.Now(),
			})
		}

		environments = append(environments, models.EnvironmentVersions{
			Environment: env,
			Versions:    versions,
		})
	}

	return &models.ApplicationVersionsDetailResponse{
		ApplicationID:   app.ID,
		ApplicationName: app.Name,
		Environments:    environments,
	}, nil
}

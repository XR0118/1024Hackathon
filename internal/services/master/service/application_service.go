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
			// 计算健康度：跨环境的节点健康度加权平均
			healthLevel := 0
			if stats.HealthCount > 0 {
				healthLevel = stats.HealthSum / stats.HealthCount
			}
			summary.Healthy = models.HealthInfo{
				Level: healthLevel,
				Msg:   fmt.Sprintf("Average health across %d environment(s), %d instance(s)", stats.EnvironmentCount, stats.TotalInstances),
			}

			// 计算覆盖度百分比（部署环境数 / 总环境数）
			if len(app.Environments) > 0 {
				summary.CoveragePercent = float64(stats.EnvironmentCount) / float64(len(app.Environments)) * 100
			}
		} else {
			// 没有统计数据时，健康度为 0
			summary.Healthy = models.HealthInfo{
				Level: 0,
				Msg:   "No deployment data",
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

		// 计算该环境的总实例数（所有版本的节点数之和）
		totalInstances := 0
		for _, versionStatus := range appStatus.Versions {
			totalInstances += len(versionStatus.Nodes)
		}

		// 转换 VersionStatus 为 EnvironmentVersionDetail
		versions := make([]models.EnvironmentVersionDetail, 0, len(appStatus.Versions))
		for _, versionStatus := range appStatus.Versions {
			// 转换实例信息
			instances := make([]models.VersionInstance, 0, len(versionStatus.Nodes))
			healthSum := 0
			for _, nodeStatus := range versionStatus.Nodes {
				instances = append(instances, models.VersionInstance{
					NodeName:      nodeStatus.Node,
					Healthy:       nodeStatus.Healthy, // 直接使用 HealthInfo
					Status:        "running",          // 默认状态
					LastUpdatedAt: time.Now(),
				})
				healthSum += nodeStatus.Healthy.Level
			}

			// 计算该版本在此环境的健康度：各实例健康度的加权平均
			versionHealthLevel := 0
			if len(instances) > 0 {
				versionHealthLevel = healthSum / len(instances)
			}

			// 计算该版本在此环境的覆盖率：该版本节点数 / 总节点数 * 100
			coverage := 0
			if totalInstances > 0 {
				coverage = len(versionStatus.Nodes) * 100 / totalInstances
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
				Version:   versionStatus.Version,
				Status:    versionStatusStr,
				GitTag:    gitTag,
				GitCommit: gitCommit,
				Instances: instances,
				Healthy: models.HealthInfo{
					Level: versionHealthLevel,
					Msg:   fmt.Sprintf("Average health of %d instance(s) in this environment", len(instances)),
				},
				Coverage:      coverage, // 根据节点数计算覆盖率
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

// GetApplicationVersionCoverage 获取应用指定版本的覆盖率（累积覆盖率）
// 计算逻辑：运行实例的版本创建时间 >= 目标版本创建时间，则认为该实例被目标版本覆盖
func (s *applicationService) GetApplicationVersionCoverage(ctx context.Context, appName string, targetVersion string) (*models.VersionCoverageResponse, error) {
	// 1. 获取应用信息（包括关联的环境）
	app, err := s.appRepo.GetByName(ctx, appName)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// 2. 获取目标版本信息（用于比较创建时间）
	targetVersionInfo, err := s.versionRepo.GetByVersion(ctx, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("target version not found: %w", err)
	}

	// 3. 初始化响应
	response := &models.VersionCoverageResponse{
		ApplicationID:     app.ID,
		ApplicationName:   app.Name,
		TargetVersion:     targetVersion,
		TotalEnvironments: len(app.Environments),
		Environments:      make([]models.EnvironmentVersionCoverage, 0, len(app.Environments)),
	}

	// 4. 遍历每个环境，计算覆盖情况
	coveredEnvCount := 0
	for _, env := range app.Environments {
		// 查询该环境的应用状态
		appStatus, err := s.operatorManager.GetApplicationStatus(ctx, env.ID, appName)
		if err != nil {
			// 查询失败，标记为未覆盖
			response.Environments = append(response.Environments, models.EnvironmentVersionCoverage{
				Environment:         env,
				IsCovered:           false,
				CurrentVersion:      "",
				TotalInstances:      0,
				CoveredInstances:    0,
				CoveragePercent:     0,
				VersionDistribution: []models.VersionInstanceCount{},
			})
			continue
		}

		// 计算该环境的覆盖情况
		envCoverage := s.calculateEnvironmentCoverage(ctx, env, appStatus, targetVersionInfo)
		response.Environments = append(response.Environments, envCoverage)

		if envCoverage.IsCovered {
			coveredEnvCount++
		}
	}

	// 5. 计算总体覆盖率
	response.CoveredEnvironments = coveredEnvCount
	if response.TotalEnvironments > 0 {
		response.CoveragePercent = float64(coveredEnvCount) / float64(response.TotalEnvironments) * 100
	}

	return response, nil
}

// calculateEnvironmentCoverage 计算单个环境的版本覆盖情况
// 基于创建时间比较：如果运行版本的创建时间 >= 目标版本的创建时间，则认为被覆盖
func (s *applicationService) calculateEnvironmentCoverage(ctx context.Context, env models.Environment, appStatus *models.ApplicationStatusResponse, targetVersionInfo *models.Version) models.EnvironmentVersionCoverage {
	// 统计总实例数
	totalInstances := 0
	for _, versionStatus := range appStatus.Versions {
		totalInstances += len(versionStatus.Nodes)
	}

	// 统计覆盖实例数和版本分布
	coveredInstances := 0
	currentVersion := ""
	var currentVersionTime *models.Version
	versionDistribution := make([]models.VersionInstanceCount, 0, len(appStatus.Versions))

	for _, versionStatus := range appStatus.Versions {
		instanceCount := len(versionStatus.Nodes)

		// 查询该版本信息以获取创建时间
		versionInfo, err := s.versionRepo.GetByVersion(ctx, versionStatus.Version)
		isCovered := false
		if err == nil {
			// 比较创建时间：运行版本创建时间 >= 目标版本创建时间
			isCovered = !versionInfo.CreatedAt.Before(targetVersionInfo.CreatedAt)

			if isCovered {
				coveredInstances += instanceCount
			}

			// 记录当前最新版本（创建时间最晚的）
			if currentVersionTime == nil || versionInfo.CreatedAt.After(currentVersionTime.CreatedAt) {
				currentVersion = versionStatus.Version
				currentVersionTime = versionInfo
			}
		}

		versionDistribution = append(versionDistribution, models.VersionInstanceCount{
			Version:       versionStatus.Version,
			InstanceCount: instanceCount,
			IsCovered:     isCovered,
		})
	}

	// 计算该环境的覆盖率
	coveragePercent := 0.0
	if totalInstances > 0 {
		coveragePercent = float64(coveredInstances) / float64(totalInstances) * 100
	}

	// 判断环境是否被覆盖（至少有一个实例运行 >= 目标版本）
	isCovered := coveredInstances > 0

	return models.EnvironmentVersionCoverage{
		Environment:         env,
		IsCovered:           isCovered,
		CurrentVersion:      currentVersion,
		TotalInstances:      totalInstances,
		CoveredInstances:    coveredInstances,
		CoveragePercent:     coveragePercent,
		VersionDistribution: versionDistribution,
	}
}

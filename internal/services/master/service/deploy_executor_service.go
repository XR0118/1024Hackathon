package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/services/master/mock"
)

type DeployExecutor interface {
	Apply(ctx context.Context) error
}

type DeployClient interface {
	Apply(ctx context.Context, app string, version string, pkg models.DeploymentPackage) (models.ApplyResponse, error)
	AppStatus(ctx context.Context, app string) ([]models.AgentAppStatus, error)
}

type SimpleDeployExecutor struct {
	task           models.Task
	deploymentRepo interfaces.DeploymentRepository
	client         DeployClient
}

func NewSimpleDeployExecutor(task models.Task, deploymentRepo interfaces.DeploymentRepository, client DeployClient) *SimpleDeployExecutor {
	return &SimpleDeployExecutor{
		task:           task,
		deploymentRepo: deploymentRepo,
		client:         mock.NewMockDeploymentClient(),
	}
}

func (e *SimpleDeployExecutor) Apply(ctx context.Context) error {
	status := func() (int, map[string]int, []models.AgentAppStatus) {
		appStatus, _ := e.client.AppStatus(ctx, e.task.AppID)
		total := 0
		versionMap := make(map[string]int)
		for _, s := range appStatus {
			total += s.Replicas
			versionMap[s.Version] += s.Replicas
		}
		return total, versionMap, appStatus
	}

	pkg := models.DeploymentPackage{}
	_ = json.Unmarshal([]byte(e.task.Payload), &pkg)

	total, versionMap, _ := status()
	if total == 0 {
		for _, s := range e.task.Deployment.Strategy {
			total += s.BatchSize
		}
		_, err := e.client.Apply(ctx, e.task.AppID, e.task.Deployment.ID, pkg)
		return err
	}

	if versionMap[e.task.Deployment.ID] == total {
		return nil
	}

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			deployment, err := e.deploymentRepo.GetByID(ctx, e.task.DeploymentID)
			if err != nil {
				return err
			}
			if deployment.Status == models.DeploymentStatusFailed ||
				deployment.Status == models.DeploymentStatusCancelled {
				return nil
			}

			total, versionMap, appStatus := status()
			if deployment.Status == models.DeploymentStatusRolledBack {
				// todo: 找到最近版本回滚
				v := ""
				for k := range versionMap {
					if k != e.task.Deployment.ID {
						v = k
						break
					}
				}
				versionMap[v] += versionMap[deployment.ID]
				versionMap[deployment.ID] = 0
				pkg.Replicas = 0
				_, err := e.client.Apply(ctx, e.task.AppID, deployment.ID, pkg)
				if err != nil {
					return err
				}
				pkg.Replicas = versionMap[v]
				_, err = e.client.Apply(ctx, e.task.AppID, v, pkg)
				return err
			}
			if versionMap[e.task.Deployment.ID] == total {
				return nil
			}

			s, idx := func() (models.DeploySteps, int) {
				replicas := versionMap[e.task.Deployment.ID]
				raito := float64(replicas) / float64(total)
				i := 0
				for j, s := range e.task.Deployment.Strategy {
					if raito < s.CanaryRatio || replicas < s.BatchSize {
						return s, j
					}
					i = j
				}
				s := e.task.Deployment.Strategy[i]
				return s, i
			}()

			var unhealthy bool
			var lastUpdate time.Time
			for _, status := range appStatus {
				if status.Version == e.task.Deployment.ID {
					lastUpdate = status.Updated
					if status.Healthy.Level < 80 {
						unhealthy = true
						break
					}
				}
			}
			if unhealthy && s.AutoRollback {
				deployment.Status = models.DeploymentStatusRolledBack
				_ = e.deploymentRepo.Update(ctx, deployment)
				continue
			}

			if time.Since(lastUpdate).Seconds() < float64(s.BatchInterval) ||
				s.ManualApprovalStatus != nil && !*s.ManualApprovalStatus {
				continue
			}

			if idx == len(e.task.Deployment.Strategy)-1 {
				continue
			}

			ns := e.task.Deployment.Strategy[idx+1]
			if ns.BatchSize > 0 {
				pkg.Replicas = ns.BatchSize
			} else {
				pkg.Replicas = int(float64(total) * ns.CanaryRatio)
			}
			_, err = e.client.Apply(ctx, e.task.AppID, e.task.Deployment.ID, pkg)
			if err != nil {
				return err
			}
		}
	}
}

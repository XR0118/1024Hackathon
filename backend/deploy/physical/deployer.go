package physical

import (
	"github.com/XR0118/1024Hackathon/backend/model"
)

type PhysicalDeployer interface {
	Deploy(app *model.Application, version *model.Version, env *model.TargetEnvironment) error
	BuildArtifact(app *model.Application, version *model.Version) (*Artifact, error)
	DeployToHost(host *model.Host, artifact *Artifact, app *model.Application) error
	HealthCheck(app *model.Application, env *model.TargetEnvironment) bool
	Rollback(app *model.Application, env *model.TargetEnvironment, targetVersion string) error
}

type Artifact struct {
	Path    string
	Name    string
	Version string
	Size    int64
	Hash    string
}

type SSHClient interface {
	Connect(host *model.Host, config *model.SSHConfig) error
	Disconnect() error
	Upload(localPath, remotePath string) error
	Download(remotePath, localPath string) error
	Execute(command string) (string, error)
	ExecuteScript(scriptPath string) (string, error)
}

type ProcessManager interface {
	Start(service string) error
	Stop(service string) error
	Restart(service string) error
	Status(service string) (*ServiceStatus, error)
	CheckHealth(endpoint string) bool
}

type ServiceStatus struct {
	Name       string
	State      string
	PID        int
	Uptime     string
	MemoryMB   float64
	CPUPercent float64
}

type DeploymentScript interface {
	GenerateDeployScript(app *model.Application, version *model.Version) string
	GenerateRollbackScript(app *model.Application, targetVersion string) string
	GenerateHealthCheckScript(app *model.Application) string
}

type HostManager interface {
	GetHostInfo(host *model.Host) (*HostInfo, error)
	CheckDiskSpace(host *model.Host, requiredMB int64) (bool, error)
	CheckPort(host *model.Host, port int) (bool, error)
	BackupCurrentVersion(host *model.Host, app *model.Application) error
}

type HostInfo struct {
	Hostname      string
	OS            string
	Architecture  string
	TotalMemoryMB int64
	FreeDiskGB    float64
	CPUCores      int
	Uptime        string
}

type FileManager interface {
	CreateDirectory(path string, perm int) error
	CopyFile(src, dst string) error
	RemoveFile(path string) error
	ExtractArchive(archivePath, destPath string) error
	SetPermissions(path string, perm int) error
}

type LogCollector interface {
	TailLog(logPath string, lines int) ([]string, error)
	StreamLog(logPath string, callback func(line string)) error
	RotateLog(logPath string) error
}

package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config Agent 配置
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
	Agent  AgentConfig  `mapstructure:"agent"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// AgentConfig Agent配置
type AgentConfig struct {
	ID       string            `mapstructure:"id"`
	Hostname string            `mapstructure:"hostname"`
	WorkDir  string            `mapstructure:"work_dir"`
	Docker   DockerConfig      `mapstructure:"docker"`
	Health   HealthConfig      `mapstructure:"health"`
	Config   map[string]string `mapstructure:"config"`
}

// DockerConfig Docker配置
type DockerConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	SocketPath  string `mapstructure:"socket_path"`
	Registry    string `mapstructure:"registry"`
	NetworkMode string `mapstructure:"network_mode"`
}

// HealthConfig 健康检查配置
type HealthConfig struct {
	CheckInterval int `mapstructure:"check_interval"` // 秒
	Timeout       int `mapstructure:"timeout"`        // 秒
	RetryCount    int `mapstructure:"retry_count"`
}

// Load 加载配置
func Load(configPath string) (*Config, error) {
	// 如果指定了配置文件路径，使用指定的路径
	if configPath != "" {
		return loadConfigFromFile(configPath)
	}

	// 否则尝试从默认位置加载
	return loadConfigFromDefault()
}

// loadConfigFromFile 从指定文件加载配置
func loadConfigFromFile(configPath string) (*Config, error) {
	// 使用 viper 加载指定配置文件
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 从环境变量覆盖配置
	overrideFromEnv(&cfg)

	return &cfg, nil
}

// loadConfigFromDefault 从默认位置加载配置
func loadConfigFromDefault() (*Config, error) {
	// 设置配置文件搜索路径
	viper.SetConfigName("agent")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/operator-pm-agent/configs")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/boreas-agent")
	viper.AddConfigPath("$HOME/.boreas-agent")

	// 设置默认值
	setDefaults()

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// 配置文件不存在，使用默认配置
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 从环境变量覆盖配置
	overrideFromEnv(&cfg)

	return &cfg, nil
}

// setDefaults 设置默认值
func setDefaults() {
	// 服务器配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8081)

	// 日志配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")

	// Agent 配置
	viper.SetDefault("agent.id", "")
	viper.SetDefault("agent.hostname", "")
	viper.SetDefault("agent.work_dir", "/var/lib/boreas-agent")
	viper.SetDefault("agent.docker.enabled", true)
	viper.SetDefault("agent.docker.socket_path", "/var/run/docker.sock")
	viper.SetDefault("agent.docker.registry", "")
	viper.SetDefault("agent.docker.network_mode", "bridge")
	viper.SetDefault("agent.health.check_interval", 30)
	viper.SetDefault("agent.health.timeout", 10)
	viper.SetDefault("agent.health.retry_count", 3)
}

// overrideFromEnv 从环境变量覆盖配置
func overrideFromEnv(cfg *Config) {
	// 服务器配置
	if host := os.Getenv("AGENT_SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("AGENT_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}

	// 日志配置
	if level := os.Getenv("AGENT_LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}
	if format := os.Getenv("AGENT_LOG_FORMAT"); format != "" {
		cfg.Log.Format = format
	}
	if output := os.Getenv("AGENT_LOG_OUTPUT"); output != "" {
		cfg.Log.Output = output
	}

	// Agent 配置
	if id := os.Getenv("AGENT_ID"); id != "" {
		cfg.Agent.ID = id
	}
	if hostname := os.Getenv("AGENT_HOSTNAME"); hostname != "" {
		cfg.Agent.Hostname = hostname
	}
	if workDir := os.Getenv("AGENT_WORK_DIR"); workDir != "" {
		cfg.Agent.WorkDir = workDir
	}
	if enabled := os.Getenv("AGENT_DOCKER_ENABLED"); enabled != "" {
		cfg.Agent.Docker.Enabled = enabled == "true"
	}
	if socketPath := os.Getenv("AGENT_DOCKER_SOCKET_PATH"); socketPath != "" {
		cfg.Agent.Docker.SocketPath = socketPath
	}
	if registry := os.Getenv("AGENT_DOCKER_REGISTRY"); registry != "" {
		cfg.Agent.Docker.Registry = registry
	}
	if networkMode := os.Getenv("AGENT_DOCKER_NETWORK_MODE"); networkMode != "" {
		cfg.Agent.Docker.NetworkMode = networkMode
	}
	if checkInterval := os.Getenv("AGENT_HEALTH_CHECK_INTERVAL"); checkInterval != "" {
		if interval, err := strconv.Atoi(checkInterval); err == nil {
			cfg.Agent.Health.CheckInterval = interval
		}
	}
	if timeout := os.Getenv("AGENT_HEALTH_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.Agent.Health.Timeout = t
		}
	}
	if retryCount := os.Getenv("AGENT_HEALTH_RETRY_COUNT"); retryCount != "" {
		if retry, err := strconv.Atoi(retryCount); err == nil {
			cfg.Agent.Health.RetryCount = retry
		}
	}
}

// GetServerAddr 获取服务器地址
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetAgentWorkDir 获取Agent工作目录
func (c *Config) GetAgentWorkDir() string {
	if c.Agent.WorkDir != "" {
		return c.Agent.WorkDir
	}
	return "/var/lib/boreas-agent"
}

// GetAgentID 获取Agent ID
func (c *Config) GetAgentID() string {
	if c.Agent.ID != "" {
		return c.Agent.ID
	}
	// 如果没有配置ID，使用主机名
	if c.Agent.Hostname != "" {
		return c.Agent.Hostname
	}
	// 最后使用默认值
	return "agent-unknown"
}

// IsDockerEnabled 检查Docker是否启用
func (c *Config) IsDockerEnabled() bool {
	return c.Agent.Docker.Enabled
}

// GetDockerSocketPath 获取Docker Socket路径
func (c *Config) GetDockerSocketPath() string {
	if c.Agent.Docker.SocketPath != "" {
		return c.Agent.Docker.SocketPath
	}
	return "/var/run/docker.sock"
}

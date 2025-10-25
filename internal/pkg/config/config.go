package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	GitHub   GitHubConfig   `mapstructure:"github"`
	K8s      K8sConfig      `mapstructure:"k8s"`

	Trigger TriggerConfig `mapstructure:"trigger"`
	Agent   AgentConfig   `mapstructure:"agent"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// GitHubConfig GitHub配置
type GitHubConfig struct {
	WebhookSecret string `mapstructure:"webhook_secret"`
	Token         string `mapstructure:"token"`
}

// K8sConfig Kubernetes配置
type K8sConfig struct {
	ConfigPath string `mapstructure:"config_path"`
	Namespace  string `mapstructure:"namespace"`
}

type TriggerConfig struct {
	WebhookSecret  string `mapstructure:"webhook_secret"`
	WorkDir        string `mapstructure:"work_dir"`
	DockerRegistry string `mapstructure:"docker_registry"`
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
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

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
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 从环境变量覆盖配置
	overrideFromEnv(&config)

	return &config, nil
}

// setDefaults 设置默认值
func setDefaults() {
	// 服务器配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)

	// 数据库配置
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "boreas")
	viper.SetDefault("database.password", "boreas123")
	viper.SetDefault("database.name", "boreas")
	viper.SetDefault("database.ssl_mode", "disable")

	// Redis配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// 日志配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")

	// GitHub配置
	viper.SetDefault("github.webhook_secret", "")
	viper.SetDefault("github.token", "")

	// Kubernetes配置
	viper.SetDefault("k8s.config_path", "")
	viper.SetDefault("k8s.namespace", "default")

	// Agent配置
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
func overrideFromEnv(config *Config) {
	// 服务器配置
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}

	// 数据库配置
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Database.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		config.Database.Name = name
	}
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		config.Database.SSLMode = sslMode
	}

	// Redis配置
	if host := os.Getenv("REDIS_HOST"); host != "" {
		config.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Redis.Port = p
		}
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.Redis.Password = password
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		if d, err := strconv.Atoi(db); err == nil {
			config.Redis.DB = d
		}
	}

	// 日志配置
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Log.Level = level
	}

	// GitHub配置
	if secret := os.Getenv("GITHUB_WEBHOOK_SECRET"); secret != "" {
		config.GitHub.WebhookSecret = secret
	}
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		config.GitHub.Token = token
	}

	// Agent配置
	if id := os.Getenv("AGENT_ID"); id != "" {
		config.Agent.ID = id
	}
	if hostname := os.Getenv("AGENT_HOSTNAME"); hostname != "" {
		config.Agent.Hostname = hostname
	}
	if workDir := os.Getenv("AGENT_WORK_DIR"); workDir != "" {
		config.Agent.WorkDir = workDir
	}
	if enabled := os.Getenv("AGENT_DOCKER_ENABLED"); enabled != "" {
		config.Agent.Docker.Enabled = enabled == "true"
	}
	if socketPath := os.Getenv("AGENT_DOCKER_SOCKET_PATH"); socketPath != "" {
		config.Agent.Docker.SocketPath = socketPath
	}
	if registry := os.Getenv("AGENT_DOCKER_REGISTRY"); registry != "" {
		config.Agent.Docker.Registry = registry
	}
	if networkMode := os.Getenv("AGENT_DOCKER_NETWORK_MODE"); networkMode != "" {
		config.Agent.Docker.NetworkMode = networkMode
	}
	if checkInterval := os.Getenv("AGENT_HEALTH_CHECK_INTERVAL"); checkInterval != "" {
		if interval, err := strconv.Atoi(checkInterval); err == nil {
			config.Agent.Health.CheckInterval = interval
		}
	}
	if timeout := os.Getenv("AGENT_HEALTH_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			config.Agent.Health.Timeout = t
		}
	}
	if retryCount := os.Getenv("AGENT_HEALTH_RETRY_COUNT"); retryCount != "" {
		if retry, err := strconv.Atoi(retryCount); err == nil {
			config.Agent.Health.RetryCount = retry
		}
	}
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetRedisAddr 获取Redis地址
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
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

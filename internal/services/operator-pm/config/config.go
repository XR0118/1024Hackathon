package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config Operator-PM 配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	PM       PMConfig       `mapstructure:"pm"`
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

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// PMConfig 物理机管理配置
type PMConfig struct {
	AgentTimeout int               `mapstructure:"agent_timeout"` // Agent 通信超时（秒）
	MaxRetries   int               `mapstructure:"max_retries"`   // 最大重试次数
	HealthCheck  HealthCheckConfig `mapstructure:"health_check"`  // 健康检查配置
	Deployment   DeploymentConfig  `mapstructure:"deployment"`    // 部署配置
	Config       map[string]string `mapstructure:"config"`        // 自定义配置
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Interval int `mapstructure:"interval"` // 检查间隔（秒）
	Timeout  int `mapstructure:"timeout"`  // 超时时间（秒）
}

// DeploymentConfig 部署配置
type DeploymentConfig struct {
	Timeout       int `mapstructure:"timeout"`        // 部署超时（秒）
	MaxConcurrent int `mapstructure:"max_concurrent"` // 最大并发部署数
	RetryInterval int `mapstructure:"retry_interval"` // 重试间隔（秒）
	StatusCheck   int `mapstructure:"status_check"`   // 状态检查间隔（秒）
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
	viper.SetConfigName("operator-pm")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/operator-pm/configs")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/boreas-operator-pm")
	viper.AddConfigPath("$HOME/.boreas-operator-pm")

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
	viper.SetDefault("server.port", 8080)

	// 数据库配置
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "boreas")
	viper.SetDefault("database.password", "boreas123")
	viper.SetDefault("database.name", "boreas")
	viper.SetDefault("database.ssl_mode", "disable")

	// 日志配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")

	// PM 配置
	viper.SetDefault("pm.agent_timeout", 30)
	viper.SetDefault("pm.max_retries", 3)
	viper.SetDefault("pm.health_check.interval", 60)
	viper.SetDefault("pm.health_check.timeout", 10)
	viper.SetDefault("pm.deployment.timeout", 300)
	viper.SetDefault("pm.deployment.max_concurrent", 5)
	viper.SetDefault("pm.deployment.retry_interval", 30)
	viper.SetDefault("pm.deployment.status_check", 30)
}

// overrideFromEnv 从环境变量覆盖配置
func overrideFromEnv(cfg *Config) {
	// 服务器配置
	if host := os.Getenv("PM_SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("PM_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}

	// 数据库配置
	if host := os.Getenv("PM_DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("PM_DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Database.Port = p
		}
	}
	if user := os.Getenv("PM_DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("PM_DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if name := os.Getenv("PM_DB_NAME"); name != "" {
		cfg.Database.Name = name
	}
	if sslMode := os.Getenv("PM_DB_SSL_MODE"); sslMode != "" {
		cfg.Database.SSLMode = sslMode
	}

	// 日志配置
	if level := os.Getenv("PM_LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}
	if format := os.Getenv("PM_LOG_FORMAT"); format != "" {
		cfg.Log.Format = format
	}
	if output := os.Getenv("PM_LOG_OUTPUT"); output != "" {
		cfg.Log.Output = output
	}

	// PM 配置
	if timeout := os.Getenv("PM_AGENT_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.PM.AgentTimeout = t
		}
	}
	if retries := os.Getenv("PM_MAX_RETRIES"); retries != "" {
		if r, err := strconv.Atoi(retries); err == nil {
			cfg.PM.MaxRetries = r
		}
	}
	if interval := os.Getenv("PM_HEALTH_CHECK_INTERVAL"); interval != "" {
		if i, err := strconv.Atoi(interval); err == nil {
			cfg.PM.HealthCheck.Interval = i
		}
	}
	if timeout := os.Getenv("PM_HEALTH_CHECK_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.PM.HealthCheck.Timeout = t
		}
	}
	if timeout := os.Getenv("PM_DEPLOYMENT_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.PM.Deployment.Timeout = t
		}
	}
	if concurrent := os.Getenv("PM_DEPLOYMENT_MAX_CONCURRENT"); concurrent != "" {
		if c, err := strconv.Atoi(concurrent); err == nil {
			cfg.PM.Deployment.MaxConcurrent = c
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

// GetServerAddr 获取服务器地址
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

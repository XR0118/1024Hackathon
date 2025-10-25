package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config Operator-K8s 配置
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
	K8s    K8sConfig    `mapstructure:"k8s"`
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

// K8sConfig Kubernetes配置
type K8sConfig struct {
	ConfigPath  string            `mapstructure:"config_path"`
	Namespace   string            `mapstructure:"namespace"`
	Context     string            `mapstructure:"context"`
	Timeout     int               `mapstructure:"timeout"`      // 操作超时（秒）
	RetryCount  int               `mapstructure:"retry_count"`  // 重试次数
	HealthCheck HealthCheckConfig `mapstructure:"health_check"` // 健康检查配置
	Deployment  DeploymentConfig  `mapstructure:"deployment"`   // 部署配置
	Config      map[string]string `mapstructure:"config"`       // 自定义配置
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Interval int `mapstructure:"interval"` // 检查间隔（秒）
	Timeout  int `mapstructure:"timeout"`  // 超时时间（秒）
}

// DeploymentConfig 部署配置
type DeploymentConfig struct {
	Timeout         int `mapstructure:"timeout"`          // 部署超时（秒）
	MaxConcurrent   int `mapstructure:"max_concurrent"`   // 最大并发部署数
	RetryInterval   int `mapstructure:"retry_interval"`   // 重试间隔（秒）
	StatusCheck     int `mapstructure:"status_check"`     // 状态检查间隔（秒）
	RollbackTimeout int `mapstructure:"rollback_timeout"` // 回滚超时（秒）
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
	viper.SetConfigName("operator-k8s")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/operator-k8s/configs")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/boreas-operator-k8s")
	viper.AddConfigPath("$HOME/.boreas-operator-k8s")

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
	viper.SetDefault("server.port", 8082)

	// 日志配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")

	// K8s 配置
	viper.SetDefault("k8s.config_path", "")
	viper.SetDefault("k8s.namespace", "default")
	viper.SetDefault("k8s.context", "")
	viper.SetDefault("k8s.timeout", 30)
	viper.SetDefault("k8s.retry_count", 3)
	viper.SetDefault("k8s.health_check.interval", 60)
	viper.SetDefault("k8s.health_check.timeout", 10)
	viper.SetDefault("k8s.deployment.timeout", 300)
	viper.SetDefault("k8s.deployment.max_concurrent", 5)
	viper.SetDefault("k8s.deployment.retry_interval", 30)
	viper.SetDefault("k8s.deployment.status_check", 30)
	viper.SetDefault("k8s.deployment.rollback_timeout", 180)
}

// overrideFromEnv 从环境变量覆盖配置
func overrideFromEnv(cfg *Config) {
	// 服务器配置
	if host := os.Getenv("K8S_SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("K8S_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}

	// 日志配置
	if level := os.Getenv("K8S_LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}
	if format := os.Getenv("K8S_LOG_FORMAT"); format != "" {
		cfg.Log.Format = format
	}
	if output := os.Getenv("K8S_LOG_OUTPUT"); output != "" {
		cfg.Log.Output = output
	}

	// K8s 配置
	if configPath := os.Getenv("K8S_CONFIG_PATH"); configPath != "" {
		cfg.K8s.ConfigPath = configPath
	}
	if namespace := os.Getenv("K8S_NAMESPACE"); namespace != "" {
		cfg.K8s.Namespace = namespace
	}
	if context := os.Getenv("K8S_CONTEXT"); context != "" {
		cfg.K8s.Context = context
	}
	if timeout := os.Getenv("K8S_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.K8s.Timeout = t
		}
	}
	if retries := os.Getenv("K8S_RETRY_COUNT"); retries != "" {
		if r, err := strconv.Atoi(retries); err == nil {
			cfg.K8s.RetryCount = r
		}
	}
	if interval := os.Getenv("K8S_HEALTH_CHECK_INTERVAL"); interval != "" {
		if i, err := strconv.Atoi(interval); err == nil {
			cfg.K8s.HealthCheck.Interval = i
		}
	}
	if timeout := os.Getenv("K8S_HEALTH_CHECK_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.K8s.HealthCheck.Timeout = t
		}
	}
	if timeout := os.Getenv("K8S_DEPLOYMENT_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.K8s.Deployment.Timeout = t
		}
	}
	if concurrent := os.Getenv("K8S_DEPLOYMENT_MAX_CONCURRENT"); concurrent != "" {
		if c, err := strconv.Atoi(concurrent); err == nil {
			cfg.K8s.Deployment.MaxConcurrent = c
		}
	}
}

// GetServerAddr 获取服务器地址
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

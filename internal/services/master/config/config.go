package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config Master 服务配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	GitHub   GitHubConfig   `mapstructure:"github"`
	K8s      K8sConfig      `mapstructure:"k8s"`
	Trigger  TriggerConfig  `mapstructure:"trigger"`
	Operator OperatorConfig `mapstructure:"operator"`
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

// TriggerConfig 触发器配置
type TriggerConfig struct {
	WebhookSecret  string `mapstructure:"webhook_secret"`
	WorkDir        string `mapstructure:"work_dir"`
	DockerRegistry string `mapstructure:"docker_registry"`
}

// OperatorConfig Operator 配置
type OperatorConfig struct {
	K8SOperatorURL string `mapstructure:"k8s_operator_url"`
	PMOperatorURL  string `mapstructure:"pm_operator_url"`
	UseMock        bool   `mapstructure:"use_mock"`
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
	viper.SetConfigName("master")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/master-service/configs")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/boreas-master")
	viper.AddConfigPath("$HOME/.boreas-master")

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

	// 触发器配置
	viper.SetDefault("trigger.webhook_secret", "")
	viper.SetDefault("trigger.work_dir", "/var/lib/boreas")
	viper.SetDefault("trigger.docker_registry", "")

	// Operator配置
	viper.SetDefault("operator.k8s_operator_url", "http://localhost:8081")
	viper.SetDefault("operator.pm_operator_url", "http://localhost:8082")
	viper.SetDefault("operator.use_mock", false)
}

// overrideFromEnv 从环境变量覆盖配置
func overrideFromEnv(cfg *Config) {
	// 服务器配置
	if host := os.Getenv("MASTER_SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("MASTER_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}

	// 数据库配置
	if host := os.Getenv("MASTER_DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("MASTER_DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Database.Port = p
		}
	}
	if user := os.Getenv("MASTER_DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("MASTER_DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if name := os.Getenv("MASTER_DB_NAME"); name != "" {
		cfg.Database.Name = name
	}
	if sslMode := os.Getenv("MASTER_DB_SSL_MODE"); sslMode != "" {
		cfg.Database.SSLMode = sslMode
	}

	// Redis配置
	if host := os.Getenv("MASTER_REDIS_HOST"); host != "" {
		cfg.Redis.Host = host
	}
	if port := os.Getenv("MASTER_REDIS_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Redis.Port = p
		}
	}
	if password := os.Getenv("MASTER_REDIS_PASSWORD"); password != "" {
		cfg.Redis.Password = password
	}
	if db := os.Getenv("MASTER_REDIS_DB"); db != "" {
		if d, err := strconv.Atoi(db); err == nil {
			cfg.Redis.DB = d
		}
	}

	// 日志配置
	if level := os.Getenv("MASTER_LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}
	if format := os.Getenv("MASTER_LOG_FORMAT"); format != "" {
		cfg.Log.Format = format
	}
	if output := os.Getenv("MASTER_LOG_OUTPUT"); output != "" {
		cfg.Log.Output = output
	}

	// GitHub配置
	if secret := os.Getenv("MASTER_GITHUB_WEBHOOK_SECRET"); secret != "" {
		cfg.GitHub.WebhookSecret = secret
	}
	if token := os.Getenv("MASTER_GITHUB_TOKEN"); token != "" {
		cfg.GitHub.Token = token
	}

	// Kubernetes配置
	if configPath := os.Getenv("MASTER_K8S_CONFIG_PATH"); configPath != "" {
		cfg.K8s.ConfigPath = configPath
	}
	if namespace := os.Getenv("MASTER_K8S_NAMESPACE"); namespace != "" {
		cfg.K8s.Namespace = namespace
	}

	// 触发器配置
	if secret := os.Getenv("MASTER_TRIGGER_WEBHOOK_SECRET"); secret != "" {
		cfg.Trigger.WebhookSecret = secret
	}
	if workDir := os.Getenv("MASTER_TRIGGER_WORK_DIR"); workDir != "" {
		cfg.Trigger.WorkDir = workDir
	}
	if registry := os.Getenv("MASTER_TRIGGER_DOCKER_REGISTRY"); registry != "" {
		cfg.Trigger.DockerRegistry = registry
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

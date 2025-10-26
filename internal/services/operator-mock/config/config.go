package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

func Load(configPath string) (*Config, error) {
	if configPath != "" {
		return loadConfigFromFile(configPath)
	}

	return loadConfigFromDefault()
}

func loadConfigFromFile(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	overrideFromEnv(&cfg)

	return &cfg, nil
}

func loadConfigFromDefault() (*Config, error) {
	viper.SetConfigName("operator-mock")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/operator-mock/configs")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/boreas-operator-mock")
	viper.AddConfigPath("$HOME/.boreas-operator-mock")

	setDefaults()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	overrideFromEnv(&cfg)

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8082)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
}

func overrideFromEnv(cfg *Config) {
	if host := os.Getenv("MOCK_SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("MOCK_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}

	if level := os.Getenv("MOCK_LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}
	if format := os.Getenv("MOCK_LOG_FORMAT"); format != "" {
		cfg.Log.Format = format
	}
	if output := os.Getenv("MOCK_LOG_OUTPUT"); output != "" {
		cfg.Log.Output = output
	}
}

func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

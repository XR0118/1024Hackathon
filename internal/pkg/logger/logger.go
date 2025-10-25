package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// Init 初始化日志器
func Init(level, format string) error {
	var config zap.Config

	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// 设置日志级别
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// 设置输出
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	// 创建日志器
	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	// 替换全局日志器
	zap.ReplaceGlobals(Logger)

	return nil
}

// GetLogger 获取日志器
func GetLogger() *zap.Logger {
	if Logger == nil {
		// 如果未初始化，使用默认配置
		_ = Init("info", "json")
	}
	return Logger
}

// Sync 同步日志
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// NewFileLogger 创建文件日志器
func NewFileLogger(filename string, level zapcore.Level) (*zap.Logger, error) {
	_, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.OutputPaths = []string{filename}
	config.ErrorOutputPaths = []string{filename}

	return config.Build()
}

package database

// DatabaseConfig 数据库配置接口
type DatabaseConfig interface {
	GetDSN() string
}

// InitWithConfig 使用配置接口初始化数据库
func InitWithConfig(cfg DatabaseConfig) error {
	dsn := cfg.GetDSN()
	return InitWithDSN(dsn)
}

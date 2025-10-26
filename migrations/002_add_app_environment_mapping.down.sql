-- 删除索引
DROP INDEX IF EXISTS idx_app_env_environment_id;
DROP INDEX IF EXISTS idx_app_env_application_id;

-- 删除应用-环境映射表
DROP TABLE IF EXISTS application_environments;


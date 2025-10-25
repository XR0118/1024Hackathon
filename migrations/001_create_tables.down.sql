-- 删除索引
DROP INDEX IF EXISTS idx_tasks_type;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_app_id;
DROP INDEX IF EXISTS idx_tasks_deployment_id;
DROP INDEX IF EXISTS idx_deployments_created_at;
DROP INDEX IF EXISTS idx_deployments_version_id;
DROP INDEX IF EXISTS idx_deployments_environment_id;
DROP INDEX IF EXISTS idx_deployments_status;
DROP INDEX IF EXISTS idx_environments_is_active;
DROP INDEX IF EXISTS idx_environments_type;
DROP INDEX IF EXISTS idx_applications_type;
DROP INDEX IF EXISTS idx_applications_repository;
DROP INDEX IF EXISTS idx_versions_app_builds_gin;
DROP INDEX IF EXISTS idx_versions_created_at;
DROP INDEX IF EXISTS idx_versions_repository;

-- 删除表
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS deployments;
DROP TABLE IF EXISTS environments;
DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS versions;

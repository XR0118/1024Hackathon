-- 创建应用-环境映射表
CREATE TABLE IF NOT EXISTS application_environments (
    id VARCHAR(255) PRIMARY KEY,
    application_id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(application_id, environment_id),
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE
);

-- 创建索引以优化查询
CREATE INDEX idx_app_env_application_id ON application_environments(application_id);
CREATE INDEX idx_app_env_environment_id ON application_environments(environment_id);

-- 添加注释
COMMENT ON TABLE application_environments IS '应用-环境映射表，表示某个应用会在哪些环境中有部署';
COMMENT ON COLUMN application_environments.application_id IS '应用ID';
COMMENT ON COLUMN application_environments.environment_id IS '环境ID';


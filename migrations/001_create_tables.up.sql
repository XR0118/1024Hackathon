-- 创建版本表
CREATE TABLE IF NOT EXISTS versions (
    id VARCHAR(255) PRIMARY KEY,
    git_tag VARCHAR(255) UNIQUE NOT NULL,
    git_commit VARCHAR(255) NOT NULL,
    repository VARCHAR(500) NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    description TEXT,
    app_builds JSONB
);

-- 创建应用表
CREATE TABLE IF NOT EXISTS applications (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    repository VARCHAR(500) NOT NULL,
    type VARCHAR(100) NOT NULL,
    config JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 创建环境表
CREATE TABLE IF NOT EXISTS environments (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    type VARCHAR(100) NOT NULL,
    config JSONB NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 创建部署表
CREATE TABLE IF NOT EXISTS deployments (
    id VARCHAR(255) PRIMARY KEY,
    version_id VARCHAR(255) NOT NULL REFERENCES versions(id),
    must_in_order JSONB,
    environment_id VARCHAR(255) NOT NULL REFERENCES environments(id),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,

    manual_approval BOOLEAN NOT NULL DEFAULT FALSE,
    strategy JSONB
);

-- 创建任务表
CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(255) PRIMARY KEY,
    deployment_id VARCHAR(255) NOT NULL REFERENCES deployments(id),
    app_id VARCHAR(255) NOT NULL REFERENCES applications(id),
    type VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    block_by VARCHAR(255),
    payload TEXT,
    result TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_versions_repository ON versions(repository);
CREATE INDEX IF NOT EXISTS idx_versions_created_at ON versions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_versions_app_builds_gin ON versions USING gin (app_builds jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_applications_repository ON applications(repository);
CREATE INDEX IF NOT EXISTS idx_applications_type ON applications(type);
CREATE INDEX IF NOT EXISTS idx_environments_type ON environments(type);
CREATE INDEX IF NOT EXISTS idx_environments_is_active ON environments(is_active);
CREATE INDEX IF NOT EXISTS idx_deployments_status ON deployments(status);
CREATE INDEX IF NOT EXISTS idx_deployments_environment_id ON deployments(environment_id);
CREATE INDEX IF NOT EXISTS idx_deployments_version_id ON deployments(version_id);
CREATE INDEX IF NOT EXISTS idx_deployments_created_at ON deployments(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_deployment_id ON tasks(deployment_id);
CREATE INDEX IF NOT EXISTS idx_tasks_app_id ON tasks(app_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_type ON tasks(type);

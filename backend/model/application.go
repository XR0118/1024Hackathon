package model

import "time"

type Application struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Repository  string    `json:"repository"`
	Description string    `json:"description"`
	Owner       string    `json:"owner"`
	Config      AppConfig `json:"config"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AppConfig struct {
	BuildConfig   BuildConfig   `json:"build_config"`
	RuntimeConfig RuntimeConfig `json:"runtime_config"`
	HealthCheck   HealthCheck   `json:"health_check"`
}

type BuildConfig struct {
	Dockerfile string            `json:"dockerfile"`
	BuildArgs  map[string]string `json:"build_args"`
	Context    string            `json:"context"`
}

type RuntimeConfig struct {
	Port      int               `json:"port"`
	Env       map[string]string `json:"env"`
	Resources Resources         `json:"resources"`
}

type Resources struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type HealthCheck struct {
	Path         string `json:"path"`
	Port         int    `json:"port"`
	InitialDelay int    `json:"initial_delay"`
	Period       int    `json:"period"`
}

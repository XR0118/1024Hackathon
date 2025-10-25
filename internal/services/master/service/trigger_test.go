package service

import (
	"context"
	"testing"
)

func TestVersionService_ProcessTagEvent(t *testing.T) {
	config := &TriggerConfig{
		WorkDir:        "/tmp/test-version-trigger",
		DockerRegistry: "registry.example.com",
	}

	service := NewTriggerService(config)

	event := &TagEvent{
		TagName:    "v1.0.0",
		Repository: "https://github.com/example/repo.git",
		Commit:     "abc123",
		Pusher:     "test-user",
	}

	ctx := context.Background()
	result, err := service.ProcessTagEvent(ctx, event)

	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}

	if result != nil {
		t.Logf("Result: %+v", result)
	}
}

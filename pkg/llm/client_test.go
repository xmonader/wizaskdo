package llm

import (
	"os"
	"testing"
)

func TestNewClientMissingAPIKey(t *testing.T) {
	// Save and restore env
	original := os.Getenv("OPENROUTER_API_KEY")
	defer os.Setenv("OPENROUTER_API_KEY", original)

	os.Unsetenv("OPENROUTER_API_KEY")
	client, err := NewClient()
	if err == nil {
		t.Error("expected error for missing API key")
	}
	if client != nil {
		t.Error("expected nil client for missing API key")
	}
}

func TestNewClientWithDefaults(t *testing.T) {
	originalKey := os.Getenv("OPENROUTER_API_KEY")
	originalModel := os.Getenv("WIZASK_MODEL")
	defer func() {
		if originalKey == "" {
			os.Unsetenv("OPENROUTER_API_KEY")
		} else {
			os.Setenv("OPENROUTER_API_KEY", originalKey)
		}
		if originalModel == "" {
			os.Unsetenv("WIZASK_MODEL")
		} else {
			os.Setenv("WIZASK_MODEL", originalModel)
		}
	}()

	os.Setenv("OPENROUTER_API_KEY", "test-key")
	os.Unsetenv("WIZASK_MODEL")

	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.APIKey != "test-key" {
		t.Errorf("expected API key 'test-key', got '%s'", client.APIKey)
	}
	if client.Model != DefaultModel {
		t.Errorf("expected default model '%s', got '%s'", DefaultModel, client.Model)
	}
}

func TestNewClientWithCustomModel(t *testing.T) {
	originalKey := os.Getenv("OPENROUTER_API_KEY")
	originalModel := os.Getenv("WIZASK_MODEL")
	defer func() {
		if originalKey == "" {
			os.Unsetenv("OPENROUTER_API_KEY")
		} else {
			os.Setenv("OPENROUTER_API_KEY", originalKey)
		}
		if originalModel == "" {
			os.Unsetenv("WIZASK_MODEL")
		} else {
			os.Setenv("WIZASK_MODEL", originalModel)
		}
	}()

	os.Setenv("OPENROUTER_API_KEY", "test-key")
	os.Setenv("WIZASK_MODEL", "custom-model")

	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Model != "custom-model" {
		t.Errorf("expected model 'custom-model', got '%s'", client.Model)
	}
}

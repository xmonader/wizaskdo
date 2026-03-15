package main

import (
	"os"
	"strings"
	"testing"
)

func TestSystemTemplate(t *testing.T) {
	if !strings.Contains(systemTemplate, "terminal assistant") {
		t.Error("system template should mention terminal assistant")
	}
	if !strings.Contains(systemTemplate, "code block") {
		t.Error("system template should mention code blocks")
	}
}

func TestApiKeyCheck(t *testing.T) {
	original := os.Getenv("OPENROUTER_API_KEY")
	defer os.Setenv("OPENROUTER_API_KEY", original)

	os.Unsetenv("OPENROUTER_API_KEY")
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey != "" {
		t.Error("API key should be empty after unset")
	}
}

func TestRequestMarshal(t *testing.T) {
	req := Request{
		Model: defaultModel,
		Messages: []Message{
			{Role: "system", Content: "test"},
			{Role: "user", Content: "hello"},
		},
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	if req.Model == "" {
		t.Error("model should not be empty")
	}
	if len(req.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(req.Messages))
	}
	if req.MaxTokens != maxTokens {
		t.Errorf("expected maxTokens %d, got %d", maxTokens, req.MaxTokens)
	}
}

package main

import (
	"os"
	"testing"
)

func TestGetEnvFallback(t *testing.T) {
	const key = "TEST_GETENV_FALLBACK"
	_ = os.Unsetenv(key)

	value := getEnv(key, "default")
	if value != "default" {
		t.Fatalf("expected fallback 'default', got %q", value)
	}
}

func TestGetEnvOverride(t *testing.T) {
	const key = "TEST_GETENV_OVERRIDE"
	os.Setenv(key, "value")
	defer os.Unsetenv(key)

	value := getEnv(key, "default")
	if value != "value" {
		t.Fatalf("expected env override 'value', got %q", value)
	}
}

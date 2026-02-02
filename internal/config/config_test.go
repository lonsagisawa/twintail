package config

import (
	"os"
	"testing"
)

func TestLoad_DefaultPort(t *testing.T) {
	os.Unsetenv("PORT")

	cfg := Load()

	if cfg.Port != "8077" {
		t.Errorf("expected default port '8077', got '%s'", cfg.Port)
	}
}

func TestLoad_CustomPort(t *testing.T) {
	os.Setenv("PORT", "9000")
	defer os.Unsetenv("PORT")

	cfg := Load()

	if cfg.Port != "9000" {
		t.Errorf("expected port '9000', got '%s'", cfg.Port)
	}
}

func TestLoad_EmptyPortFallsBackToDefault(t *testing.T) {
	os.Setenv("PORT", "")
	defer os.Unsetenv("PORT")

	cfg := Load()

	if cfg.Port != "8077" {
		t.Errorf("expected default port '8077', got '%s'", cfg.Port)
	}
}

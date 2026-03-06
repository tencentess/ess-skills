package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetConfigDir_Default(t *testing.T) {
	os.Unsetenv("TSIGN_CONFIG_PATH")
	dir := GetConfigDir()
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		up := os.Getenv("USERPROFILE")
		if up == "" {
			up = home
		}
		expected := filepath.Join(up, ".tsign")
		if dir != expected {
			t.Errorf("expected %s, got %s", expected, dir)
		}
	} else {
		expected := filepath.Join(home, ".tsign")
		if dir != expected {
			t.Errorf("expected %s, got %s", expected, dir)
		}
	}
}

func TestGetConfigDir_EnvOverride(t *testing.T) {
	t.Setenv("TSIGN_CONFIG_PATH", "/tmp/custom/.tsign/config.yaml")
	dir := GetConfigDir()
	if dir != "/tmp/custom/.tsign" {
		t.Errorf("expected /tmp/custom/.tsign, got %s", dir)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	tmpPath := filepath.Join(tmpDir, "config.yaml")
	t.Setenv("TSIGN_CONFIG_PATH", tmpPath)

	cfg := &Config{
		Credentials: Credentials{SecretID: "test-id", SecretKey: "test-key"},
		Operator:    Operator{UserID: "test-user"},
		Env:         "online",
	}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load("")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Credentials.SecretID != "test-id" {
		t.Errorf("expected test-id, got %s", loaded.Credentials.SecretID)
	}
	if loaded.Operator.UserID != "test-user" {
		t.Errorf("expected test-user, got %s", loaded.Operator.UserID)
	}
	if loaded.Env != "online" {
		t.Errorf("expected online, got %s", loaded.Env)
	}
}

func TestLoadProfile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpPath := filepath.Join(tmpDir, "config.yaml")
	t.Setenv("TSIGN_CONFIG_PATH", tmpPath)

	cfg := &Config{
		Credentials: Credentials{SecretID: "default-id", SecretKey: "default-key"},
		Operator:    Operator{UserID: "default-user"},
		Env:         "online",
		Profiles: map[string]Profile{
			"staging": {
				Credentials: Credentials{SecretID: "staging-id", SecretKey: "staging-key"},
				Operator:    Operator{UserID: "staging-user"},
				Env:         "test",
			},
		},
	}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load("staging")
	if err != nil {
		t.Fatalf("Load staging failed: %v", err)
	}
	if loaded.Credentials.SecretID != "staging-id" {
		t.Errorf("expected staging-id, got %s", loaded.Credentials.SecretID)
	}
	if loaded.Env != "test" {
		t.Errorf("expected test, got %s", loaded.Env)
	}
}

func TestLoadProfile_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	tmpPath := filepath.Join(tmpDir, "config.yaml")
	t.Setenv("TSIGN_CONFIG_PATH", tmpPath)

	cfg := &Config{
		Credentials: Credentials{SecretID: "id", SecretKey: "key"},
		Operator:    Operator{UserID: "user"},
	}
	Save(cfg)

	_, err := Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent profile")
	}
}

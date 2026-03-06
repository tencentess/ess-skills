package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Credentials 凭证配置
type Credentials struct {
	SecretID  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
}

// Operator 经办人配置
type Operator struct {
	UserID string `yaml:"user_id"`
}

// Profile 命名配置档
type Profile struct {
	Credentials Credentials `yaml:"credentials"`
	Operator    Operator    `yaml:"operator"`
	Env         string      `yaml:"env,omitempty"`
}

// Config 完整配置文件结构
type Config struct {
	Credentials Credentials        `yaml:"credentials"`
	Operator    Operator           `yaml:"operator"`
	Env         string             `yaml:"env,omitempty"`
	Profiles    map[string]Profile `yaml:"profiles,omitempty"`
}

// GetConfigDir 返回跨平台配置目录
func GetConfigDir() string {
	if p := os.Getenv("TSIGN_CONFIG_PATH"); p != "" {
		return filepath.Dir(p)
	}
	switch runtime.GOOS {
	case "windows":
		if up := os.Getenv("USERPROFILE"); up != "" {
			return filepath.Join(up, ".tsign")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".tsign")
	default:
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".tsign")
	}
}

// GetConfigPath 返回跨平台配置文件路径
func GetConfigPath() string {
	if p := os.Getenv("TSIGN_CONFIG_PATH"); p != "" {
		return p
	}
	return filepath.Join(GetConfigDir(), "config.yaml")
}

// Load 从文件加载配置，可指定 profile
func Load(profile string) (*Config, error) {
	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败 %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	if profile != "" {
		if p, ok := cfg.Profiles[profile]; ok {
			cfg.Credentials = p.Credentials
			cfg.Operator = p.Operator
			if p.Env != "" {
				cfg.Env = p.Env
			}
		} else {
			return nil, fmt.Errorf("profile '%s' 不存在", profile)
		}
	}
	return &cfg, nil
}

// Save 保存配置到文件（自动设置权限）
func Save(cfg *Config) error {
	path := GetConfigPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	if runtime.GOOS != "windows" {
		return os.WriteFile(path, data, 0600)
	}
	return os.WriteFile(path, data, 0644)
}

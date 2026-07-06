package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type JWTConfig struct {
	Secret        string `yaml:"secret"`
	ExpireSeconds int64  `yaml:"expire_seconds"`
}

func Load(path string) (Config, error) {
	data, err := readConfigFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	applyDefaults(&cfg)
	if err := applyEnvOverrides(&cfg); err != nil {
		return Config{}, err
	}
	applyDefaults(&cfg)
	return cfg, nil
}

func readConfigFile(path string) ([]byte, error) {
	if envPath := strings.TrimSpace(os.Getenv("CONFIG_PATH")); envPath != "" {
		path = envPath
	}
	if strings.TrimSpace(path) == "" {
		path = "config.yaml"
	}

	data, err := os.ReadFile(path)
	if err == nil {
		return data, nil
	}
	if errors.Is(err, os.ErrNotExist) && path == "config.yaml" {
		return os.ReadFile("config.example.yaml")
	}
	return nil, err
}

func applyDefaults(cfg *Config) {
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.JWT.ExpireSeconds == 0 {
		cfg.JWT.ExpireSeconds = 86400
	}
}

func applyEnvOverrides(cfg *Config) error {
	if value := strings.TrimSpace(os.Getenv("SERVER_PORT")); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		cfg.Server.Port = port
	}

	if value := strings.TrimSpace(os.Getenv("DB_DSN")); value != "" {
		cfg.Database.DSN = value
	}

	if value := strings.TrimSpace(os.Getenv("JWT_SECRET")); value != "" {
		cfg.JWT.Secret = value
	}

	if value := strings.TrimSpace(os.Getenv("JWT_EXPIRE_SECONDS")); value != "" {
		expireSeconds, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		cfg.JWT.ExpireSeconds = expireSeconds
	}

	return nil
}

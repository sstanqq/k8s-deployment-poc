package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

var (
	ErrPathNotFound     = errors.New("path not found")
	ErrPathNotDirectory = errors.New("path is not a directory")
	ErrNegativeDuration = errors.New("duration must be positive")
	ErrInvalidHost      = errors.New("invalid host")
	ErrInvalidPort      = errors.New("invalid port")
)

type HostConfig struct {
	NodeName string `env:"NODE_NAME"`
	NodeIP   string `env:"NODE_IP"`
}

type Config struct {
	SrvHost string `env:"HTTP_ADDR_HOST" envDefault:"0.0.0.0"`
	SrvPort int    `env:"HTTP_ADDR_PORT" envDefault:"8000"`

	SrvName    string `env:"SERVER_NAME"`
	SrvVersion string `env:"SERVER_VERSION" envDefault:"1.0.0"`

	LogFilePath string `env:"LOG_FILE_PATH" envDefault:"./logs/"`

	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`

	HstConfig *HostConfig
}

func LoadConfig() (*Config, error) {
	envPath := os.Getenv("ENV_FILE_PATH")
	if envPath == "" {
		envPath = "configs/.env"
	}

	godotenv.Load(envPath)

	var conf Config
	if err := env.Parse(&conf); err != nil {
		return nil, err
	}

	var hostConf HostConfig
	if err := env.Parse(&hostConf); err != nil {
		return nil, fmt.Errorf("failed to parse host config: %w", err)
	}
	conf.HstConfig = &hostConf

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &conf, nil
}

func (c *Config) Validate() error {
	if err := validateHost(c.SrvHost); err != nil {
		return fmt.Errorf("failed to validate HTTP_ADDR_HOST: %w", err)
	}
	if err := validatePort(c.SrvPort); err != nil {
		return fmt.Errorf("failed to validate HTTP_ADDR_PORT: %w", err)
	}

	if err := validateOSPath(c.LogFilePath); err != nil {
		return fmt.Errorf("failed to validate LOG_FILE_PATH: %w", err)
	}

	if err := validateTimeout(c.ShutdownTimeout); err != nil {
		return fmt.Errorf("failed to validate SHUTDOWN_TIMEOUT: %w", err)
	}

	return nil
}

func validateOSPath(path string) error {
	if path == "" {
		return ErrPathNotFound
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrPathNotFound
		}
		return fmt.Errorf("failed to access path: %w", err)
	}

	if !info.IsDir() {
		return ErrPathNotDirectory
	}

	return nil
}

func validateTimeout(d time.Duration) error {
	if d <= 0 {
		return ErrNegativeDuration
	}

	return nil
}

func validateHost(host string) error {
	if host == "" {
		return ErrInvalidHost
	}

	if _, err := net.LookupHost(host); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidHost, err)
	}

	return nil
}

func validatePort(port int) error {
	if port <= 0 || port > 65535 {
		return ErrInvalidPort
	}
	return nil
}

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	version    string = "dev"
	commitHash string = "-"
)

type Project struct {
	Name       string `yaml:"name"`
	LogLevel   int    `yaml:"logLevel"`
	Version    string
	CommitHash string
}

type Server struct {
	Address           string        `yaml:"address"`
	Port              int           `yaml:"port"`
	ConnectionTimeout time.Duration `yaml:"connectionTimeout"`
}

type ProofOfWork struct {
	Difficulty int `yaml:"difficulty"`
}

type Config struct {
	Project     Project     `yaml:"project"`
	Server      Server      `yaml:"server"`
	ProofOfWork ProofOfWork `yaml:"proofOfWork"`
}

var cfg *Config

func ReadFile(path string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("open config file: %w", err)
	}
	defer func() {
		_ = file.Close() //nolint:errcheck
	}()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return fmt.Errorf("decode config file: %w", err)
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return nil
}

func GetInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

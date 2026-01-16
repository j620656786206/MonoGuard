// Package config provides configuration management using Viper
package config

import "github.com/spf13/viper"

// Config represents the MonoGuard configuration structure
type Config struct {
	Workspaces []string   `mapstructure:"workspaces" json:"workspaces"`
	Rules      Rules      `mapstructure:"rules" json:"rules"`
	Thresholds Thresholds `mapstructure:"thresholds" json:"thresholds"`
}

// Rules defines validation rules configuration
type Rules struct {
	CircularDependencies string `mapstructure:"circularDependencies" json:"circularDependencies"`
	BoundaryViolations   string `mapstructure:"boundaryViolations" json:"boundaryViolations"`
}

// Thresholds defines threshold configuration
type Thresholds struct {
	HealthScore int `mapstructure:"healthScore" json:"healthScore"`
}

// Load reads configuration from Viper
func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

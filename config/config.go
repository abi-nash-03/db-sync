package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Source      DatabaseConfig `yaml:"source"`
	Destination DatabaseConfig `yaml:"destination"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func LoadConfig(configPath string) (*Config, error) {

	// Read the YAML file content
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	fmt.Println("Config loaded successfully")
	fmt.Println(config)
	return &config, nil
}

// func Validate(c *Config) error {
// 	if c.Source.Host == "" {
// 		return fmt.Errorf("source host is required")
// 	}
// 	if c.Source.Port == 0 {
// 		return fmt.Errorf("source port is required")
// 	}
// 	if c.Source.User == "" {
// 		return fmt.Errorf("source user is required")
// 	}
// 	if c.Source.Password == "" {
// 		return fmt.Errorf("source password is required")
// 	}
// 	if c.Source.Database == "" {
// 		return fmt.Errorf("source database is required")
// 	}
// 	if c.Destination.Host == "" {
// 		return fmt.Errorf("destination host is required")
// 	}
// 	if c.Destination.Port == 0 {
// 		return fmt.Errorf("destination port is required")
// 	}
// 	if c.Destination.User == "" {
// 		return fmt.Errorf("destination user is required")
// 	}
// 	if c.Destination.Password == "" {
// 		return fmt.Errorf("destination password is required")
// 	}
// 	if c.Destination.Database == "" {
// 		return fmt.Errorf("destination database is required")
// 	}
// 	return nil
// }

func (c *Config) Validate() error {
	required := map[string]string{
		"source.host":     c.Source.Host,
		"source.user":     c.Source.User,
		"source.password": c.Source.Password,
		"source.database": c.Source.Database,
		"dest.host":       c.Destination.Host,
		"dest.user":       c.Destination.User,
		"dest.password":   c.Destination.Password,
		"dest.database":   c.Destination.Database,
	}

	for field, value := range required {
		if value == "" {
			return fmt.Errorf("missing required config field: %s", field)
		}
	}
	return nil
}

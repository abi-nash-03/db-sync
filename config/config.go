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

var AppConfig *Config

type SSHConfig struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	User string `yaml:"user"`
	KeyPath string `yaml:"key_path"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSH      SSHConfig `yaml:"ssh"`

}

func LoadConfig(configPath string) (error) {

	// Read the YAML file content
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading YAML file: %w", err)
	}


	if err := yaml.Unmarshal(yamlFile, &AppConfig); err != nil {
		return fmt.Errorf("error unmarshalling YAML: %w", err)
	}

	if err := AppConfig.Validate(); err != nil {
		return fmt.Errorf("error validating config: %w", err)
	}

	fmt.Println("Config loaded successfully")
	fmt.Println(AppConfig)	
	
	return nil
}


func (c *Config) Validate() error {
    required := map[string]string{
        "source.host":     c.Source.Host,
        "source.user":     c.Source.User,
        "source.password": c.Source.Password,
        "source.database": c.Source.Database,
		"source.ssh.host": c.Source.SSH.Host,
		"source.ssh.user": c.Source.SSH.User,
		"source.ssh.key_path": c.Source.SSH.KeyPath,
        "dest.host":       c.Destination.Host,
        "dest.user":       c.Destination.User,
        "dest.password":   c.Destination.Password,
        "dest.database":   c.Destination.Database,
		"dest.ssh.host": c.Destination.SSH.Host,
		"dest.ssh.user": c.Destination.SSH.User,
		"dest.ssh.key_path": c.Destination.SSH.KeyPath,
    }

    for field, value := range required {
        if value == "" {
            return fmt.Errorf("missing required config field: %s", field)
        }
    }

    // validate integers separately
    if c.Source.Port <= 0 || c.Source.Port > 65535 {
        return fmt.Errorf("invalid source.port: %d (must be 1-65535)", c.Source.Port)
    }
    if c.Destination.Port <= 0 || c.Destination.Port > 65535 {
        return fmt.Errorf("invalid dest.port: %d (must be 1-65535)", c.Destination.Port)
    }

	// validate ssh ports
	if c.Source.SSH.Port <= 0 || c.Source.SSH.Port > 65535 {
		return fmt.Errorf("invalid source.ssh.port: %d (must be 1-65535)", c.Source.SSH.Port)
	}
	if c.Destination.SSH.Port <= 0 || c.Destination.SSH.Port > 65535 {
		return fmt.Errorf("invalid dest.ssh.port: %d (must be 1-65535)", c.Destination.SSH.Port)
	}

    return nil
}

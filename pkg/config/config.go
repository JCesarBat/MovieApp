package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// DBConfig contains the database connection configuration.
type DBConfig struct {
	ConnectionString string `json:"connection_string"`
}

// Config represent all the configuration in the jsonfile.
type Config struct {
	DB DBConfig `json:"DB"`
}

var (
	instance *Config
	once     sync.Once
)

// LoadConfig load the configuration from a jsonFile.
func LoadConfig(filename string) (*Config, error) {
	var config Config

	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading the configuration file: %w", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing the configuration: %w", err)
	}

	return &config, nil
}

// GetConfig return a singleton instance of a config file.
func GetConfig() *Config {
	once.Do(func() {
		var err error
		instance, err = LoadConfig("config.json")
		if err != nil {
			instance, err = LoadConfig("../config.json")
			if err != nil {
				panic(fmt.Sprintf("Error loading configuración: %v", err))
			}
		}
	})
	return instance
}

// GetDBConnectionString return the string in the configuration
func (c *Config) GetDBConnectionString() string {
	return c.DB.ConnectionString
}

// GetDBConfig return the all configuration of the database.
func (c *Config) GetDBConfig() DBConfig {
	return c.DB
}

// New generate a new configuration
func New(connectionString string) *Config {
	return &Config{
		DB: DBConfig{
			ConnectionString: connectionString,
		},
	}
}

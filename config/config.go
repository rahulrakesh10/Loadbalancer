package config

import (
	"encoding/json"
	"load-balancer/server"
	"os"
)

// ServerConfig represents a backend server configuration
type ServerConfig struct {
	URL string `json:"url"`
}

// Config represents the load balancer configuration
type Config struct {
	Backends []ServerConfig `json:"backends"`
	Port     int            `json:"port"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetBackendServers converts config to BackendServer instances
func (c *Config) GetBackendServers() []*server.BackendServer {
	backends := make([]*server.BackendServer, len(c.Backends))
	for i, cfg := range c.Backends {
		backends[i] = server.NewBackendServer(cfg.URL)
	}
	return backends
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Backends: []ServerConfig{
			{URL: "http://localhost:9001"},
			{URL: "http://localhost:9002"},
			{URL: "http://localhost:9003"},
		},
		Port: 8080,
	}
}



package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ServerURL         string `json:"server_url"`
	AdminUsername     string `json:"admin_username"`
	AdminPassword     string `json:"admin_password"`
	TestcasesDir      string `json:"testcases_dir"`
	Concurrency       int    `json:"concurrency"`
	RetryCount        int    `json:"retry_count"`
	RetryDelaySeconds int    `json:"retry_delay_seconds"`
	TimeoutSeconds    int    `json:"timeout_seconds"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}

	cfg.normalizeURL()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	cfg.setDefaults()
	return &cfg, nil
}

func (c *Config) normalizeURL() {
	if !strings.HasPrefix(c.ServerURL, "http://") && !strings.HasPrefix(c.ServerURL, "https://") {
		c.ServerURL = "http://" + c.ServerURL
	}
	c.ServerURL = strings.TrimRight(c.ServerURL, "/")
}

func (c *Config) validate() error {
	if c.ServerURL == "" {
		return fmt.Errorf("server_url is required")
	}
	if c.AdminUsername == "" {
		return fmt.Errorf("admin_username is required")
	}
	if c.AdminPassword == "" {
		return fmt.Errorf("admin_password is required")
	}
	return nil
}

func (c *Config) setDefaults() {
	if c.TestcasesDir == "" {
		c.TestcasesDir = "./test-cases"
	}
	if c.Concurrency <= 0 {
		c.Concurrency = 4
	}
	if c.RetryCount <= 0 {
		c.RetryCount = 3
	}
	if c.RetryDelaySeconds <= 0 {
		c.RetryDelaySeconds = 2
	}
	if c.TimeoutSeconds <= 0 {
		c.TimeoutSeconds = 30
	}
}

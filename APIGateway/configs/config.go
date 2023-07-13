package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	APIEndpoint    string
	NacosIpAddr    string
	NacosPort      uint64
	CachingEnabled bool
	MaxQPS         int
	BurstSize      int
}

func InitConfig() (Configuration, error) {
	var config Configuration
	file, err := os.ReadFile("./configs/configs.json")
	if err != nil {
		return config, fmt.Errorf("unable to read config file: %w", err)
	}
	if err := json.Unmarshal(file, &config); err != nil {
		return config, fmt.Errorf("unable to parse JSON file: %w", err)
	}
	return config, nil
}

package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	GoogleConfig GoogleConfig `json:"google_config"`
	GithubConfig GithubConfig `json:"github_config"`
}

type GoogleConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type GithubConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func Load(path string) (*Config, error) {
	confJSON, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	conf := &Config{}
	err = json.Unmarshal(confJSON, conf)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file: %w", err)
	}

	return conf, nil
}

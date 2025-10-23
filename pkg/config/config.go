package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	UpstreamProxy *url.URL `mapstructure:"upstream_proxy"`
	NoProxy       []string `mapstructure:"no_proxy"`
	ListenAddr    string   `mapstructure:"listen_addr"`
}

func Load() (*Config, error) {
	viper.SetConfigName("proxyone")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/proxyone")
	viper.AddConfigPath("$HOME/.proxyone")
	viper.AddConfigPath(".")

	viper.SetDefault("listen_addr", "localhost:8080")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
		// Config file not found, use defaults
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Parse upstream proxy URL if provided
	if upstreamStr := viper.GetString("upstream_proxy"); upstreamStr != "" {
		u, err := url.Parse(upstreamStr)
		if err != nil {
			return nil, fmt.Errorf("invalid upstream proxy URL: %w", err)
		}
		config.UpstreamProxy = u
	}

	return &config, nil
}

func (c *Config) ShouldUseProxy(host string) bool {
	for _, noProxy := range c.NoProxy {
		if strings.Contains(host, noProxy) {
			return false
		}
	}
	return c.UpstreamProxy != nil
}

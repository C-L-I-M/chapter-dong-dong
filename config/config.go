package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type (
	ScrappingMode string

	Saga struct {
		Name           string         `mapstructure:"name" validate:"required"`
		SchedulingMode ScrappingMode  `mapstructure:"scheduling_mode" validate:"required"`
		Interval       time.Duration  `mapstructure:"interval" validate:"required"`
		Parameters     map[string]any `mapstructure:"parameters"`
		State          map[string]any `mapstructure:"state"`
	}

	Config struct {
		Sagas          map[string]*Saga `mapstructure:"sagas" validate:"required"`
		DiscordToken   string           `mapstructure:"discord_token" validate:"required"`
		ServerId       string           `mapstructure:"server_id" validate:"required"`
		StaticChannels []string         `mapstructure:"static_channels"`
	}
)

const (
	// SchedulingModeSequentialPageNotFound is a mode for scrapping sagas where urls have a sequential chapter number
	// in them. When a 404 is returned on the next chapter, no new chapter is detected, when the page is found, the
	// next page is scrapped.
	// Required parameters:
	// - "url" - the url of the first page, the %v will be replaced with the chapter number
	// - "start" - the first chapter number
	// - "not_found_status_code" - the status code to check for to detect a chapter not found
	// - "found_status_code" - the status code to check for to detect a chapter found
	SchedulingModeSequentialPageNotFound ScrappingMode = "sequential"
)

var configMutex sync.Mutex

func Load() (*Config, error) {
	configMutex.Lock()
	defer configMutex.Unlock()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %v ", err)
	}

	return &cfg, nil
}

func Save() error {
	configMutex.Lock()
	defer configMutex.Unlock()

	return viper.WriteConfig()
}

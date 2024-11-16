package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type (
	ScrappingMode string

	Saga struct {
		Name           string         `json:"name"`
		Slug           string         `json:"slug"`
		SchedulingMode ScrappingMode  `json:"scheduling_mode"`
		Parameters     map[string]any `json:"parameters"`
		State          map[string]any `json:"state"`
		Interval       time.Duration  `json:"interval"`
	}

	Config struct {
		Sagas []*Saga `json:"series"`
	}
)

const (
	// SchedulingModeSequentialPageNotFound is a mode for scrapping sagas where urls have a sequential chapter number
	// in them. When a 404 is returned on the next chapter, no new chapter is detected, when the page is found, the
	// next page is scrapped.
	// Required parameters:
	// - "url" - the url of the first page, the %i will be replaced with the chapter number
	// - "start" - the first chapter number
	// - "not_found_status_code" - the status code to check for to detect a chapter not found
	// - "found_status_code" - the status code to check for to detect a chapter found
	SchedulingModeSequentialPageNotFound ScrappingMode = "sequential"
)

var configMutex sync.Mutex

func Load(path string) (*Config, error) {
	configMutex.Lock()
	defer configMutex.Unlock()

	fp, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open config file: %v", path, err)
	}

	var cfg Config
	if err := json.NewDecoder(fp).Decode(&cfg); err != nil {
		log.Fatalf("%s: failed to decode config file: %v", path, err)
	}

	return &cfg, nil
}

func Save(path string, cfg *Config) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	fp, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%s: failed to create config file: %v", path, err)
	}

	if err := json.NewEncoder(fp).Encode(cfg); err != nil {
		return fmt.Errorf("%s: failed to encode config file: %v", path, err)
	}

	return nil
}

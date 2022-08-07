package main

import (
	"encoding/json"
	"os"
)

type GlobalConfig struct {
	DuplicateSongs []string `json:"duplicate_songs"`
	UpdateOrder []string `json:"update_order"`
	ExcludeFromAll []string `json:"exclude_from_all"`
}

var globalConfig GlobalConfig

func init() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		logFatalAndAlert(err)
	}

	err = json.Unmarshal(data, &globalConfig)
	if err != nil {
		logFatalAndAlert(err)
	}
}

func (*GlobalConfig) isExcludedFromAll(configFile string) bool {
	for _, excludedSong := range globalConfig.ExcludeFromAll {
		if configFile == excludedSong {
			return true
		}
	}
	return false
}

func (*GlobalConfig) isPriority(configFile string) bool {
	for _, priority := range globalConfig.UpdateOrder {
		if configFile == priority {
			return true
		}
	}
	return false
}
package main

import (
	"encoding/json"
	"os"
)

type GlobalConfig struct {
	Categories []string `json:"categories"`
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
	
	// exclude some files by default
	globalConfig.ExcludeFromAll = append(globalConfig.ExcludeFromAll, "images")
	for _, category := range globalConfig.Categories {
		globalConfig.ExcludeFromAll = append(globalConfig.ExcludeFromAll, category + ".json")
	}
}

func (c *GlobalConfig) isCategory(category string) bool {
	for _, c := range c.Categories {
		if c == category {
			return true
		}
	}
	return false
}

func (c *GlobalConfig) isExcludedFromAll(configFile string) bool {
	for _, excludedSong := range c.ExcludeFromAll {
		if configFile == excludedSong {
			return true
		}
	}
	return false
}

func (c *GlobalConfig) isPriority(configFile string) bool {
	for _, priority := range c.UpdateOrder {
		if configFile == priority {
			return true
		}
	}
	return false
}
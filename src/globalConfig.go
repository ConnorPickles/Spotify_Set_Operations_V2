package main

import (
	"encoding/json"
	"log"
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
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &globalConfig)
	if err != nil {
		log.Fatal(err)
	}
}
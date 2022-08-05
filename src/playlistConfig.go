package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type PlaylistConfig struct {
	Playlist1Name string    `json:"playlist1_name"`
	Playlist2Name string    `json:"playlist2_name"`
	SetPublic     bool      `json:"set_public"`
	UseExplicit   bool      `json:"use_explicit"`
	Operation     Operation `json:"operation"`
	Description   string    `json:"description"`
	Image         string    `json:"image"`
}

func loadPlaylistConfig(playlistName string) PlaylistConfig {
	var config PlaylistConfig
	playlistName += ".json"
	data, err := os.ReadFile("playlists/" + playlistName)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	config.setDescriptionFromOperation()
	return config
}

func createPringleConfig(playlistName string) PlaylistConfig {
	config := loadPlaylistConfig("Pringle")
	config.Playlist1Name = playlistName
	config.Description += "\"" + playlistName + "\""
	return config
}

func (config *PlaylistConfig) setDescriptionFromOperation() {
	if config.Description != "" {
		return
	}
	
	config.Description += "This playlist contains all songs that are in "
	switch config.Operation {
		case Intersection:
			config.Description += fmt.Sprintf("\"%s\" and \"%s\"", config.Playlist1Name, config.Playlist2Name)
		case Union:
			config.Description += fmt.Sprintf("\"%s\" or \"%s\"", config.Playlist1Name, config.Playlist2Name)
		case Difference:
			config.Description += fmt.Sprintf("\"%s\" but not \"%s\"", config.Playlist1Name, config.Playlist2Name)
	}
}
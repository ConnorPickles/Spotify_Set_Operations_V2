package main

import (
	"encoding/json"
	"fmt"
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
	var playlistConfig PlaylistConfig
	playlistName += ".json"
	data, err := os.ReadFile("playlists/" + playlistName)
	if err != nil {
		logFatalAndAlert(err)
	}
	err = json.Unmarshal(data, &playlistConfig)
	if err != nil {
		logFatalAndAlert(err)
	}
	if (playlistConfig.Operation < 0 || playlistConfig.Operation > 2) {
		logFatalAndAlert("Invalid operation. Must be 0, 1, or 2 (Intersection, Union, Difference).")
	}
	playlistConfig.setDescriptionFromOperation()
	return playlistConfig
}

func createPringleConfig(playlistName string) PlaylistConfig {
	playlistConfig := loadPlaylistConfig("Pringle")
	playlistConfig.Playlist1Name = playlistName
	playlistConfig.Description += "\"" + playlistName + "\""
	return playlistConfig
}

func (playlistConfig *PlaylistConfig) setDescriptionFromOperation() {
	if playlistConfig.Description != "" {
		return
	}
	
	playlistConfig.Description += "This playlist contains all songs that are in "
	switch playlistConfig.Operation {
		case Intersection:
			playlistConfig.Description += fmt.Sprintf("\"%s\" and \"%s\"", playlistConfig.Playlist1Name, playlistConfig.Playlist2Name)
		case Union:
			playlistConfig.Description += fmt.Sprintf("\"%s\" or \"%s\"", playlistConfig.Playlist1Name, playlistConfig.Playlist2Name)
		case Difference:
			playlistConfig.Description += fmt.Sprintf("\"%s\" but not \"%s\"", playlistConfig.Playlist1Name, playlistConfig.Playlist2Name)
	}
}
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify/v2"
)


func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: spotify_set_operations [update|create] [all|<playlist_name>]")
		return
	}

	client := authenticate()
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	
	createPlaylist(client, user.ID, os.Args[2])
}

func createPlaylist(client *spotify.Client, userID string, playlistName string) {
	var config PlaylistConfig
	if os.Args[2] == "Pringle" {
		if (len(os.Args) < 4) {
			fmt.Println("Usage: spotify_set_operations [update|create] Pringle <playlist_name>")
			os.Exit(1)
		}
		config = createPringleConfig(os.Args[3])
		playlistName = "Pringle " + os.Args[3]
	} else {
		config = loadPlaylistConfig(playlistName)
	}
		
	allPlaylists := getAllPlaylists(client, userID)
	playlistID1 := getPlaylistIDFromName(allPlaylists, config.Playlist1Name)
	playlistID2 := getPlaylistIDFromName(allPlaylists, config.Playlist2Name)
	tracks1 := getTracks(client, playlistID1)
	tracks2 := getTracks(client, playlistID2)
	tracksToAdd := executeOperation(config.Operation, nil, tracks1, tracks2, config.UseExplicit)
	
	playlist := createNewPlaylist(client, config, userID, playlistName)
	setPlaylistImage(client, playlist, config.Image)
	addTracksToPlaylist(client, playlist, tracksToAdd)
}
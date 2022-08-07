package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gen2brain/beeep"
	"github.com/zmb3/spotify/v2"
)


func main() {
	if len(os.Args) < 3 {
		printUsageAndExit()
	}
	if os.Args[1] != "update" && os.Args[1] != "create" {
		printUsageAndExit()
	}
	var createNew bool
	if os.Args[1] == "create" {
		createNew = true
	} else if os.Args[1] == "update" {
		createNew = false
	} else {
		printUsageAndExit()
	}

	client := authenticate()
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		logFatalAndAlert(err)
	}
	
	if os.Args[2] != "all" {
		var playlistConfig PlaylistConfig
		playlistName := os.Args[2]
		if os.Args[2] == "Pringle" {
			if (len(os.Args) != 4) {
				printPringleUsageAndExit()
			}
			playlistConfig = createPringleConfig(os.Args[3])
			playlistName += " " + os.Args[3]
		} else {
			playlistConfig = loadPlaylistConfig(os.Args[2])
		}
		
		operateOnPlaylist(client, playlistConfig, user.ID, playlistName, createNew)
		return
	}
	
	allPlaylistConfigs, allPlaylistNames := loadAllPlaylistConfigs()
	for i, playlistConfig := range allPlaylistConfigs {
		fmt.Printf("Working on \"%s\"...\n", allPlaylistNames[i])
		operateOnPlaylist(client, playlistConfig, user.ID, allPlaylistNames[i], createNew)
	}
}

func operateOnPlaylist(client *spotify.Client, playlistConfig PlaylistConfig, userID string, playlistName string, createNew bool) {
	allPlaylists := getAllPlaylists(client, userID)
	playlistID1 := getPlaylistIDFromName(allPlaylists, playlistConfig.Playlist1Name)
	playlistID2 := getPlaylistIDFromName(allPlaylists, playlistConfig.Playlist2Name)
	tracks1 := getTracks(client, playlistID1)
	tracks2 := getTracks(client, playlistID2)
	
	var tracksToAdd []spotify.SimpleTrack
	var tracksToRemove []spotify.SimpleTrack
	var playlist spotify.ID
	if createNew {
		playlist = createNewPlaylist(client, playlistConfig, userID, playlistName)
		setPlaylistImage(client, playlist, playlistConfig.Image)
		tracksToAdd, tracksToRemove = executeOperation(playlistConfig, nil, tracks1, tracks2)
	} else {
		playlist = getPlaylistIDFromName(allPlaylists, playlistName)
		existingTracks := getTracks(client, playlist)
		tracksToAdd, tracksToRemove = executeOperation(playlistConfig, existingTracks, tracks1, tracks2)
	}
	
	addTracksToPlaylist(client, playlist, tracksToAdd)
	removeTracksFromPlaylist(client, playlist, tracksToRemove)
}

func logFatalAndAlert(v ...any) {
	beeep.Alert("Spotify Set Operations", fmt.Sprint(v...), "")
	log.Fatal(v...)
}

func printUsageAndExit() {
	fmt.Println("Usage: spotify_set_operations [update|create] [all|<playlist_name>]")
	os.Exit(1)
}

func printPringleUsageAndExit() {
	fmt.Println("Usage: spotify_set_operations [update|create] Pringle <playlist_name>")
	os.Exit(1)
}
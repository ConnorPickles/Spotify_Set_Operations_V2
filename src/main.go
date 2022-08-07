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
		if globalConfig.isCategory(playlistName) {
			if (len(os.Args) != 4) {
				printUsageAndExit()
			}
			playlistConfig = createCategoryConfig(os.Args[3], playlistName)
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
	playlistID1, err := getPlaylistIDFromName(allPlaylists, playlistConfig.Playlist1Name)
	if err != nil {
		logFatalAndAlert(err)
	}
	playlistID2, err := getPlaylistIDFromName(allPlaylists, playlistConfig.Playlist2Name)
	if err != nil {
		logFatalAndAlert(err)
	}
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
		playlist, err = getPlaylistIDFromName(allPlaylists, playlistName)
		if err != nil && playlistConfig.CreateOnUpdate {
			// Create new playlist if it doesn't exist
			playlist = createNewPlaylist(client, playlistConfig, userID, playlistName)
			setPlaylistImage(client, playlist, playlistConfig.Image)
			tracksToAdd, tracksToRemove = executeOperation(playlistConfig, nil, tracks1, tracks2)
		} else if err != nil {
			logFatalAndAlert(err)
		} else {
			// Update existing playlist
			existingTracks := getTracks(client, playlist)
			tracksToAdd, tracksToRemove = executeOperation(playlistConfig, existingTracks, tracks1, tracks2)
			
			if (playlistConfig.DeleteIfEmpty && (len(existingTracks) + len(tracksToAdd) - len(tracksToRemove)) == 0) {
				deletePlaylist(client, playlist)
				return
			}
		}
	}
	
	addTracksToPlaylist(client, playlist, tracksToAdd)
	removeTracksFromPlaylist(client, playlist, tracksToRemove)
}

func logFatalAndAlert(v ...any) {
	beeep.Alert("Spotify Set Operations", fmt.Sprint(v...), "")
	log.Fatal(v...)
}

func printUsageAndExit() {
	fmt.Println("\nUSAGE:")
	fmt.Println("spotify_set_operations [update|create] [all|<playlist_name>]")
	fmt.Println("OR")
	fmt.Println("spotify_set_operations [update|create] <category> <playlist_name>")
	os.Exit(1)
}
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
		fmt.Println("Usage: spotify_set_operations [update|create] [all|<playlist_name>]")
		return
	}

	client := authenticate()
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		logFatalAndAlert(err)
	}
	
	createPlaylist(client, user.ID, os.Args[2])
}

func createPlaylist(client *spotify.Client, userID string, playlistName string) {
	var playlistConfig PlaylistConfig
	if os.Args[2] == "Pringle" {
		if (len(os.Args) < 4) {
			fmt.Println("Usage: spotify_set_operations [update|create] Pringle <playlist_name>")
			os.Exit(1)
		}
		playlistConfig = createPringleConfig(os.Args[3])
		playlistName = "Pringle " + os.Args[3]
	} else {
		playlistConfig = loadPlaylistConfig(playlistName)
	}
		
	allPlaylists := getAllPlaylists(client, userID)
	playlistID1 := getPlaylistIDFromName(allPlaylists, playlistConfig.Playlist1Name)
	playlistID2 := getPlaylistIDFromName(allPlaylists, playlistConfig.Playlist2Name)
	tracks1 := getTracks(client, playlistID1)
	tracks2 := getTracks(client, playlistID2)
	
	var tracksToAdd []spotify.SimpleTrack
	var tracksToRemove []spotify.SimpleTrack
	var playlist spotify.ID
	if os.Args[1] == "update" {
		playlist = getPlaylistIDFromName(allPlaylists, playlistName)
		existingTracks := getTracks(client, playlist)
		tracksToAdd, tracksToRemove = executeOperation(playlistConfig.Operation, existingTracks, tracks1, tracks2, playlistConfig.UseExplicit)
	} else {
		playlist = createNewPlaylist(client, playlistConfig, userID, playlistName)
		setPlaylistImage(client, playlist, playlistConfig.Image)
		tracksToAdd, tracksToRemove = executeOperation(playlistConfig.Operation, nil, tracks1, tracks2, playlistConfig.UseExplicit)
	}
	// check for something that isn't update or create
	
	addTracksToPlaylist(client, playlist, tracksToAdd)
	removeTracksFromPlaylist(client, playlist, tracksToRemove)
}

func logFatalAndAlert(v ...any) {
	beeep.Alert("Spotify Set Operations", fmt.Sprint(v...), "")
	log.Fatal(v...)
}
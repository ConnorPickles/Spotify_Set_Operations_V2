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
			if len(os.Args) != 4 {
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

	if globalConfig.RemoveNotLiked {
		removeNotLiked(client, user.ID)
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

	// doesn't need to be in an if statement because it will get overwritten if a new playlistID is created anyway
	playlistID, err := getPlaylistIDFromName(allPlaylists, playlistName)

	var tracksToAdd []spotify.SimpleTrack
	if createNew || (err != nil && playlistConfig.CreateOnUpdate) {
		tracksToAdd, _ = executeOperation(playlistConfig, nil, tracks1, tracks2)
		if len(tracksToAdd) == 0 && !createNew {
			// We still want to create a new empty playlist if we've been told to create a playlist.
			// It's useful when setting up playlists to see that the program did create a playlist,
			// it just didn't put any tracks in it.
			return
		}
		playlistID = createNewPlaylist(client, playlistConfig, userID, playlistName)
		setPlaylistImage(client, playlistID, playlistConfig.Image)
		addTracksToPlaylist(client, playlistID, tracksToAdd)
	} else {
		existingTracks := getTracks(client, playlistID)
		tracksToAdd, tracksToRemove := executeOperation(playlistConfig, existingTracks, tracks1, tracks2)

		if playlistConfig.DeleteIfEmpty && (len(existingTracks)+len(tracksToAdd)-len(tracksToRemove)) == 0 {
			deletePlaylist(client, playlistID)
			return
		}

		addTracksToPlaylist(client, playlistID, tracksToAdd)
		removeTracksFromPlaylist(client, playlistID, tracksToRemove)
	}
}

func removeNotLiked(client *spotify.Client, userID string) {
	allPlaylists := getAllPlaylists(client, userID)
	likedSongsID, err := getPlaylistIDFromName(allPlaylists, "Liked Songs")
	if err != nil {
		logFatalAndAlert(err)
	}

	likedTracks := Tracks(getTracks(client, likedSongsID)).toSimpleTracks()
	for _, playlist := range allPlaylists {
		if globalConfig.isUsingNotLikedSongs(playlist.Name) {
			continue
		}
		fmt.Printf("Removing un-liked songs from \"%s\"...\n", playlist.Name)
		tracks := Tracks(getTracks(client, playlist.ID)).toSimpleTracks()
		var onlyLiked []spotify.SimpleTrack
		for _, track := range tracks {
			for _, likedTrack := range likedTracks {
				if track.ID == likedTrack.ID {
					onlyLiked = append(onlyLiked, track)
					break
				}
			}
		}
		notLiked := difference(tracks, onlyLiked)
		removeTracksFromPlaylist(client, playlist.ID, notLiked)
	}
}

func logFatalAndAlert(v ...any) {
	beeep.Alert("Spotify Set Operations", fmt.Sprint(v...), "")
	log.Fatal(v...)
}

func printUsageAndExit() {
	fmt.Println("\nUSAGE:")
	fmt.Println("spotify_set_operations.exe [update|create] [all|<playlist_name>]")
	fmt.Println("OR")
	fmt.Println("spotify_set_operations.exe [update|create] <category> <playlist_name>\n")
	os.Exit(1)
}

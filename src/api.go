package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zmb3/spotify/v2"
)

// api wrapper functions
// --------------------------------------------------------------------------------------------------------------------

func getAllPlaylists(client *spotify.Client, userID string) []spotify.SimplePlaylist {
	var results []spotify.SimplePlaylist
	playlists, err := client.GetPlaylistsForUser(context.Background(), userID)
	if err != nil {
		logFatalAndAlert(err)
	}
	results = append(results, playlists.Playlists...)
	
	for {
		err = client.NextPage(context.Background(), playlists)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			logFatalAndAlert(err)
		}
		results = append(results, playlists.Playlists...)
	}
	
	return results
}


func getTracks(client *spotify.Client, playlistID spotify.ID) []Track {
	if playlistID == "Liked Songs" {
		return getSavedTracks(client)	
	}
	
	var results []spotify.PlaylistItem
	tracks, err := client.GetPlaylistItems(context.Background(), playlistID)
	if err != nil {
		logFatalAndAlert(err)
	}
	results = append(results, tracks.Items...)
	
	for {
		err = client.NextPage(context.Background(), tracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			logFatalAndAlert(err)
		}
		results = append(results, tracks.Items...)
	}
	
	return PlaylistItems(results).toTracks()
}

func getSavedTracks(client *spotify.Client) []Track {
	var results []spotify.SavedTrack
	tracks, err := client.CurrentUsersTracks(context.Background())
	if err != nil {
		logFatalAndAlert(err)
	}
	results = append(results, tracks.Tracks...)
	
	for {
		err = client.NextPage(context.Background(), tracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			logFatalAndAlert(err)
		}
		results = append(results, tracks.Tracks...)
	}
	
	return SavedTracks(results).toTracks()
}

func createNewPlaylist(client *spotify.Client, playlistConfig PlaylistConfig, userID string, playlistName string) spotify.ID {
	playlist, err := client.CreatePlaylistForUser(context.Background(), userID, playlistName, playlistConfig.Description, playlistConfig.SetPublic, false)
	if err != nil {
		logFatalAndAlert(err)
	}
	return playlist.ID
}

func setPlaylistImage(client *spotify.Client, playlist spotify.ID, image string) {
	if image == "" {
		return
	}
	
	r, err := os.Open("playlists/images/" + image)
	if err != nil {
		logFatalAndAlert(err)
	}
	defer r.Close()
	
	err = client.SetPlaylistImage(context.Background(), playlist, r)
	if err != nil {
		logFatalAndAlert(err)
	}
}

func addTracksToPlaylist(client *spotify.Client, playlistID spotify.ID, tracks []spotify.SimpleTrack) {
	if len(tracks) == 0 {
		return
	}
	
	convertedTracks := convertToTrackIDs(tracks)
	for i := 0; i < len(tracks); i += 100 {
		var currTracks []spotify.ID
		if i + 100 > len(tracks) {
			currTracks = convertedTracks[i:]
		} else {
			currTracks = convertedTracks[i:i+100]
		}
		
		_, err := client.AddTracksToPlaylist(context.Background(), playlistID, currTracks...)
		if err != nil {
			logFatalAndAlert(err)
		}
	}
}

func removeTracksFromPlaylist(client *spotify.Client, playlistID spotify.ID, tracks []spotify.SimpleTrack) {
	if len(tracks) == 0 {
		return
	}
	
	convertedTracks := convertToTrackIDs(tracks)
	for i := 0; i < len(tracks); i += 100 {
		var currTracks []spotify.ID
		if i + 100 > len(tracks) {
			currTracks = convertedTracks[i:]
		} else {
			currTracks = convertedTracks[i:i+100]
		}
		
		_, err := client.RemoveTracksFromPlaylist(context.Background(), playlistID, currTracks...)
		if err != nil {
			logFatalAndAlert(err)
		}
	}
}

func deletePlaylist(client *spotify.Client, playlistID spotify.ID) {
	err := client.UnfollowPlaylist(context.Background(), playlistID)
	if err != nil {
		logFatalAndAlert(err)
	}
}


// helper functions
// --------------------------------------------------------------------------------------------------------------------

func getPlaylistIDFromName(playlists []spotify.SimplePlaylist, playlistName string) (spotify.ID, error) {
	if playlistName == "Liked Songs" {
		return "Liked Songs", nil
	}
	
	for _, playlist := range playlists {
		if playlist.Name == playlistName {
			return playlist.ID, nil
		}
	}
	return "", fmt.Errorf("could not find playlist %s", playlistName)
}

func convertToTrackIDs(tracks []spotify.SimpleTrack) []spotify.ID {
	converted := []spotify.ID{}
	for _, track := range tracks {
		converted = append(converted, track.ID)
	}
	return converted
}
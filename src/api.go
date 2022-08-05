package main

import (
	"context"
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


func getTracks(client *spotify.Client, playlistID spotify.ID) []spotify.SimpleTrack {
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
	
	return convertToSimpleTracks(results)
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


// helper functions
// --------------------------------------------------------------------------------------------------------------------

func getPlaylistIDFromName(playlists []spotify.SimplePlaylist, playlistName string) spotify.ID {
	for _, playlist := range playlists {
		if playlist.Name == playlistName {
			return playlist.ID
		}
	}
	logFatalAndAlert("Could not find playlist " + playlistName)
	return ""
}

func convertToSimpleTracks(tracks []spotify.PlaylistItem) []spotify.SimpleTrack {
	converted := []spotify.SimpleTrack{}
	for _, track := range tracks {
		converted = append(converted, track.Track.Track.SimpleTrack)
	}
	return converted
}

func convertToTrackIDs(tracks []spotify.SimpleTrack) []spotify.ID {
	converted := []spotify.ID{}
	for _, track := range tracks {
		converted = append(converted, track.ID)
	}
	return converted
}
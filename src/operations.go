package main

import (
	"strings"
	"time"

	"github.com/zmb3/spotify/v2"
)

type Operation string

const (
	Intersection Operation = "intersection"
	Union = "union"
	Difference = "difference"
)


func executeOperation(playlistConfig PlaylistConfig, existingTracks, playlist1, playlist2 []Track) (add, remove []spotify.SimpleTrack) {
	if playlistConfig.OnlyNewSongs {
		var moreRecentOldestTime time.Time
		if oldest1, oldest2 := oldestAddedAt(playlist1), oldestAddedAt(playlist2); oldest1.After(oldest2) {
			moreRecentOldestTime = oldest1
		} else {
			moreRecentOldestTime = oldest2
		}
		
		playlist1 = removeSongsOlderThan(playlist1, moreRecentOldestTime)
		playlist2 = removeSongsOlderThan(playlist2, moreRecentOldestTime)
	}
	
	var op_result []spotify.SimpleTrack
	opPlaylist1 := Tracks(playlist1).toSimpleTracks()
	opPlaylist2 := Tracks(playlist2).toSimpleTracks()
	switch playlistConfig.Operation {
		case Intersection:
			op_result = intersection(opPlaylist1, opPlaylist2, playlistConfig.UseExplicit)
		case Union:
			op_result = union(opPlaylist1, opPlaylist2, playlistConfig.UseExplicit)
		case Difference:
			op_result = difference(opPlaylist1, opPlaylist2)
	}
	
	add = difference(op_result, Tracks(existingTracks).toSimpleTracks())
	remove = difference(Tracks(existingTracks).toSimpleTracks(), op_result)
	return add, remove
}

func intersection(playlist1, playlist2 []spotify.SimpleTrack, explicit bool) []spotify.SimpleTrack {
	var result []spotify.SimpleTrack
	
	for _, track1 := range playlist1 {
		for _, track2 := range playlist2 {
			if track1.ID == track2.ID {
				result = append(result, track1)
				continue
			}
			
			if !differentIDSameSong(track1, track2) {
				continue
			}
			
			if track1.Explicit == track2.Explicit {
				result = append(result, track1)
				continue
			}
			
			if explicit == track1.Explicit {
				result = append(result, track1)
			} else {
				result = append(result, track2)
			}
		}
	}
	
	return result
}

func union(playlist1, playlist2 []spotify.SimpleTrack, explicit bool) []spotify.SimpleTrack {
	var result []spotify.SimpleTrack
	result = append(result, playlist2...)
	
	for _, track1 := range playlist1 {
		for j, track2 := range playlist2 {
			if track1.ID == track2.ID {
				break
			}
			
			if differentIDSameSong(track1, track2) {
				if (track1.Explicit == track2.Explicit) {
					// this exact track has already been added
					break
				}
				
				// this song has been added, it may not be the correct version
				// update the result accordingly
				if explicit == track1.Explicit {
					result[j] = track1
				} else {
					result[j] = track2
				}
				
				break
			}
			
			if j == len(playlist2)-1 {
				result = append(result, track1)
			}
		}
	}
	
	return result
}

func difference(base, remove []spotify.SimpleTrack) []spotify.SimpleTrack {
	if remove == nil {
		return base
	}
	
	var result []spotify.SimpleTrack
	
	for _, trackBase := range base {
		for j, trackRemove := range remove {
			if trackBase.ID == trackRemove.ID || differentIDSameSong(trackBase, trackRemove) {
				break
			}
			
			if j == len(remove)-1 {
				result = append(result, trackBase)
			}
		}
	}
	
	return result
}

// used to identify songs that have explicit and non-explicit versions
func differentIDSameSong(track1, track2 spotify.SimpleTrack) bool {
	for _, name := range globalConfig.DuplicateSongs {
		if (strings.Contains(track1.Name, name) && strings.Contains(track2.Name, name)) {
			return true
		}	
	}
	
	return false
}
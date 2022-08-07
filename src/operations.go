package main

import (
	"strings"

	"github.com/zmb3/spotify/v2"
)

type Operation int

const (
	Intersection Operation = iota
	Union
	Difference
)


func executeOperation(operation Operation, existingTracks, playlist1, playlist2 []spotify.SimpleTrack, explicit bool) (add, remove []spotify.SimpleTrack) {
	var op_result []spotify.SimpleTrack
	switch operation {
		case Intersection:
			op_result = intersection(playlist1, playlist2, explicit)
		case Union:
			op_result = union(playlist1, playlist2, explicit)
		case Difference:
			op_result = difference(playlist1, playlist2)
	}
	
	add = difference(op_result, existingTracks)
	remove = difference(existingTracks, op_result)
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
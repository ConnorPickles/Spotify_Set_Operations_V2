package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/zmb3/spotify/v2"
)

type Operation int

const (
	Intersection Operation = iota
	Union
	Difference
)


func executeOperation(operation Operation, existingTracks, playlist1, playlist2 []spotify.SimpleTrack, explicit bool) []spotify.SimpleTrack {
	var result []spotify.SimpleTrack
	switch operation {
		case Intersection:
			result = intersection(playlist1, playlist2, explicit)
		case Union:
			result = union(playlist1, playlist2, explicit)
		case Difference:
			result = difference(playlist1, playlist2)
	}
	return difference(result, existingTracks)
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
	for _, name := range dups.Names {
		if (strings.Contains(track1.Name, name) && strings.Contains(track2.Name, name)) {
			return true
		}	
	}
	
	return false
}

type DuplicateSongs struct {
	Names []string `json:"names"`
}

var dups DuplicateSongs

func init() {
	// do this once on init so we're not doing file I/O every time the function is called
	data, err := os.ReadFile("duplicate_songs.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &dups)
	if err != nil {
		log.Fatal(err)
	}
}
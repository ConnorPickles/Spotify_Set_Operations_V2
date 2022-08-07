package main

import (
	"time"

	"github.com/zmb3/spotify/v2"
)

// A wrapper type to hold a simple track and other relevant information for this program
type Track struct {
	Track spotify.SimpleTrack
	AddedAt time.Time
}

type Tracks []Track

type PlaylistItems []spotify.PlaylistItem

type SavedTracks []spotify.SavedTrack

func oldestAddedAt(tracks []Track) time.Time {
	oldest := time.Now()
	for _, track := range tracks {
		if track.AddedAt.Before(oldest) {
			oldest = track.AddedAt
		}
	}
	return oldest
}

func removeSongsOlderThan(tracks []Track, cutoff time.Time) []Track {
	var result []Track
	for _, track := range tracks {
		if track.AddedAt.Before(cutoff) {
			continue
		}
		result = append(result, track)
	}
	return result
}

func (t Tracks) toSimpleTracks() []spotify.SimpleTrack {
	var converted []spotify.SimpleTrack
	for _, track := range t {
		converted = append(converted, track.Track)
	}
	return converted
}

func (items PlaylistItems) toTracks() []Track {
	converted := []Track{}
	for _, item := range items {
		time, err := time.Parse(spotify.TimestampLayout, item.AddedAt)
		if err != nil {
			logFatalAndAlert(err)
		}
		track := Track{
			Track: item.Track.Track.SimpleTrack,
			AddedAt: time,
		}
		converted = append(converted, track)
	}
	return converted
}

func (items SavedTracks) toTracks() []Track {
	converted := []Track{}
	for _, item := range items {
		time, err := time.Parse(spotify.TimestampLayout, item.AddedAt)
		if err != nil {
			logFatalAndAlert(err)
		}
		track := Track{
			Track: item.SimpleTrack,
			AddedAt: time,
		}
		converted = append(converted, track)
	}
	return converted
}
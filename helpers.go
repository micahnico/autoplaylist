package autoplaylist

import "github.com/zmb3/spotify"

func convertToSimpleArtists(list []spotify.FullArtist) []spotify.SimpleArtist {
	result := make([]spotify.SimpleArtist, len(list))
	for i, item := range list {
		result[i] = item.SimpleArtist
	}
	return result
}

func convertFromFullToSimpleTracks(list []spotify.FullTrack) []spotify.SimpleTrack {
	result := make([]spotify.SimpleTrack, len(list))
	for i, item := range list {
		result[i] = item.SimpleTrack
	}
	return result
}

func convertFromSavedToSimpleTracks(list []spotify.SavedTrack) []spotify.SimpleTrack {
	result := make([]spotify.SimpleTrack, len(list))
	for i, item := range list {
		result[i] = item.SimpleTrack
	}
	return result
}

func convertFromSimpleArtistsToIDs(list []spotify.SimpleArtist) []spotify.ID {
	result := make([]spotify.ID, len(list))
	for i, item := range list {
		result[i] = item.ID
	}
	return result
}

func convertFromSimpleTracksToIDs(list []spotify.SimpleTrack) []spotify.ID {
	result := make([]spotify.ID, len(list))
	for i, item := range list {
		result[i] = item.ID
	}
	return result
}

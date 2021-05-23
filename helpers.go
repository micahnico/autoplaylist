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

func convertFromPlaylistToSimpleTracks(list []spotify.PlaylistTrack) []spotify.SimpleTrack {
	result := make([]spotify.SimpleTrack, len(list))
	for i, item := range list {
		result[i] = item.Track.SimpleTrack
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

func uniqueArtists(artistSlice []spotify.SimpleArtist) []spotify.SimpleArtist {
	keys := make(map[spotify.ID]bool)
	list := []spotify.SimpleArtist{}
	for _, entry := range artistSlice {
		if _, value := keys[entry.ID]; !value {
			keys[entry.ID] = true
			list = append(list, entry)
		}
	}
	return list
}

func uniqueTracks(trackSlice []spotify.SimpleTrack) []spotify.SimpleTrack {
	keys := make(map[spotify.ID]bool)
	list := []spotify.SimpleTrack{}
	for _, entry := range trackSlice {
		if _, value := keys[entry.ID]; !value {
			keys[entry.ID] = true
			list = append(list, entry)
		}
	}
	return list
}

func remove(s []spotify.ID, i int) ([]spotify.ID, spotify.ID) {
	item := s[i]
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1], item
}

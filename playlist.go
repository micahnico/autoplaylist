package autoplaylist

import (
	"fmt"

	"github.com/zmb3/spotify"
)

type AutoPlaylist struct {
	spotifyClient *spotify.Client
	numTracks     *int
}

func NewAutoPlaylist(client *spotify.Client, numTracks int) *AutoPlaylist {
	p := new(AutoPlaylist)
	p.spotifyClient = client
	p.SetNumTracks(numTracks)
	return p
}

func (p *AutoPlaylist) SetNumTracks(num int) {
	if num < 0 {
		num = 0
	}
	p.numTracks = &num
}

func (p *AutoPlaylist) Create() error {
	var artists []spotify.SimpleArtist
	var tracks []spotify.SimpleTrack

	// get artists and tracks from playlists
	currPlaylists, err := p.spotifyClient.CurrentUsersPlaylists()
	if err != nil {
		return err
	}
	artists, tracks, err = p.getArtistsAndTracksFromPlaylists(currPlaylists.Playlists)

	// get top artists
	topArtistsPage, err := p.spotifyClient.CurrentUsersTopArtists()
	if err != nil {
		return err
	}
	artists = append(artists, convertToSimpleArtists(topArtistsPage.Artists)...)

	// get top tracks
	topTracksPage, err := p.spotifyClient.CurrentUsersTopTracks()
	if err != nil {
		return err
	}
	tracks = append(tracks, convertFromFullToSimpleTracks(topTracksPage.Tracks)...)

	// get saved tracks and artists
	savedTracks, err := p.spotifyClient.CurrentUsersTracks()
	tempArtists, tempTracks, err := p.getArtistsAndTracksFromSavedTracks(savedTracks)
	if err != nil {
		return err
	}
	artists = append(artists, tempArtists...)
	tracks = append(tracks, tempTracks...)

	// get the seeds and set options for the playlist
	artistIDs := convertFromSimpleArtistsToIDs(uniqueArtists(artists))
	trackIDs := convertFromSimpleTracksToIDs(uniqueTracks(tracks))

	// TODO: choose at random 2 artists and 3 tracks 5x to get recommendations for
	playlistSeeds := spotify.Seeds{
		Artists: artistIDs,
		Tracks:  trackIDs,
	}
	playlistOptions := spotify.Options{
		Limit: p.numTracks,
	}

	recommendations, err := p.spotifyClient.GetRecommendations(playlistSeeds, nil, &playlistOptions)
	if err != nil {
		return err
	}
	playlistTracks := recommendations.Tracks

	fmt.Println(playlistTracks) // TODO: create the playlist

	return nil
}

func (p *AutoPlaylist) getArtistsAndTracksFromPlaylists(playlists []spotify.SimplePlaylist) ([]spotify.SimpleArtist, []spotify.SimpleTrack, error) {
	var returnArtists []spotify.SimpleArtist
	var returnTracks []spotify.SimpleTrack

	for _, playlist := range playlists {
		tracks, err := p.spotifyClient.GetPlaylistTracks(playlist.ID)
		if err != nil {
			return nil, nil, err
		}

		// each page only has up to 100 results so go through all of them
		for page := 1; ; page++ {
			err = p.spotifyClient.NextPage(tracks)
			if err == spotify.ErrNoMorePages {
				break
			}
			if err != nil {
				return nil, nil, err
			}
			returnTracks = append(returnTracks, convertFromPlaylistToSimpleTracks(tracks.Tracks)...)
		}
	}

	for _, track := range returnTracks {
		for _, artist := range track.Artists {
			returnArtists = append(returnArtists, artist)
		}
	}

	return uniqueArtists(returnArtists), uniqueTracks(returnTracks), nil
}

func (p *AutoPlaylist) getArtistsAndTracksFromSavedTracks(savedTracks *spotify.SavedTrackPage) ([]spotify.SimpleArtist, []spotify.SimpleTrack, error) {
	var err error
	var returnArtists []spotify.SimpleArtist
	var returnTracks []spotify.SimpleTrack

	// each page only has up to 100 results so go through all of them
	for page := 1; ; page++ {
		err = p.spotifyClient.NextPage(savedTracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		returnTracks = append(returnTracks, convertFromSavedToSimpleTracks(savedTracks.Tracks)...)
	}

	for _, track := range returnTracks {
		for _, artist := range track.Artists {
			returnArtists = append(returnArtists, artist)
		}
	}

	return uniqueArtists(returnArtists), uniqueTracks(returnTracks), nil
}

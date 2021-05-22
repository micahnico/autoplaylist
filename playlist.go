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

	// get followed artists
	followedArtistsPage, err := p.spotifyClient.CurrentUsersFollowedArtists()
	if err != nil {
		return err
	}
	artists = append(artists, convertToSimpleArtists(followedArtistsPage.Artists)...)

	// get top tracks
	topTracksPage, err := p.spotifyClient.CurrentUsersTopTracks()
	if err != nil {
		return err
	}
	tracks = append(tracks, convertFromFullToSimpleTracks(topTracksPage.Tracks)...)

	// get saved tracks
	savedTracksPage, err := p.spotifyClient.CurrentUsersTracks()
	if err != nil {
		return err
	}
	tracks = append(tracks, convertFromSavedToSimpleTracks(savedTracksPage.Tracks)...)

	// TODO: use GetRelatedArtists(id) on artists

	// get the seeds and set options for the playlist
	artistIDs := convertFromSimpleArtistsToIDs(artists)
	trackIDs := convertFromSimpleTracksToIDs(tracks)
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
		page, err := p.spotifyClient.GetPlaylistTracks(playlist.ID)
		if err != nil {
			return nil, nil, err
		}
		pTracks := page.Tracks
		for _, track := range pTracks {
			returnTracks = append(returnTracks, track.Track.SimpleTrack)
		}
	}

	for _, track := range returnTracks {
		for _, artist := range track.Artists {
			returnArtists = append(returnArtists, artist)
		}
	}

	return returnArtists, returnTracks, nil
}

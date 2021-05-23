package autoplaylist

import (
	"math/rand"
	"time"

	"github.com/zmb3/spotify"
)

type AutoPlaylist struct {
	spotifyClient *spotify.Client
	name          string
	description   string
	numTracks     int
	playlistID    spotify.ID
}

func NewAutoPlaylist(client *spotify.Client, name string, description string, numTracks int, playlistID spotify.ID) *AutoPlaylist {
	p := new(AutoPlaylist)
	p.spotifyClient = client
	p.name = name
	p.description = description
	p.setNumTracks(numTracks)
	p.playlistID = playlistID
	return p
}

func (p *AutoPlaylist) setNumTracks(num int) {
	if num < 1 {
		num = 1
	}
	if num > 100 {
		num = 100
	}
	p.numTracks = num
}

func (p *AutoPlaylist) Create() error {
	var artists []spotify.SimpleArtist
	var tracks []spotify.SimpleTrack
	var err error

	if p.playlistID == "" {
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
	} else {
		artists, tracks, err = p.getArtistsAndTracksFromPlaylist(p.playlistID)
		if err != nil {
			return err
		}
	}

	// get stuff from one playlist

	playlistTracks, err := p.getNewTracks(artists, tracks)
	if err != nil {
		return err
	}

	// create a new playlist and populate it
	currentUser, err := p.spotifyClient.CurrentUser()
	newPlaylist, err := p.spotifyClient.CreatePlaylistForUser(currentUser.ID, p.name, p.description, false)
	_, err = p.spotifyClient.AddTracksToPlaylist(newPlaylist.ID, convertFromSimpleTracksToIDs(playlistTracks)...)
	if err != nil {
		return err
	}

	return nil
}

func (p *AutoPlaylist) getArtistsAndTracksFromPlaylist(playlistID spotify.ID) ([]spotify.SimpleArtist, []spotify.SimpleTrack, error) {
	var returnArtists []spotify.SimpleArtist
	var returnTracks []spotify.SimpleTrack

	tracks, err := p.spotifyClient.GetPlaylistTracks(playlistID)
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

	// get all the artists from the playlist
	for _, track := range returnTracks {
		for _, artist := range track.Artists {
			returnArtists = append(returnArtists, artist)
		}
	}

	return uniqueArtists(returnArtists), uniqueTracks(returnTracks), nil
}

func (p *AutoPlaylist) getArtistsAndTracksFromPlaylists(playlists []spotify.SimplePlaylist) ([]spotify.SimpleArtist, []spotify.SimpleTrack, error) {
	var returnArtists []spotify.SimpleArtist
	var returnTracks []spotify.SimpleTrack

	for _, playlist := range playlists {
		pArtists, pTracks, err := p.getArtistsAndTracksFromPlaylist(playlist.ID)
		returnArtists = append(returnArtists, pArtists...)
		returnTracks = append(returnTracks, pTracks...)
		if err != nil {
			return nil, nil, err
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

func (p *AutoPlaylist) getNewTracks(artists []spotify.SimpleArtist, tracks []spotify.SimpleTrack) ([]spotify.SimpleTrack, error) {
	var playlistTracks []spotify.SimpleTrack
	var playlistSeeds spotify.Seeds
	tracksPerPass := p.numTracks / 5
	playlistOptions := spotify.Options{
		Limit: &tracksPerPass,
	}

	// get the ids of artists and tracks for the seeds
	artistIDs := convertFromSimpleArtistsToIDs(uniqueArtists(artists))
	trackIDs := convertFromSimpleTracksToIDs(uniqueTracks(tracks))

	// get recommendations 5 times since Spotify limits seeds
	for i := 0; i < 5; i++ {
		s := rand.NewSource(time.Now().UnixNano())
		r := rand.New(s)

		// get random artists and tracks
		seedArtistIDs := make([]spotify.ID, 2)
		for i := range seedArtistIDs {
			n := r.Intn(len(artistIDs))
			artistIDs, seedArtistIDs[i] = remove(artistIDs, n)
		}
		seedTrackIDs := make([]spotify.ID, 3)
		for i := range seedTrackIDs {
			n := r.Intn(len(trackIDs))
			trackIDs, seedTrackIDs[i] = remove(trackIDs, n)
		}

		// setup the seeds
		playlistSeeds = spotify.Seeds{
			Artists: seedArtistIDs,
			Tracks:  seedTrackIDs,
		}

		// get recommendations
		recommendations, err := p.spotifyClient.GetRecommendations(playlistSeeds, nil, &playlistOptions)
		if err != nil {
			return nil, err
		}
		playlistTracks = append(playlistTracks, recommendations.Tracks...)
	}

	return uniqueTracks(playlistTracks), nil
}

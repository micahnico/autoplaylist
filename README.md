# AutoPlaylist (for Spotify)

### This is a go package that will create a Spotify playlist based on your saved artists, saved tracks, tracks from your playlists, and your top tracks and artists.

It requires that you have a spotify client to pass into the constructor.
Can be request heavy depending on how many tracks you have. A library of around 5,000 tracks (saved and in playlists) could use around 100 or more requests to the client.

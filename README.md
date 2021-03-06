# AutoPlaylist (for Spotify)

### This is a Go package that will create a Spotify playlist based on saved artists, saved tracks, tracks from playlists, and top tracks and artists. It also has an option to create a new playlist based off of only one existing playlist.

&nbsp;

**Notes:**
1. It uses [github.com/zmb3/spotify](https://github.com/zmb3/spotify)
2. It requires that you have an authenticated spotify client with the correct permissions to pass into the function.
3. The range for the number of the new playlist's tracks is 5 to 500.

&nbsp;
## Example for whole music library
```go
	p, err := autoplaylist.NewAutoPlaylist(client, "name", "description", 100, "")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success:", p.ID)
	}
```
## Example for one playlist
```go
	p, err := autoplaylist.NewAutoPlaylist(client, "name", "description", 100, "3hPeTyReRDrbVqtUHvwfSp")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success:", p.ID)
	}
```

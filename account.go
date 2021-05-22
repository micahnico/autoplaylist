package spotifyplaylist

import (
	"github.com/zmb3/spotify"
)

func ConnectAccount(url string, clientID string, secretKey string) {
	auth := spotify.NewAuthenticator(url, spotify.ScopeUserReadPrivate)
	auth.SetAuthInfo(clientID, secretKey)
	// TODO: finish setting this up
}

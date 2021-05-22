package test

import (
	"os"
	"testing"

	"github.com/micahnico/spotifyplaylist"
)

func TestConnectAccount(t *testing.T) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT_URI")

	_, err := spotifyplaylist.ConnectAccount(redirectURI, clientID, clientSecret)
	if err != nil {
		t.Errorf("could not connect account")
	}
}

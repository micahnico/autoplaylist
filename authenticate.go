package spotifyplaylist

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/zmb3/spotify"
)

var (
	ch    = make(chan *spotify.Client)
	state = randStringRunes(42)
)
var auth spotify.Authenticator

func ConnectAccount(redirectURI string, clientID string, secretKey string) (*spotify.Client, error) {
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistReadPrivate, spotify.ScopePlaylistModifyPrivate, spotify.ScopeUserTopRead, spotify.ScopeUserLibraryRead)
	auth.SetAuthInfo(clientID, secretKey)

	// start an HTTP server to authenticate
	callbackEndpoint := strings.Split(redirectURI, "/")[3]
	http.HandleFunc(fmt.Sprintf("/%v", callbackEndpoint), completeAuth)
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	openBrowser(url)

	// wait for auth to complete
	client := <-ch

	return client, nil
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return err
	}
	return nil
}

func randStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

package spotifyplaylist

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/zmb3/spotify"
)

var (
	ch    = make(chan *spotify.Client)
	state = "abc123"
)
var auth spotify.Authenticator

func ConnectAccount(redirectURI string, clientID string, secretKey string) (*spotify.Client, error) {
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate)
	auth.SetAuthInfo(clientID, secretKey)

	// start an HTTP server to authenticate
	// TODO: make the /callback dynamic alongwith redirectURI
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
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

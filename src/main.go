package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
)

const redirectURI = "http://localhost:8888/callback"

var (
	scopes = []string{
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopePlaylistReadPrivate,
		spotifyauth.ScopePlaylistReadCollaborative,
		spotifyauth.ScopePlaylistModifyPublic,
		spotifyauth.ScopePlaylistModifyPrivate,
		spotifyauth.ScopeImageUpload,
		spotifyauth.ScopeUserLibraryRead,
		spotifyauth.ScopeUserLibraryModify,
	}
	state = "spotify_set_operations"
	auth  *spotifyauth.Authenticator
	ch    chan *spotify.Client
)

func init() {
	const clientId = "2f364f15087c4793886f0da1a331b2e8"
	arr, err := os.ReadFile("client_secret.txt")
	if err != nil {
		log.Print("Could not read client_secret file")
		log.Fatal(err)
	}
	clientSecret := string(arr)
	
	auth = spotifyauth.New(spotifyauth.WithClientID(clientId), spotifyauth.WithClientSecret(clientSecret), spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(scopes...))
	ch = make(chan *spotify.Client)
}

func main() {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8888", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gen2brain/beeep"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
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
		logFatalAndAlert(err)
	}
	clientSecret := string(arr)
	
	auth = spotifyauth.New(spotifyauth.WithClientID(clientId), spotifyauth.WithClientSecret(clientSecret), spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(scopes...))
	ch = make(chan *spotify.Client)
}

func authenticate() *spotify.Client {
	tok := loadToken()
	if tok == nil {
		return login()
	} else {
		return spotify.New(auth.Client(context.Background(), tok))
	}
}

func login() *spotify.Client{
	// first start an HTTP server
	server := &http.Server{Addr: ":8888"}
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logFatalAndAlert(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	beeep.Alert("Spotify Set Operations", "Login required! Please locate the output for this program and follow the provided link", "")

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		logFatalAndAlert(err)
	}
	fmt.Println("You are logged in as:", user.ID)
	
	if err := server.Close(); err != nil {
		log.Println(err)
	}
	return client
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		logFatalAndAlert(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		logFatalAndAlert("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	cacheToken(tok)
	ch <- client
}

func cacheToken(tok *oauth2.Token) {
	data, err := json.MarshalIndent(tok, "", "    ")
	if err != nil {
		log.Println(err)
	}
	err = os.WriteFile("token.json", data, 0666)
	if err != nil {
		log.Println(err)
	}
}

func loadToken() *oauth2.Token {
	data, err := os.ReadFile("token.json")
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		log.Println(err)
	}
	tok := &oauth2.Token{}
	err = json.Unmarshal(data, tok)
	if err != nil {
		log.Println(err)
	}
	return tok
}
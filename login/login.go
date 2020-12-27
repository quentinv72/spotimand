package login

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const (
	address     = "localhost:9000"
	redirectURI = "http://" + address + "/redirect"
)

var (
	state  string
	server = &http.Server{Addr: address, Handler: nil}
	// Auth is the login authorization strutc that will be used to create clients.
	Auth spotify.Authenticator
)

// Initialize the environment and assign value to auth variable
func init() {
	os.Setenv("SPOTIFY_ID", "d651facb9c8c47188af996f8e0816764")
	Auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate, spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserReadCurrentlyPlaying)
}

// Handle signin
func handleSignin(w http.ResponseWriter, r *http.Request) {
	codeChallenge := codeChallenger()
	generateState(&state)
	url := Auth.AuthURLWithOpts(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	http.Redirect(w, r, url, http.StatusFound)
}

// Obtain token from redirectURI
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	tok, err := Auth.TokenWithOpts(state, r, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	fmt.Fprintf(w, "Login Completed! \nYou can now go back to your command line!")
	// Write tokens to JSON file
	json, _ := json.Marshal(tok)
	ioutil.WriteFile("tokens.json", json, 0600)
	go func() {
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

}

// Login manages the login process
func Login() (bool, spotify.Client) {
	fmt.Printf("Please sign-in at http://%s/signin\n", address)
	http.HandleFunc("/redirect", handleRedirect)
	http.HandleFunc("/signin", handleSignin)
	// return the authenticator to make accessible in main
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	// Verify the tokens were written to JSON file.
	var tokens oauth2.Token
	var client spotify.Client
	jsonData, err := ioutil.ReadFile("tokens.json")
	if err != nil {
		fmt.Println(err)
		return false, client
	}
	if err := json.Unmarshal(jsonData, &tokens); err != nil {
		fmt.Println(err)
		return false, client
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		fmt.Println("Seems like we weren't able to fetch your tokens... :/")
		return false, client
	}
	client = Auth.NewClient(&tokens)
	return true, client
}

package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	// Auth is the login authorization struct that will be used to create clients.
	Auth spotify.Authenticator
	// TokenFile is a string of the filepath to tokens.json
	TokenFile string
)

// Initialize the environment and assign value to auth variable
func init() {
	os.Setenv("SPOTIFY_ID", "d651facb9c8c47188af996f8e0816764")
	Auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate, spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserReadCurrentlyPlaying)
	executablePath, _ := os.Executable()
	executableFolder, _ := filepath.Split(executablePath)
	TokenFile = filepath.Join(executableFolder, "tokens.json")
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
	err = ioutil.WriteFile(TokenFile, json, 0600)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

}

// Login manages the login process
func Login() (bool, spotify.Client) {
	var tokens oauth2.Token
	var client spotify.Client
	// Check if user was logged in past
	if _, err := os.Stat(TokenFile); err == nil {
		err := RefreshToken()
		if err != nil {
			// Just print the error and continue to signin
			fmt.Println(err)
		} else {
			err = validateTokens(&tokens)
			if err == nil {
				return true, Auth.NewClient(&tokens)
			}

		}

	}
	// If not previously logged in then signin through Spotify portal
	fmt.Printf("Please sign-in at http://%s/signin\n", address)
	http.HandleFunc("/redirect", handleRedirect)
	http.HandleFunc("/signin", handleSignin)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	// Verify the tokens were written to JSON file.
	if err := validateTokens(&tokens); err != nil {
		return false, client
	}
	client = Auth.NewClient(&tokens)
	return true, client
}

func validateTokens(tokens *oauth2.Token) error {
	jsonData, err := ioutil.ReadFile(TokenFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonData, tokens); err != nil {
		return err
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		return errors.New("No tokens")
	}

	return nil
}

package login

import (
	"fmt"
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
	auth   spotify.Authenticator
)

// Initialize the environment and assign value to auth variable
func init() {
	os.Setenv("SPOTIFY_ID", "d651facb9c8c47188af996f8e0816764")
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate)
}

// Handle signin
func handleSignin(w http.ResponseWriter, r *http.Request) {
	codeChallenge, _ := codeChallenger()
	err := generateState(&state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	url := auth.AuthURLWithOpts(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	http.Redirect(w, r, url, http.StatusFound)
}

// Obtain token from redirectURI
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.TokenWithOpts(state, r, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	client := auth.NewClient(tok)
	user, _ := client.CurrentUser()
	fmt.Fprintf(w, "Login Completed! %s", user.ID)

}

// Login manages the login process
func Login() {
	fmt.Printf("Please login at http://%s/signin\n", address)
	http.HandleFunc("/redirect", handleRedirect)
	http.HandleFunc("/signin", handleSignin)
	log.Fatal(server.ListenAndServe())
	// return the authenticator to make accessible in main

}

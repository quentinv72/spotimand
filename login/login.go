package login

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"net/http"
	"time"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_.~"
const clientID = "d651facb9c8c47188af996f8e0816764"
const address = "localhost:9000"
const redirectURI = "http://" + address + "/redirect"

var codeVerifier string = ""
var client = &http.Client{}
var server = &http.Server{Addr: address, Handler: nil}

// CodeChallenger creates a base64url code challenger using SHA256 hash.
func CodeChallenger() (string, error) {
	codeVerifierBytes, err := codeVerifierGenerator(&codeVerifier) // Possible test by mocking the codeVerifierGenerator
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(codeVerifierBytes)
	return base64.URLEncoding.EncodeToString(hash[:]), nil
}

// codeVerifierGenerator generates a code verifier of length between 43 and 127.
// The code verifier is made up of the characters in the chars constant.
// The function updates the codeVerifier variable.
func codeVerifierGenerator(s *string) ([]byte, error) {
	rand.Seed(time.Now().UnixNano())
	length := rand.Intn(85) + 43
	bytes := make([]byte, length)
	if _, err := crand.Read(bytes); err != nil {
		return []byte{}, err
	}
	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}
	*s = string(bytes)
	return bytes, nil
}

// loginURL Return the login query
// Might be easier to use string formatting instead of URL to object
func loginRequest() *http.Request {
	req, _ := http.NewRequest("GET", "https://accounts.spotify.com/authorize", nil) // Need to handle the error
	query := req.URL.Query()
	codeChallenge, _ := CodeChallenger() // Need to handle error
	query.Set("client_id", clientID)
	query.Set("response_type", "code")
	query.Set("redirect_uri", redirectURI)
	query.Set("code_challenge_method", "S256")
	query.Set("code_challenge", codeChallenge)
	req.URL.RawQuery = query.Encode()
	// Need to add state and scope
	return req
}

// Create handler for /signin, /redirect, /success, /failure

func handleSignin(w http.ResponseWriter, r *http.Request) {
	url := loginRequest().URL
	http.Redirect(w, r, url.String(), http.StatusFound)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if _, ok := query["code"]; ok {
		http.Redirect(w, r, "/success", http.StatusFound)
	}

}

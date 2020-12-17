package main

import (
	crand "crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_.~"

func main() {
	random, err := randomString()
	if err != nil {
		log.Fatalln(err)
	}

	signInURL := "https://accounts.spotify.com/authorize?" + "client_id=d651facb9c8c47188af996f8e0816764&response_type=code&redirect_uri=http://localhost:9000/success&code_challenge_method=S256&code_challenge=" + random
	handleSignIn := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, signInURL, http.StatusPermanentRedirect)
	}
	fmt.Println("Sign in to your Spotify account at http://localhost:9000/signin")
	http.HandleFunc("/signin", handleSignIn)
	http.HandleFunc("/success", handleLogged)
	log.Fatal(http.ListenAndServe("localhost:9000", nil))
}

func randomString() (string, error) {
	rand.Seed(time.Now().UnixNano())
	length := rand.Intn(85) + 43
	bytes := make([]byte, length)

	if _, err := crand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return fmt.Sprintf("%x", sha256.Sum256(bytes)), nil
}

func handleLogged(w http.ResponseWriter, r *http.Request) {
	// query := r.URL.Query()
	// if value, ok := query["code"]; ok {
	// http.Post("https://accounts.spotify.com/api/token")

}

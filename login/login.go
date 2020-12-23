package login

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_.~"
const clientID = "d651facb9c8c47188af996f8e0816764"
const address = "localhost:9000"
const redirectURI = "http://" + address + "/redirect"

var codeVerifier string
var client = &http.Client{}
var server = &http.Server{Addr: address, Handler: nil}

// CodeChallenger creates a base64url code challenger using SHA256 hash.
func CodeChallenger() (string, error) {
	codeVerifierBytes, err := codeVerifierGenerator(&codeVerifier) // Possible test by mocking the codeVerifierGenerator
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(codeVerifierBytes)
	fmt.Println(codeVerifier, base64.URLEncoding.EncodeToString(hash[:]))
	encoding := base64.StdEncoding.EncodeToString(hash[:])
	encoding = strings.Replace(encoding, "+", "-", -1) // 62nd char of encoding
	encoding = strings.Replace(encoding, "/", "_", -1) // 63rd char of encoding
	encoding = strings.Replace(encoding, "=", "", -1)  // Remove any trailing '='s
	return encoding, nil
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

// Create handler for /signin, /redirect, /success

func handleSignin(w http.ResponseWriter, r *http.Request) {
	// Need to ad state

	codeChallenge, _ := CodeChallenger()
	url := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%v&code_challenge_method=S256&code_challenge=%v", clientID, redirectURI, codeChallenge)

	http.Redirect(w, r, url, http.StatusFound)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if _, ok := query["code"]; ok {
		successURL := fmt.Sprintf("/success?code=%v", query["code"][0])
		http.Redirect(w, r, successURL, http.StatusFound)

	} else {
		w.Write([]byte(query["error"][0]))
		w.Write([]byte(fmt.Sprintf("\nPlease try logging in again at http://%v/signin", address)))
	}

}

func handleSuccess(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"][0]
	url := "https://accounts.spotify.com/api/token"
	data := fmt.Sprintf(
		"grant_type=authorization_code&client_id=%s"+
			"&code_verifier=%s"+
			"&code=%s"+
			"&redirect_uri=%s",
		clientID, codeVerifier, code, redirectURI)
	w.Write([]byte(data))
	payload := strings.NewReader(data)
	// create the request and execute it
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var tokenBody map[string]interface{}
	err = json.Unmarshal(body, &tokenBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(body)

}

// Login function
func Login() {
	fmt.Printf("Please login at http://%s/signin\n", address)
	http.HandleFunc("/redirect", handleRedirect)
	http.HandleFunc("/success", handleSuccess)
	http.HandleFunc("/signin", handleSignin)
	log.Fatal(server.ListenAndServe())

}

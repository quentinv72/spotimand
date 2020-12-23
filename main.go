package main

import (
	"github.com/quentinv72/spotimand/login"
)

func main() {
	login.Login()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// signInURL := "https://accounts.spotify.com/authorize?" + "client_id=d651facb9c8c47188af996f8e0816764&response_type=code&redirect_uri=http://localhost:9000/success&code_challenge_method=S256&code_challenge=" + random
	// handleSignIn := func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, signInURL, http.StatusPermanentRedirect)
	// }
	// fmt.Println("Sign in to your Spotify account at http://localhost:9000/signin")
	// http.HandleFunc("/signin", handleSignIn)
	// http.HandleFunc("/success", handleLogged)
	// log.Fatal(http.ListenAndServe("localhost:9000", nil))
}

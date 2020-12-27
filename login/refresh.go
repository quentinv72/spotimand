package login

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/zmb3/spotify"
)

// RefreshToken refreshes the token and updates the spotify client
func RefreshToken(client *spotify.Client) {
	clientID := os.Getenv("SPOTIFY_ID")
	for {
		var tokens oauth2.Token
		time.Sleep(20 * time.Minute)
		jsonData, _ := ioutil.ReadFile("tokens.json")
		json.Unmarshal(jsonData, &tokens)
		data := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&client_id=%s",
			tokens.RefreshToken, clientID)
		payload := strings.NewReader(data)
		req, _ := http.NewRequest("POST", spotify.TokenURL, payload)
		req.Header.Add("content-type", "application/x-www-form-urlencoded")
		res, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &tokens)
		fmt.Println(tokens)
		res.Body.Close()
		*client = Auth.NewClient(&tokens)
		json, _ := json.Marshal(tokens)
		ioutil.WriteFile("tokens.json", json, 0600)

	}
}

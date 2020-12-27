package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/zmb3/spotify"
)

// RefreshToken refreshes the token and updates the spotify client
func RefreshToken() error {
	clientID := os.Getenv("SPOTIFY_ID")
	var tokens oauth2.Token
	jsonData, err := ioutil.ReadFile("tokens.json")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jsonData, &tokens); err != nil {
		return err
	}
	data := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&client_id=%s",
		tokens.RefreshToken, clientID)
	payload := strings.NewReader(data)
	req, err := http.NewRequest("POST", spotify.TokenURL, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return errors.New("No valid token found")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(body, &tokens)
	res.Body.Close()
	json, err := json.Marshal(tokens)
	if err != nil {
		return err
	}
	ioutil.WriteFile("tokens.json", json, 0600)
	return nil
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/quentinv72/spotimand/login"
	"github.com/quentinv72/spotimand/player"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var client spotify.Client

func main() {
	var logged bool
	logged, client = login.Login()
	if !logged {
		fmt.Println("There was an issue logging you in")
		return
	}
	fmt.Println("You are successfully logged in :)")
	// Goroutine to update the client with the refreshed token
	go func(client *spotify.Client) {
		for {
			var newTokens oauth2.Token
			time.Sleep(20 * time.Minute)
			login.RefreshToken()
			jsonData, _ := ioutil.ReadFile("tokens.json")
			json.Unmarshal(jsonData, &newTokens)
			*client = login.Auth.NewClient(&newTokens)
		}
	}(&client)
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s@spotimand> ", user.ID)
		// Read the keyboad input.
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		execInput(input) // Maybe handle error here??
	}
}

func execInput(input string) {
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\r\n")

	// Split the input separate the command and the arguments.
	args := strings.Split(input, " ")

	switch args[0] {
	case "play":
		player.Play(&client)
	case "pause":
		player.Pause(&client)
	case "next":
		player.Next(&client)
	case "previous":
		player.Previous(&client)
	case "current":
		player.CurrentlyPlaying(&client)
	case "exit":
		os.Exit(0)
	default:
		fmt.Fprintln(os.Stderr, "Not a command")
	}
}

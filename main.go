package main

import (
	"bufio"
	"encoding/json"
	"errors"
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
			jsonData, _ := ioutil.ReadFile(login.TokenFile)
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

		err = execInput(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execInput(input string) error {
	// Remove the newline character and carriage character.
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	// Split the input separate the command and the arguments.
	args := strings.Split(input, " ")

	switch args[0] {
	case "play":
		return player.Play(&client, args[1:])
	case "pause":
		return player.Pause(&client)
	case "next":
		return player.Next(&client)
	case "previous":
		return player.Previous(&client)
	case "current":
		return player.SongCurrentlyPlaying(&client)
	case "device":
		return player.Devices(&client, args[1:])
	case "exit":
		defer os.Exit(0)
		return nil
	case "":
		// Do nothing
		return nil
	default:
		return errors.New("not a command")
	}
}

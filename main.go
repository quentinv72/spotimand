package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/quentinv72/spotimand/login"
	"github.com/zmb3/spotify"
)

var client spotify.Client

func main() {
	logged, client := login.Login()
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
	user, _ := client.CurrentUser()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s@spotimand> ", user.ID)
		// Read the keyboad input.
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Handle the execution of the input.
		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execInput(input string) error {
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")

	// Prepare the command to execute.
	cmd := exec.Command(input)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}

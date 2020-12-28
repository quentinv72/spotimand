package player

import (
	"fmt"
	"os"

	"github.com/zmb3/spotify"
)

// Play plays the song that is currently playing
func Play(client *spotify.Client) {
	go CurrentlyPlaying(client)
	err := client.Play()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

// Pause pauses the osng that is currently playing
func Pause(client *spotify.Client) {
	err := client.Pause()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintln(os.Stdout, "Successfully paused")
}

// Next jumps to the next song in the queue
func Next(client *spotify.Client) {
	err := client.Next()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintln(os.Stdout, "Playing next song")
}

// Previous plays the previous song in the "queue"
func Previous(client *spotify.Client) {
	err := client.Previous()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintln(os.Stdout, "Playing previous song")
}

// CurrentlyPlaying prints the name of the current song
func CurrentlyPlaying(client *spotify.Client) {
	song, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		fmt.Fprintln(os.Stdout, "Couldn't fetch song name")
		return
	}
	trackName := song.Item.SimpleTrack.Name
	artist := song.Item.SimpleTrack.Artists[0].Name
	fmt.Fprintf(os.Stdout, "Currently playing: %s - %s\n", trackName, artist)
}

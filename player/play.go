package player

import (
	"flag"
	"fmt"
	"os"

	"github.com/zmb3/spotify"
)

// TODO if a device is available then default to playing on it
var flagsPlay = flag.NewFlagSet("Play", flag.ContinueOnError)
var flagsDevices = flag.NewFlagSet("Devices", flag.ContinueOnError)
var deviceID string
var deviceList bool

func init() {
	flagsPlay.StringVar(&deviceID, "device", "", "ID of device to play music on.")
	flagsDevices.BoolVar(&deviceList, "list", false, "List all availble devices.")
}

// Play plays the song that is currently playing
// and can also switch to play on another available device
func Play(client *spotify.Client, command []string) {
	err := flagsPlay.Parse(command)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// Switch and play on different device
	if deviceID != "" {
		err = client.TransferPlayback(spotify.ID(deviceID), true)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			// reset the default flag value (might be a cleaner way to do this)
			deviceID = ""
			return
		}
		// reset the default flag value (might be a cleaner way to do this)
		deviceID = ""
		return
	}
	err = client.Play()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

// Pause pauses the song that is currently playing
func Pause(client *spotify.Client) {
	err := client.Pause()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
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

// SongCurrentlyPlaying prints the name of the current song
func SongCurrentlyPlaying(client *spotify.Client) {
	song, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "Couldn't fetch song name")
		return
	}
	// Check if the song struct is populated
	if song.Timestamp == 0 {
		fmt.Fprintln(os.Stderr, "Seems like there aren't any active devices at the moment.")
		return
	}
	trackName := song.Item.SimpleTrack.Name
	artist := song.Item.SimpleTrack.Artists[0].Name
	fmt.Fprintf(os.Stdout, "Currently playing: %s - %s\n", trackName, artist)
	fmt.Fprintln(os.Stdout, "hello")

}

// Devices otputs teh current device the music is playing on
// or the list of all available devices
func Devices(client *spotify.Client, command []string) {
	err := flagsDevices.Parse(command)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// List all available devices
	if deviceList {
		devices, err := client.PlayerDevices()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			// reset the default flag value (might be a cleaner way to do this)
			deviceList = false
			return
		}
		if len(devices) == 0 {
			fmt.Fprintln(os.Stderr, "There are no available devices...")
			// reset the default flag value (might be a cleaner way to do this)
			deviceList = false
			return
		}
		fmt.Fprintln(os.Stdout, "The available devices are:")

		for _, device := range devices {
			fmt.Fprintf(os.Stdout, "Name: %s --- ID: %s\n", device.Name, device.ID)
		}
		// reset the default flag value (might be a cleaner way to do this)
		deviceList = false
		return
	}
	info, err := client.PlayerState()
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		return
	}
	if info.Device.ID == "" {
		fmt.Fprintln(os.Stderr, "Seems like you aren't playing music on any device...")
		return
	}
	fmt.Fprintf(os.Stdout, "Name: %s --- ID: %s\n", info.Device.Name, info.Device.ID)

}

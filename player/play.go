package player

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/zmb3/spotify"
)

var flagsPlay = flag.NewFlagSet("Play", flag.ContinueOnError)
var flagsDevices = flag.NewFlagSet("Devices", flag.ContinueOnError)
var deviceID string
var deviceList bool
var playStart bool

func init() {
	// initialize `play` flags
	flagsPlay.StringVar(&deviceID, "device", "", "ID of device to play music on.")
	flagsPlay.BoolVar(&playStart, "start", false, "Find and start playing music on an available device")
}

// Play plays the song that is currently playing
// and can also switch to play on another available device
func Play(client *spotify.Client, command []string) error {
	defer resetFlags()
	err := flagsPlay.Parse(command)
	if err != nil {
		return err
	}
	if playStart {
		devices, err := client.PlayerDevices()
		if err != nil {
			return err
		}
		if len(devices) == 0 {
			return errors.New("there are no available devices")
		}
		return client.PlayOpt(&spotify.PlayOptions{DeviceID: &devices[0].ID})
	}
	// Switch and play on different device
	if deviceID != "" {
		return client.TransferPlayback(spotify.ID(deviceID), true)
	}
	return client.Play()

}

// Pause pauses the song that is currently playing
func Pause(client *spotify.Client) error {
	return client.Pause()

}

// Next jumps to the next song in the queue
func Next(client *spotify.Client) error {
	return client.Next()

}

// Previous plays the previous song in the "queue"
func Previous(client *spotify.Client) error {
	return client.Previous()

}

// SongCurrentlyPlaying prints the name of the current song
func SongCurrentlyPlaying(client *spotify.Client) error {
	song, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		return err
	}
	// Check if the song struct is populated
	if song.Timestamp == 0 {
		return errors.New("seems like there aren't any active devices at the moment")
	}
	trackName := song.Item.SimpleTrack.Name
	artist := song.Item.SimpleTrack.Artists[0].Name
	fmt.Fprintf(os.Stdout, "Currently playing: %s - %s\n", trackName, artist)
	return nil
}

func init() {
	// initialize `device` flags
	flagsDevices.BoolVar(&deviceList, "list", false, "List all availble devices.")

}

// Devices outputs the current device the music is playing on
// or the list of all available devices
func Devices(client *spotify.Client, command []string) error {
	err := flagsDevices.Parse(command)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	// List all available devices
	if deviceList {
		defer resetFlags()
		devices, err := client.PlayerDevices()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		if len(devices) == 0 {
			return errors.New("there are no available devices")
		}
		fmt.Fprintln(os.Stdout, "The available devices are:")

		for _, device := range devices {
			fmt.Fprintf(os.Stdout, "Name: %s --- ID: %s\n", device.Name, device.ID)
		}
		return nil
	}
	// Show current device
	info, err := client.PlayerState()
	if err != nil {
		return err
	}
	if info.Device.ID == "" {
		return errors.New("seems like you aren't playing music on any device")
	}
	fmt.Fprintf(os.Stdout, "Name: %s --- ID: %s\n", info.Device.Name, info.Device.ID)
	return nil
}

// resetFlags resets the flags to their default values
func resetFlags() {
	deviceID = ""
	deviceList = false
	playStart = false
}

# Spotimand

Spotimand is a simple command-line interface to control your Spotify player. It requires you to start at least one Spotify client (iOS app, Android app, web app, etc), which you will then be able to control directly from the Spotimand shell.

At the moment, you can play a song, pause a song, switch to the previous or next song, as well as switch devices.

## How to get started

1. Download Golang
2. Clone repository
3. Run the following commands to start the Spotimand shell.

```shell
go install
spotimand
```
Don't forget to add your GOPATH binary to your PATH (you can add these lines in your `.zprofile` or `.profile` or etc)! 
```
export PATH=$PATH:$(go env GOPATH)/bin
```
## Current commands

```shell
play "Plays the current song on your active spotify client"
pause "Pauses the song that is currently playing"
device "Displays the name and ID of the device the music is currently playing on"
next "Switch to next song"
previous "Switch back previous song"
current "Shows the song that is currently playing"
play --start "Start playing music on an active device (if there is one)"
play --device <device_id> "Play music on device <device_id>"
device --list "List all available devices"
device -h "Help for device command"
play -h "Help for play command"
```

## Next Steps and improvements

1. Add search functionality
2. Add some tests

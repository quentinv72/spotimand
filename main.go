package main

import (
	"fmt"
	"os"

	"github.com/quentinv72/spotimand/login"
)

func main() {
	fmt.Println(os.Getenv("SPOTIFY_ID"))
	login.Login()
}

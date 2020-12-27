package main

import (
	"fmt"

	"github.com/quentinv72/spotimand/login"
)

func main() {
	if logged := login.Login(); !logged {
		fmt.Println("There was an issue logging you in")
	}
	fmt.Println("You are successfully logged in :)")

}

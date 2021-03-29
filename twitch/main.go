package twitch

import (
	"log"

	twitch "github.com/gempir/go-twitch-irc/v2"
)

// StartTwitchListner
func StartTwitchListner() {
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(parser)
	client.OnConnect(func() {
		log.Print("connected")
	})

	client.Join("koalalorenzo")

	log.Print("Connecting")
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

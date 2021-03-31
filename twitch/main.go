package twitch

import (
	"log"

	twitch "github.com/gempir/go-twitch-irc/v2"
)

// StartTwitchListner
func StartTwitchListner(channel string) {
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(parser)
	client.OnConnect(func() {
		log.Print("Twitch Client Connected")
	})

	client.Join(channel)

	log.Print("Twitch Client Connecting")
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
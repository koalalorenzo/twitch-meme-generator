package twitch

import (
	log "github.com/sirupsen/logrus"

	twitch "github.com/gempir/go-twitch-irc/v2"
)

// StartTwitchListner
func StartTwitchListner(channel string) {
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(parser)
	client.OnConnect(func() {
		log.Printf("Twitch Client Connected to channel %s", channel)
	})

	client.Join(channel)

	log.Print("Twitch Client Connecting")
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

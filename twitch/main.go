package twitch

import (
	twitch "github.com/gempir/go-twitch-irc/v2"
)

// StartTwitchListner
func StartTwitchListner() {
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(parser)

	client.Join("koalalorenzo")

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

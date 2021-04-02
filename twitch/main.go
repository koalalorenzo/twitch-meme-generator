package twitch

import (
	log "github.com/sirupsen/logrus"

	twitch "github.com/gempir/go-twitch-irc/v2"
)

// StartTwitchListner
func StartTwitchListner(channel string) {
	logWF := log.WithFields(log.Fields{
		"f":       "twitch.StartTwitchListner",
		"channel": channel,
	})

	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(parser)
	client.OnConnect(func() {
		logWF.Info("Twitch Client Connected")
	})

	client.Join(channel)

	logWF.Debug("Twitch Client Connecting")
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

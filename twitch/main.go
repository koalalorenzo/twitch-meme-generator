package twitch

import (
	"time"

	log "github.com/sirupsen/logrus"

	twitch "github.com/gempir/go-twitch-irc/v2"
)

var Client *twitch.Client
var ClientIsConnected = false

// StartTwitchListner
func StartTwitchListner(channel string) {
	logWF := log.WithFields(log.Fields{
		"f":       "twitch.StartTwitchListner",
		"channel": channel,
	})

	Client = twitch.NewAnonymousClient()
	Client.IdlePingInterval = 20 * time.Second
	Client.SendPings = true

	Client.OnPrivateMessage(parser)
	Client.OnConnect(func() {
		logWF.Info("Twitch Client Connected")
		ClientIsConnected = true
	})

	Client.Join(channel)

	logWF.Debug("Twitch Client Connecting")
	err := Client.Connect()
	if err != nil {
		panic(err)
	}
}

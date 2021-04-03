package http

import (
	log "github.com/sirupsen/logrus"
)

var (
	WSListnerChannels []chan string
)

// channelPipe will allow to spread the messages incoming on the main channel
// on multiple channels opened when there is a new connection.
// This allows to have multiple browser source connecting on different
// web sockets at the same time.
func channelPipe(mainChannel chan string) {
	logWF := log.WithFields(log.Fields{
		"f": "http.channelPipe",
	})
	logWF.Debug("Starting the Channel Pipe")

	for msg := range mainChannel {

		log.WithFields(log.Fields{
			"msg":      msg,
			"channels": len(WSListnerChannels),
		}).Debug("Sending messages to web socket channels")

		for _, ch := range WSListnerChannels {
			ch <- msg
		}
	}
}

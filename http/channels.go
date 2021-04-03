package http

var (
	WSListnerChannels []chan string
)

// channelPipe will allow to spread the messages incoming on the main channel
// on multiple channels opened when there is a new connection.
// This allows to have multiple browser source connecting on different
// web sockets at the same time.
func channelPipe(mainChannel chan string) {
	for msg := range mainChannel {
		for _, ch := range WSListnerChannels {
			ch <- msg
		}
	}
}

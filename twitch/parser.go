package twitch

import (
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v2"

	"gitlab.com/koalalorenzo/twitch-meme-generator/generator"
)

var (
	BotCommandPrefix = "!meme"
)

func parser(msg twitch.PrivateMessage) {
	if !strings.HasPrefix(msg.Message, BotCommandPrefix) {
		return
	}

	var cmdSlice []string
	for _, el := range strings.Split(msg.Message, " ") {
		if len(el) == 0 {
			continue
		}
		cmdSlice = append(cmdSlice, el)
	}

	// Test if the lenght of the command is actually long enough
	if len(cmdSlice) < 3 {
		return
	}

	// ["!meme", "rufu", "this is the text"]
	memeKind := cmdSlice[1]
	memeText := strings.Join(cmdSlice[2:], " ")

	go generator.GenerateMeme(memeKind, memeText)
}

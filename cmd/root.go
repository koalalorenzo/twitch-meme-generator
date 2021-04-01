package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"

	"gitlab.com/koalalorenzo/twitch-meme-generator/generator"
	"gitlab.com/koalalorenzo/twitch-meme-generator/http"
	"gitlab.com/koalalorenzo/twitch-meme-generator/twitch"
)

var rootCmd = &cobra.Command{
	Use:   "koalalorenzo-meme-generator",
	Short: "Koalalorenzo's Twitch Bot that generates meme for your stream",
	Long: `Run a server and a twitch bot capable of generating memes and display them
on a HTTP page.

Author: https://who.is.lorenzo.setale.me/?
Twitch Channel: https://twitch.tv/koalalorenzo
Source: https://gitlab.com/koalalorenzo/twitch-meme-generator
License: https://gitlab.com/koalalorenzo/twitch-meme-generator/-/blob/main/LICENSE
`,
	Run: runApp,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringP("assets", "a", "assets", "Path of the directory containing the images")

	rootCmd.Flags().StringP("channel", "c", "koalalorenzo", "Channel to comnnect to")
	rootCmd.Flags().StringP("server", "s", "0.0.0.0:8001", "address (hostname and port) to listen to")
	rootCmd.Flags().DurationP("display-time", "d", 10*time.Second, "The time a meme is displayed on screen")
}

func runApp(cmd *cobra.Command, args []string) {
	urlChan := make(chan string, 5)
	twitchChannelName, err := cmd.Flags().GetString("channel")
	if err != nil {
		log.Fatalf("Unable to read the channel: %s", err)
		return
	}

	serverAddr, err := cmd.Flags().GetString("server")
	if err != nil {
		log.Fatalf("Unable to read the server flag value: %s", err)
		return
	}

	assetsDirPath, err := cmd.Flags().GetString("assets")
	if err != nil {
		log.Fatalf("Unable to read the assets flag value: %s", err)
		return
	}

	displayTimeDuration, err := cmd.Flags().GetDuration("display-time")
	if err != nil {
		log.Fatalf("Unable to read the assets flag value: %s", err)
		return
	}

	generator.SetPkgConfig(urlChan, assetsDirPath)
	http.SetPkgConfig(urlChan, displayTimeDuration)

	// Start listening for messages
	go twitch.StartTwitchListner(twitchChannelName)

	// Start the HTTP server
	http.StartServer(serverAddr)
}

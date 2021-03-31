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
	Use:   "twitch-meme-bot",
	Short: "A Twitch Bot that generates meme",
	Long: `Run a server and a twitch bot capable of generating memes and display them
on a HTTP page.`,
	Run: runApp,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmd.yaml)")

	rootCmd.Flags().StringP("channel", "c", "koalalorenzo", "Channel to comnnect to")
	rootCmd.Flags().StringP("server", "s", "0.0.0.0:8001", "address (hostname and port) to listen to")
	rootCmd.Flags().StringP("assets", "a", "assets", "Path of the directory containing the images")

	rootCmd.Flags().DurationP("display-time", "d", 10*time.Second, "The time a meme is displayed on screen")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

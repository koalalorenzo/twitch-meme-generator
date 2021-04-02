package cmd

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
	cobra.OnInitialize(initViperEnvConfig)

	// Set common Flags
	rootCmd.PersistentFlags().StringP("loglevel", "l", "info", "sets log level to warn, info or debug")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.PersistentFlags().StringP("logformat", "f", "text", "sets log format to json or text")
	viper.BindPFlag("logformat", rootCmd.PersistentFlags().Lookup("logformat"))

	// Sets common
	rootCmd.Flags().StringP("assets", "a", "assets", "Path of the directory containing the images")
	viper.BindPFlag("assets", rootCmd.Flags().Lookup("assets"))

	rootCmd.Flags().StringP("channel", "c", "koalalorenzo", "Channel to comnnect to")
	viper.BindPFlag("channel", rootCmd.Flags().Lookup("channel"))

	rootCmd.Flags().StringP("server", "s", "0.0.0.0:8001", "address (hostname and port) to listen to")
	viper.BindPFlag("server", rootCmd.Flags().Lookup("server"))

	rootCmd.Flags().DurationP("display-time", "d", 10*time.Second, "The time a meme is displayed on screen")
	viper.BindPFlag("server", rootCmd.Flags().Lookup("server"))
}

func initViperEnvConfig() {
	viper.SetEnvPrefix("KTMG")
	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	switch viper.GetString("logformat") {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	}

	switch viper.GetString("loglevel") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func runApp(cmd *cobra.Command, args []string) {
	urlChan := make(chan string, 5)
	twitchChannelName := viper.GetString("channel")
	serverAddr := viper.GetString("server")
	assetsDirPath := viper.GetString("assets")
	displayTimeDuration := viper.GetDuration("display-time")

	generator.SetPkgConfig(urlChan, assetsDirPath)
	http.SetPkgConfig(urlChan, displayTimeDuration)

	// Start listening for messages
	go twitch.StartTwitchListner(twitchChannelName)

	// Start the HTTP server
	http.StartServer(serverAddr)
}

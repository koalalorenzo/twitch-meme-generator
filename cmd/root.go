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

	rootCmd.Flags().String("host", "0.0.0.0", "sets the Host to listen to (HTTP)")
	viper.BindPFlag("host", rootCmd.Flags().Lookup("host"))

	rootCmd.Flags().StringP("port", "p", "8000", "sets the port to listen to (HTTP)")
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))

	rootCmd.Flags().DurationP("display-time", "d", 10*time.Second, "The time a meme is displayed on screen")
	viper.BindPFlag("display_time", rootCmd.Flags().Lookup("display-time"))

	rootCmd.Flags().String("temp-path", "", "Specify a custom temporary dir path")
	viper.BindPFlag("temp_path", rootCmd.Flags().Lookup("temp-path"))

	rootCmd.Flags().BoolP("webhook-enable", "w", false, "Expose endpoint to trigger memes")
	viper.BindPFlag("webhook_enable", rootCmd.Flags().Lookup("webhook-enable"))

	rootCmd.Flags().String("webhook-username", "", "(Basic Auth) forces basic authentication for webhook")
	viper.BindPFlag("webhook_username", rootCmd.Flags().Lookup("webhook-username"))

	rootCmd.Flags().String("webhook-password", "", "(Basic Auth) forces basic authentication for webhook")
	viper.BindPFlag("webhook_password", rootCmd.Flags().Lookup("webhook-password"))
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

	// Get the PORT env variable (required by Heroku). This is required because
	// we have a perfix KTMG
	herokuPortEnv := os.Getenv("PORT")
	if herokuPortEnv != "" {
		viper.Set("port", herokuPortEnv)
	}

	// Show debug configuration
	log.WithFields(log.Fields(viper.AllSettings())).Debug("configuration")
}

func runApp(cmd *cobra.Command, args []string) {
	urlChan := make(chan string, 5)

	twitchChannelName := viper.GetString("channel")
	host := viper.GetString("host")
	port := viper.GetString("port")
	assetsDirPath := viper.GetString("assets")
	displayTimeDuration := viper.GetDuration("display_time")
	outputdir := viper.GetString("temp_path")

	if outputdir == "" {
		var err error
		outputdir, err = os.MkdirTemp(os.TempDir(), "meme-generator")
		if err != nil {
			log.Errorf("Error creating temp dir: %s", err.Error())
			return
		}
	}

	generator.SetPkgConfig(urlChan, assetsDirPath, outputdir)

	httpConf := &http.Config{
		MainChannel:       urlChan,
		ChannelName:       twitchChannelName,
		DisplayTimePeriod: displayTimeDuration,
	}

	httpConf.Webhook.Enabled = viper.GetBool("webhook_enable")
	httpConf.Webhook.Username = viper.GetString("webhook_username")
	httpConf.Webhook.Password = viper.GetString("webhook_password")

	http.SetPkgConfig(httpConf)

	// Start listening for messages
	go twitch.StartTwitchListner(twitchChannelName)

	// Start the HTTP server
	serverAddr := fmt.Sprintf("%s:%s", host, port)
	http.StartServer(serverAddr)
}

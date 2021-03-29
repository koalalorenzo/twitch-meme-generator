package cmd

import (
	"github.com/spf13/cobra"

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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runApp(cmd *cobra.Command, args []string) {
	// Start listening for messages
	twitch.StartTwitchListner()
}

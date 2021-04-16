package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/koalalorenzo/twitch-meme-generator/generator"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a meme manually",
	Long: `Generate a meme and save it into a file. Use this to test the assets.
The usage is the same as the chat, meaning the first word/argument is the meme
name, followed by the text. Example:

   koalalorenzo-meme-generator generate cat hello human

Will generate a new image using the "cat" file with the text "hello human".
`,
	Run:  runGenerate,
	Args: cobra.MinimumNArgs(2),
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("assets", "a", "assets", "Path of the directory containing the images")
	viper.BindPFlag("assets", generateCmd.Flags().Lookup("assets"))
}

func runGenerate(cmd *cobra.Command, args []string) {
	assetsDirPath := viper.GetString("assets")

	urlChan := make(chan string, 1)
	memeKind := args[0]
	phrase := strings.Join(args[1:], " ")

	outputdir, err := os.MkdirTemp(os.TempDir(), "meme-generator")
	if err != nil {
		log.Errorf("Error creating temp dir: %s", err.Error())
		return
	}

	generator.SetPkgConfig(urlChan, assetsDirPath, outputdir)
	go generator.GenerateMeme(memeKind, phrase)

	// Copy the file here, wait for the file to be generated
	filename := <-urlChan
	tempFilePath := path.Join(generator.OutputTempDir, filename)

	tmpFileContent, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		log.Fatalf("Error reading content of temp file: %s", err.Error())
	}

	newFileName := fmt.Sprintf("%s%s", memeKind, path.Ext(tempFilePath))
	ioutil.WriteFile(newFileName, tmpFileContent, 0644)
	log.Infof("New file created at %s", newFileName)
}

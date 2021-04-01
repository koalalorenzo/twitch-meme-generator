package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/spf13/cobra"

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
}

func runGenerate(cmd *cobra.Command, args []string) {
	assetsDirPath, err := cmd.Flags().GetString("assets")
	if err != nil {
		log.Fatalf("Unable to read the assets flag value: %s", err)
		return
	}

	urlChan := make(chan string, 1)
	memeKind := args[0]
	phrase := strings.Join(args[1:], " ")

	generator.SetPkgConfig(urlChan, assetsDirPath)
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
	log.Printf("New file created at %s", newFileName)
}

package generator

import (
	"log"
	"strings"
)

// MemeFiles these are the file names of the memes available from the bot
var MemeFiles *[]string

func tickerFilesLiveLoade() {
	mfs, err := getMemesFilesData()
	if err != nil {
		log.Panicf("Erorr loading files: %s", err.Error())
	}

	mfsstrign := strings.Join(mfs, ", ")
	if mfsstrign != strings.Join(*MemeFiles, ", ") {
		log.Printf("Loaded memes: %s", strings.Join(mfs, ", "))
	}
	MemeFiles = &mfs
}

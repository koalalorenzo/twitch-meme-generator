package generator

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// MemeFiles these are the file names of the memes available from the bot
var MemeFiles *[]string

func tickerFilesLiveLoade() {
	mfs, err := getMemesFilesData()
	if err != nil {
		log.Errorf("Erorr loading files: %s", err.Error())
	}

	mfsstrign := strings.Join(mfs, ", ")
	if mfsstrign != strings.Join(*MemeFiles, ", ") {
		log.Debugf("Loaded memes: %s", strings.Join(mfs, ", "))
	}
	MemeFiles = &mfs
}

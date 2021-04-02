package generator

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// MemeFiles these are the file names of the memes available from the bot
var MemeFiles *[]string

func tickerFilesLiveLoade() {
	logWF := log.WithFields(log.Fields{
		"f":      "generator.tickerFilesLiveLoade",
		"assets": assetsDirPath,
	})

	mfs, err := getMemesFilesData()
	if err != nil {
		logWF.Errorf("Erorr loading files: %s", err.Error())
	}

	mfsstrign := strings.Join(mfs, ", ")
	if mfsstrign != strings.Join(*MemeFiles, ", ") {
		logWF.Debugf("Loaded memes: %s", strings.Join(mfs, ", "))
	}
	MemeFiles = &mfs
}

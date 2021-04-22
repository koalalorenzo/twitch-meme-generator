package generator

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jpoz/gomeme"
)

var (
	// OutputTempDir is a temporary directory where we will write our images
	OutputTempDir string
	//
	urlChan chan string
)

func SetPkgConfig(ch chan string, assetPath, tempPath string) {
	urlChan = ch
	AssetsDirPath = assetPath
	OutputTempDir = tempPath

	MemeFiles = []*MemeFile{}
	tickerFilesLiveLoade()
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for _ = range ticker.C {
			tickerFilesLiveLoade()
		}
	}()
}

// GenerateMeme does what it says
func GenerateMeme(kind, text string) {
	logWF := log.WithFields(log.Fields{
		"f":    "generator.GenerateMeme",
		"kind": kind,
		"text": text,
	})

	logWF.Infof(`Generating...`)

	memeCfg := gomeme.NewConfig()
	memeCfg.BottomText = text

	memeFile := &MemeFile{}
	for _, m := range MemeFiles {
		if strings.HasPrefix(m.Kind, kind) {
			memeFile = m
			break
		}
	}

	if memeFile.Filename == "" {
		logWF.Warn("Meme kind not found")
		return
	}

	// Updating the logs with relevant extra Fields
	logWF = logWF.WithFields(log.Fields{
		"fontSize":      memeFile.FontSize,
		"fileName":      memeFile.Filename,
		"fileExtension": memeFile.Extension,
	})

	// Generating a "predictable output file name" so that we can cache images
	// generated.
	h := sha1.New()
	h.Write([]byte(kind + text))
	hashFileName := base64.URLEncoding.EncodeToString(h.Sum(nil))
	hashFileName = fmt.Sprintf("./%s.%s", hashFileName, memeFile.Extension)
	outputFile := path.Join(OutputTempDir, hashFileName)

	if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
		logWF.Debugf("File already exists: %s", outputFile)
		urlChan <- hashFileName
		return
	}

	// Ensuring that the font size is the one specified after the  first dot
	memeCfg.FontSize = memeFile.FontSize
	meme := &gomeme.Meme{
		Config: memeCfg,
	}

	mof, err := ioutil.ReadFile(memeFile.Filename)
	if err != nil {
		logWF.Warn("Error Opening %s: %s", memeFile.Filename, err.Error())
		return
	}

	meme.Memeable, err = detectFileType(mof)
	if err != nil {
		logWF.Warnf("Error making meme meemable: %s", err.Error())
		return
	}

	output, err := os.Create(outputFile)
	//ioutil.TempFile(os.TempDir(), "meme-")
	if err != nil {
		logWF.Warnf("Error creating temporary file: %s", err.Error())
		return
	}

	err = meme.Write(output)
	if err != nil {
		logWF.Warnf("Unable to create meme %s", err.Error())
		return
	}

	logWF.Debugf("File available: http://localhost:8001/static/%s", hashFileName)
	logWF.Infof("Generated")
	urlChan <- hashFileName
}

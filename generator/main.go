package generator

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
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

func init() {
	var err error
	OutputTempDir, err = os.MkdirTemp(os.TempDir(), "meme-generator")
	if err != nil {
		log.Errorf("Error creating temp dir: %s", err.Error())
	}
}

func SetPkgConfig(ch chan string, assetPath string) {
	urlChan = ch
	assetsDirPath = assetPath

	MemeFiles = &[]string{}
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

	// TODO: DO IT
	logWF.Infof(`Generating...`)

	memeCfg := gomeme.NewConfig()
	memeCfg.BottomText = text

	var fontSizeStr string
	var fileName string
	var fileExtension string
	for _, mfn := range *MemeFiles {
		if strings.HasPrefix(mfn, kind) {
			fileName = fmt.Sprintf("assets/%s", mfn)
			// Get the font size from the fileName
			fontSizeStr = strings.Split(mfn, ".")[1]
			fileExtension = strings.Split(mfn, ".")[2]
			break
		}
	}

	// Updating the logs with relevant extra Fields
	logWF = logWF.WithFields(log.Fields{
		"fontSize":      fontSizeStr,
		"fileName":      fileName,
		"fileExtension": fileExtension,
	})

	// Generating a "predictable output file name" so that we can cache images
	// generated.
	h := sha1.New()
	h.Write([]byte(kind + text))
	hashFileName := base64.URLEncoding.EncodeToString(h.Sum(nil))
	hashFileName = fmt.Sprintf("./%s.%s", hashFileName, fileExtension)
	outputFile := path.Join(OutputTempDir, hashFileName)

	if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
		logWF.Debugf("File already exists: %s", outputFile)
		urlChan <- hashFileName
		return
	}

	// Ensuring that the font size is the one specified after the  first dot
	fontSize, err := strconv.ParseFloat(fontSizeStr, 64)
	if err != nil {
		logWF.Warnf("Error defining the fontSize: %s", err.Error())
		return
	}
	memeCfg.FontSize = fontSize
	meme := &gomeme.Meme{
		Config: memeCfg,
	}

	memeFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		logWF.Warn("Error Opening %s: %s", fileName, err.Error())
		return
	}

	meme.Memeable, err = detectFileType(memeFile)
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

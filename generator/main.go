package generator

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

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
		log.Fatalf("Error creating temp dir: %s", err.Error())
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
	// TODO: DO IT
	log.Printf(`[%s] %s`, kind, text)

	memeCfg := gomeme.NewConfig()
	memeCfg.BottomText = text

	var fontSizeStr string
	fileName := "assets/wink.jpg"
	fileExtension := "jpg"
	for _, mfn := range *MemeFiles {
		if strings.HasPrefix(mfn, kind) {
			fileName = fmt.Sprintf("assets/%s", mfn)
			// Get the font size from the fileName
			fontSizeStr = strings.Split(mfn, ".")[1]
			fileExtension = strings.Split(mfn, ".")[2]
			break
		}
	}

	// Generating a "predictable output file name" so that we can cache images
	// generated.
	h := sha1.New()
	h.Write([]byte(kind + text))
	hashFileName := base64.URLEncoding.EncodeToString(h.Sum(nil))
	hashFileName = fmt.Sprintf("./%s.%s", hashFileName, fileExtension)
	outputFile := path.Join(OutputTempDir, hashFileName)

	if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
		log.Printf("File already exists: %s", outputFile)
		urlChan <- hashFileName
		return
	}

	// Ensuring that the font size is the one specified after the  first dot
	fontSize, err := strconv.ParseFloat(fontSizeStr, 64)
	if err != nil {
		log.Printf("Error defining the fontSize: %s", err.Error())
		return
	}
	memeCfg.FontSize = fontSize
	meme := &gomeme.Meme{
		Config: memeCfg,
	}

	memeFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Error Opening %s: %s", fileName, err.Error())
		return
	}

	meme.Memeable, err = detectFileType(memeFile)
	if err != nil {
		log.Printf("Error making meme meemable: %s", err.Error())
		return
	}

	output, err := os.Create(outputFile)
	//ioutil.TempFile(os.TempDir(), "meme-")
	if err != nil {
		log.Printf("Error creating temporary file: %s", err.Error())
		return
	}

	err = meme.Write(output)
	if err != nil {
		log.Printf("Unable to create meme %s", err.Error())
		return
	}

	log.Printf("File available: http://localhost:8001/static/%s", hashFileName)
	urlChan <- hashFileName
}

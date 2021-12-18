package generator

import (
	"bytes"
	"errors"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/chai2010/webp"
	log "github.com/sirupsen/logrus"
	"gitlab.com/koalalorenzo/gomeme"
)

// PermittedFileExtension contaisn the extension that are allowed to use
var PermittedFileExtension = []string{
	".gif",
	".jpg",
	".jpeg",
	".png",
	".webp",
}

// MemeFiles these are the file names of the memes available from the bot
var MemeFiles []*MemeFile

type MemeFile struct {
	Kind      string
	Filename  string
	FontSize  float64
	Extension string
}

var (
	// This is the path where we store all the images.
	AssetsDirPath string
)

func getMemesFilesData() (l []*MemeFile, err error) {
	logWF := log.WithFields(log.Fields{
		"f": "generator.getMemesFilesData",
	})

	l = []*MemeFile{}
	files, err := os.ReadDir(AssetsDirPath)
	if err != nil {
		return []*MemeFile{}, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		mfn := file.Name()

		// Ensure that we check only images and gifs formats
		isExtValid := false
		for _, ext := range PermittedFileExtension {
			if strings.HasSuffix(mfn, ext) {
				isExtValid = true
				break
			}
		}

		// Go to the next file if the extension is not allowed
		if !isExtValid {
			continue
		}

		logWF = logWF.WithField("file", mfn)

		mfsp := strings.Split(mfn, ".")

		if len(mfsp) < 3 {
			logWF.Debugf("Error, file has wrong fomrat. Expecting <kind>.<font-size>.<extension>")
			continue
		}

		fontSize, err := strconv.ParseFloat(mfsp[1], 64)
		if err != nil {
			logWF.Warnf("Error defining the fontSize: %s", err.Error())
			continue
		}

		m := &MemeFile{
			Filename:  fmt.Sprintf("assets/%s", mfn),
			Kind:      mfsp[0],
			FontSize:  fontSize,
			Extension: mfsp[2],
		}
		l = append(l, m)
	}
	return l, nil
}

func tickerFilesLiveLoade() {
	logWF := log.WithFields(log.Fields{
		"f":      "generator.tickerFilesLiveLoade",
		"assets": AssetsDirPath,
	})

	mfs, err := getMemesFilesData()
	if err != nil {
		logWF.Errorf("Erorr loading files: %s", err.Error())
	}

	MemeFiles = mfs
}

func detectFileType(in []byte) (n gomeme.Memeable, err error) {
	contentType := http.DetectContentType(in)
	buff := bytes.NewBuffer(in)

	switch contentType {
	case "image/gif":
		g, err := gif.DecodeAll(buff)
		if err != nil {
			return nil, err
		}
		return gomeme.GIF{g}, nil

	case "image/jpeg":
		j, err := jpeg.Decode(buff)
		if err != nil {
			return nil, err
		}
		return gomeme.JPEG{j}, nil

	case "image/png":
		p, err := png.Decode(buff)
		if err != nil {
			return nil, err
		}
		return gomeme.PNG{p}, nil

	case "image/webp":
		p, err := webp.Decode(buff)
		if err != nil {
			return nil, err
		}
		return gomeme.WebP{p}, nil
	}

	return nil, errors.New("Unable to identify file type")
}

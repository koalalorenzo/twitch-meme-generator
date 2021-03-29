package generator

import (
	"bytes"
	"errors"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"

	"github.com/jpoz/gomeme"
)

func getMemesFilesData() (l []string, err error) {
	l = []string{}
	files, err := os.ReadDir("assets")
	if err != nil {
		return []string{}, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		l = append(l, file.Name())
	}
	return l, nil
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
	}
	return nil, errors.New("Unable to identify file type")
}

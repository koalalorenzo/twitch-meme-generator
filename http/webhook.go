package http

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"gitlab.com/koalalorenzo/twitch-meme-generator/generator"
)

type WebHookRequest struct {
	Kind string `json:"kind"`
	Text string `json:"text"`
}

func serveWebHook(w http.ResponseWriter, r *http.Request) {
	logWF := log.WithFields(log.Fields{
		"f":          "http.serveWebHook",
		"RemoteAddr": r.RemoteAddr,
		// "Host":       r.Host,
		"UserAgent": r.UserAgent(),
		"URI":       r.RequestURI,
		"Method":    r.Method,
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Making sure this uses the right method
	if r.Method != http.MethodPost {
		logWF.Warn("Got Wrong method (Not POSt")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"methods accept "}`))
		return
	}

	// Decoding the json value in the struct
	reqMeme := WebHookRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqMeme); err != nil {
		logWF.Warnf("Error decoding input:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Bad json format"}`))
		return
	}

	// Generate the meme
	logWF.Debugf("Received request via HTTP webhook")
	go generator.GenerateMeme(reqMeme.Kind, reqMeme.Text)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"meme queued"}`))
}

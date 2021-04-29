package http

import (
	"net/http"

	"gitlab.com/koalalorenzo/twitch-meme-generator/twitch"
)

func healthLiveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}

func healthReady(w http.ResponseWriter, r *http.Request) {
	if twitch.ClientIsConnected {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte{})
}

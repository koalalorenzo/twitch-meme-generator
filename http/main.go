package http

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"

	"gitlab.com/koalalorenzo/twitch-meme-generator/generator"
)

var (
	urlChan chan string
)

func init() {
	homeTempl = template.Must(template.New("").Parse(homeHTML))
}

func SetPkgConfig(ch chan string, displayTime time.Duration) {
	urlChan = ch
	displayTimePeriod = displayTime
}

func StartServer(addr string) {
	r := mux.NewRouter()

	sfHandler := http.FileServer(http.Dir(generator.OutputTempDir))
	staticHandler := http.StripPrefix("/static/", sfHandler)
	r.PathPrefix("/static/").Handler(staticHandler)
	r.PathPrefix("/ws").HandlerFunc(serveWs)
	r.PathPrefix("/").HandlerFunc(serveHome)

	srv := &http.Server{
		Handler: r,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}

	log.Printf("Using http://%s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
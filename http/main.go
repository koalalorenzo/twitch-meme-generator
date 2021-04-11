package http

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"

	"gitlab.com/koalalorenzo/twitch-meme-generator/generator"
)

type Config struct {
	MainChannel       chan string
	DisplayTimePeriod time.Duration
	Webhook           struct {
		Enabled  bool
		Username string
		Password string
	}
}

var conf = &Config{}

func init() {
	// Set defaults config
	conf.Webhook.Enabled = true
	conf.DisplayTimePeriod = 10 * time.Second
}

func SetPkgConfig(c *Config) {
	conf = c
}

func StartServer(addr string) {
	logWF := log.WithFields(log.Fields{
		"f":    "http.StartServer",
		"addr": fmt.Sprintf("http://%s", addr),
	})

	// Prepare a router for Webhook
	whr := mux.NewRouter()
	if conf.Webhook.Enabled {
		logWF.Debug("WebHook enabled")
		whr.Methods(http.MethodPost).HandlerFunc(serveWebHook)
	}

	// If Username & Password were passed then check basic auth against them
	if conf.Webhook.Username != "" && conf.Webhook.Password != "" {
		logWF.Debug("WebHook has basic http authentication enabled")
		whr.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user, pass, ok := r.BasicAuth()
				if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(conf.Webhook.Username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(conf.Webhook.Password)) != 1 {
					w.Header().Set("WWW-Authenticate", `Basic realm=""`)
					http.Error(w, "Unauthorised", http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r)
			})
		})
	} else {
		logWF.Warn("WebHook does not have basic authentication enabled")
	}

	// setting main router
	r := mux.NewRouter()

	sfHandler := http.FileServer(http.Dir(generator.OutputTempDir))
	staticHandler := http.StripPrefix("/static/", sfHandler)
	r.PathPrefix("/static/").Handler(staticHandler)
	r.PathPrefix("/ws").HandlerFunc(serveWs)
	r.PathPrefix("/wh").Handler(whr)
	r.PathPrefix("/").HandlerFunc(serveHome)

	srv := &http.Server{
		Handler: r,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}

	go channelPipe(conf.MainChannel)

	logWF.Infof("Starting HTTP Server")
	log.Fatal(srv.ListenAndServe())
}

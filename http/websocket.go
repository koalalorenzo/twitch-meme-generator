package http

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, msgChan chan string) {
	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()

	for {
		select {
		case url := <-msgChan:
			// Send instantly the new image
			ws.SetWriteDeadline(time.Now().Add(writeWait))

			var err error
			err = ws.WriteMessage(websocket.TextMessage, []byte(url))
			if err != nil {
				return
			}

			// This should not make possible to have other memes untill this one
			// disappears.
			time.Sleep(conf.DisplayTimePeriod)
			// Sending an empty message to make the image disappear
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err = ws.WriteMessage(websocket.TextMessage, []byte{})
			if err != nil {
				return
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	logWF := log.WithFields(log.Fields{
		"f":          "http.serveWs",
		"RemoteAddr": r.RemoteAddr,
		"URI":        r.RequestURI,
	})

	logWF.Infof("New WebSocket connection")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logWF.Warn(err)
		}
		return
	}

	// Add a new channel that will receive the memes for this web socket.
	wsChannel := make(chan string, 5)
	WSListnerChannels = append(WSListnerChannels, wsChannel)

	go writer(ws, wsChannel)
	reader(ws)
}

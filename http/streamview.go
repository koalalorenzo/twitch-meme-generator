package http

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
)

var streamViewTempl *template.Template

func init() {
	// Prepare the HTML template
	streamViewTempl = template.Must(template.New("").Parse(streamViewHTML))
}

func serveStreamView(w http.ResponseWriter, r *http.Request) {
	logWF := log.WithFields(log.Fields{
		"f":          "http.serveStreamView",
		"RemoteAddr": r.RemoteAddr,
		// "Host":       r.Host,
		"UserAgent": r.UserAgent(),
		"URI":       r.RequestURI,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Preparing values for our HTML template
	var v = struct {
		Host    string
		LastMod string
	}{
		r.Host,
		strconv.FormatInt(time.Now().UnixNano(), 16),
	}

	logWF.Infof("")
	streamViewTempl.Execute(w, &v)
}

const streamViewHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>koalalorenzo's Twitch Meme Generator</title>
    </head>
		<style>
			body {
				background-color: transparent;
				background-repeat: no-repeat;
				background-position: center;
				background-size: contain;
				transition: background-image 1s ease-in-out;

				width: 100%;
				height: 100%;
				margin: 0;
				padding: 0;
			}
		</style>
    <body>
        <script type="text/javascript">
            (function() {
								var conn;
							  var body = document.body;
								body.style.height = window.innerHeight - 50 + "px";
								body.style.height = window.innerHeight - 50 + "px";

								function setImage(url) {
									// Clear the image if empty
									if(url === "") {
										body.style.backgroundImage = "";
										return;
									}
									body.style.backgroundImage = 'url('+"http://{{.Host}}/static/" + url +')';
								}

								function startWebSocket() {
									console.log("Starting a new connection");
									// Clean the image in case there is one showing...
									setImage("");
									var conn = new WebSocket("ws://{{.Host}}/ws?lastMod={{.LastMod}}");
									
									conn.onclose = function(evt) {
										console.log("Connection closed... Reconnecting...")
										setTimeout(function(){
											startWebSocket();
										}, 5000);
									}
	
									conn.onmessage = function(evt) {
										console.log("Received new Image filename:", evt.data);
										setImage(evt.data);
									}
								}

								startWebSocket()
            })();
        </script>
    </body>
</html>
`
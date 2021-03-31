package http

import (
	"net/http"
	"strconv"
	"text/template"
	"time"
)

var homeTempl *template.Template

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Preparing values for our HTML template
	var v = struct {
		Host    string
		LastMod string
	}{
		r.Host,
		strconv.FormatInt(time.Now().UnixNano(), 16),
	}

	homeTempl.Execute(w, &v)
}

const homeHTML = `<!DOCTYPE html>
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

								function startWebSocket() {
									console.log("starting the baby");
									var conn = new WebSocket("ws://{{.Host}}/ws?lastMod={{.LastMod}}");
									
									conn.onclose = function(evt) {
										setTimeout(function(){
											startWebSocket();
										}, 5000);
									}
	
									conn.onmessage = function(evt) {
											console.log("received message:", evt.data);
											if(evt.data === "") {
												body.style.backgroundImage = "";
												return;
											}
											body.style.backgroundImage = 'url('+"http://{{.Host}}/static/" + evt.data +')'
									}
								}

								startWebSocket()
            })();
        </script>
    </body>
</html>
`

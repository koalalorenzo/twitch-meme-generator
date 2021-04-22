package http

import (
	"net/http"
	"text/template"

	log "github.com/sirupsen/logrus"

	"gitlab.com/koalalorenzo/twitch-meme-generator/generator"
)

var listTempl *template.Template

func init() {
	// Prepare the HTML template
	listTempl = template.Must(template.New("").Parse(listHTML))
}

func serveListMeme(w http.ResponseWriter, r *http.Request) {
	logWF := log.WithFields(log.Fields{
		"f":          "http.serveListMeme",
		"RemoteAddr": r.RemoteAddr,
		"UserAgent":  r.UserAgent(),
		"URI":        r.RequestURI,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Preparing values for our HTML template
	var v = struct {
		ChannelName string
		Images      []*generator.MemeFile
	}{
		conf.ChannelName,
		generator.MemeFiles,
	}

	logWF.Infof("")
	listTempl.Execute(w, &v)
}

const listHTML = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <title>koalalorenzo's Twitch Meme Generator</title>
        <link href="//cdn.muicss.com/mui-0.10.3/css/mui.min.css" rel="stylesheet" type="text/css" />
        <script src="//cdn.muicss.com/mui-0.10.3/js/mui.min.js"></script>
        </head>
    <body>
      <div class="mui-container">
        <div class="mui-row">
          <div class="mui-col-md-12">
            <h1>{{ .ChannelName }} Twitch Meme Generator Bot</h1>
            <p>
              This page will show the images and their <code>NAME</code> to use when
              generating a new <i>Meme</i> on the
              <a href="http://twitch.tv/{{ .ChannelName }}">live stream</a>.
              To generate a new meme, you need to write in the channel something
              similar to the following text:
            </p>
            <p>
              <code>!meme NAME WRITE YOUR SENTENCE HERE</code>
            </p>
            <p>
              Writing this in the chat channel will generate a new image of kind
              <code>NAME</code>,inserting the text
              <code>WRITE YOUR SENTENCE HERE</code> at the bottom  of the picture.
            </p>
            <p>
              <a class="mui-btn mui-btn--raised mui-btn--primary" href="http://twitch.tv/{{ .ChannelName }}">
                Watch {{ .ChannelName }}
              </a>
              <a class="mui-btn mui-btn--raised mui-btn--accent" href="https://gitlab.com/koalalorenzo/twitch-meme-generator">
                Add this to your channel
              </a>
            </p>
          </div>
        <div>
        <div class="mui-row">
          <div class="mui-col-md-12">
            <h2>List of Meme names</h2>
          <div>
          {{ range .Images }}
            <div class="mui-col-md-3 mui-panel">
              <img width="100%" src="/{{ .Filename }}"><br/>
              <code>{{ .Kind }}</code>
            </div>
          {{ end }}
        </div>
      </div>

    </body>
</html>
`

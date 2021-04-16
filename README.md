# Koalalorenzo's Twitch Meme Bot
This Twitch bot is capable to listen for `!meme` commands in a Twitch chat and
generate on the fly memes to display as an overlay for OBS. 

If you like it, please follow me on 
[Twitch channel](https://twitch.tv/koalalorenzo). Original source code 
[available here](https://gitlab.com/koalalorenzo/twitch-meme-generator)

See it live [here](https://clips.twitch.tv/VibrantHotZucchiniAsianGlow-iTtCXFtzvn8cBljd)

## Building
To build the code into a portable binary you need the following things installed

- [GNU Make](https://www.gnu.org/software/make/) 
  (`brew install make`, but it is often pre-installed)
- [Go 1.16](https://go.dev) 
  (`brew install golang`)

After that you can run:

```bash
make clean build
```

and you will find all the files needed in the `build` directory. You can also
cross compile for Raspberry Pi by running:

```bash
make clean build -e BUILD_TARGET=rpi
```

## Usage
### Start the Server locally
To start the server locally you can use the built binary as in the following
exampoe:

```
./build/koalalorenzo-meme-generator --channel koalalorenzo
```

If you need more help feel free to run the command with `--help` flag.
Have a look at it because you can customize things like display time and 
host address (default to: `0.0.0.0:8000`).

### OBS Setup
Once the server is listening and ready you can load a new Browser Source in
OBS. The URL of the source should the one printed in the logs. 

The size of the window can be arbitrary based on the custom images. I suggest
to use:

- Width: 800
- Height: 600

By default this url will be [http://localhost:8000/](http://localhost:8000/).
If you are running on a different platform like Heroku or DigitalOcean please
use the URL provided by those platforms.

Remember to have dibled the following options: 
  - Refresh browser when scene becomes active
  - Shutdown source when not visible

The reason is this that this uses WebSockets and it is better to no have 
_multiple listners_ at the same time (aka: reuse the same source by 
_Copy Reference_). In case of issues the web page should reconnect automagically

### Adding files
You can add files any time into the `assets` directory (please note that the
path of the directory can be customized using flags). 

**IMPORTANT**: The file name convention is the following:

```xml
<meme>.<font-size>.<extension>
```

where:
- `meme` is the 1-word used after the `!meme` chat command to slect the image
- `font-size` is the size of the text that will be in the image, this is 
  a value that depends on the size of the image. Test it!
- `extension` is the image format. 
  It can be one of the following: `gif`, `png`, `jpeg` or `jpg`.

You can test the meme by running `generate` subcommand:

```bash
./build/koalalorenzo-meme-generator generate deal Dogs are cool, DEAL WITH IT
```

### Chat usage
The server  is connected to the Twitch Channel (IRC channel) using an anonymous
client. Therefore is not capable of announcing himself in the chat.

If the server is running you should be able to use the bot by writing a message
with the following syntax:

```xml
!meme <meme> <phrase to write...>
```

For example:

```xml
!meme pathetic Humans? Pathetic servants
```

Will display in the Browser source for a few seconds the following image:

![Picture of a cat master](example.jpg)

### WebHook usage
It is possible to set up a custom webhook that will queue memes and show them.
This is useful if you want to hook up some other tools or manually generate 
memes without using the chat (ex: From Apple Shortcuts or from `curl` for 
testing).

By default the webhook is disabled, to enable it you can pass the flag 
`--webhook-enable` to `true` or set the env variable `KTMG_WEBHOOK_ENABLE=true`.
This endpoint can be _protected_ by basic http authentication. Please use
`--help` flag to check all the options availalbe.

For example if I set the following env variables and run the app:

```bash
export KTMG_LOGLEVEL=debug
export KTMG_WEBHOOK_ENABLE=true
export KTMG_WEBHOOK_USERNAME=koalalorenzo
export KTMG_WEBHOOK_USERNAME=SuperSecretPassword
./build/koalalorenzo-meme-generator
```

I will then be able to queue memes to be generated on the fly by running 
http POST request to the endpoint `/wh`, like this:

```bash
curl http://127.0.0.1:8000/wh \
  -X POST --data '{ "kind":"pathetic", "text":"Humans? Pathetic servants"}' \
  -u koalalorenzo:SuperSecretPassword
```

## Deploy on Heroku

To quickly deploy on Heroku you can click on the following button and follow
the instructions:

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/koalalorenzo/twitch-meme-generator/tree/main)

Follow the instructions and remember to customize the twitch channel and display
time as env variables.

**Important**: Please note that the name of the application will impact on the
URL that will be used in the OBS Browser Source. Do not share this URL/address
with anybody and try to keep it hard to find. (_hint_: try adding random
characters in the name might help)

**Important**: Due to security restrictions on some browser, using HTTP instead
of HTTPs in the OBS Browser Source URL will fix connectivity issues.

## Docker usage
It is possible to run Lorenzo's Twitch Bot meme generator using Docker.
The docker image is downloadable from 
[the project GitLab Docker Registry](https://gitlab.com/koalalorenzo/twitch-meme-generator/container_registry)

You can use both env variables as well as command line to configure the app:

```bash
docker run -p 8005:8005 registry.gitlab.com/koalalorenzo/twitch-meme-generator:latest
```

You can also mount the temporary directory (under path `/tmp`) to save the cache
of generated memes, this will speed up future requests.

You can mount the `/assets` path to a custom directory path where you can drop
the images that will be used as base for the memes.

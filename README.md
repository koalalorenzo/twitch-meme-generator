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

Remember to enable the following options: 
  - Refresh browser when scene becomes active
  - Shutdown source when not visible

The reason is this that this uses WebSockets and it is better to no have 
_multiple listners_ at the same time (aka: reuse the same source by 
_Copy Reference_)

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


## Deploy on Heroku
(WIP)

# Twitch Meme Bot

This project contains the source code of a Twitch Bot capable of generating
memes and displaying them on a stream using a browser source.

Please follow the development live on my 
[Twitch channel](https://twitch.tv/koalalorenzo). Original source code 
[available here](https://gitlab.com/koalalorenzo/twitch-meme-generator)

See it live [here](https://clips.twitch.tv/VibrantHotZucchiniAsianGlow-iTtCXFtzvn8cBljd)

What this Bot does:

* Connect to the twitch channel and listen to `!meme` commands
* It will generate a PNG/JPEG Image with custom text
* It will update an HTTP endpoint to show the meme for a few seconds

## Building

To build the code into a portable binary you need the following things installed

- GNU Make (`brew install make`, but it is often pre-installed)
- Go 1.16

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
(WIP)

### OBS Setup
(WIP)

### Adding files
(WIP)
BUILD_TARGET ?=

ifeq (${BUILD_TARGET},rpi)
GOARCH := arm
GOOS := linux
GOARM=7
endif

.EXPORT_ALL_VARIABLES:

build:
	go build -o build/twitch-meme-bot 
	cp -R ./assets ./build/assets

run:
	go run ./*.go

clean:
	rm -rf build
.PHONY: clean
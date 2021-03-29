build:
	go build -o build/twitch-meme-bot 

run:
	go run ./*.go

clean:
	rm -rf build
.PHONY: clean
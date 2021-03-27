build:
	go build -o build/twitch-meme-bot 

run:
	go run

clean:
	rm -rf build
.PHONY: clean
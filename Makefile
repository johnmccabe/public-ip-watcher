.PHONY: all clean build run image

all: clean build

clean:
	rm -f db/ipwatcher.db
	rm -f ipwatcher

build:
	go build -o ipwatcher .

run:
	go run main.go

image:
	docker build -t johnmccabe/ipwatcher:dev .
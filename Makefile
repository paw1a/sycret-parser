build:
	go mod download && go build -o ./.bin/app ./cmd/app/main.go

run: build
	./.bin/app

build-linux:
	GOOS=linux GOARCH=arm64 go build -o linux-app.exe cmd/app/main.go

test:
	go test

clean:
	rm -rf .bin

.DEFAULT_GOAL := run

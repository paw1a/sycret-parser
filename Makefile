build:
	go mod download && go build -o ./.bin/app ./cmd/app/main.go

run: build
	./.bin/app

test:
	go test ./... -cover

clean:
	rm -rf .bin

.DEFAULT_GOAL := run

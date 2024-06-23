build:
	mkdir -p bin && go build -o bin/gogrep cmd/gogrep/main.go

.DEFAULT_GOAL := build

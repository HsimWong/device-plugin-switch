#.PHONY: build

all:
	go build -o build/device-plugin-switch app/device-plugin-switch.go

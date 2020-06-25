.PHONY: build
build:
	go build -v ./cmd/gsheet-crm

.DEFAULT_GOAL := build

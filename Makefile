.PHONY: build

rpm: build
	nfpm pkg --packager rpm --target .
deb: build
	nfpm pkg --packager deb --target .

build:
	go build -v ./cmd/gsheet-crm

.DEFAULT_GOAL := build

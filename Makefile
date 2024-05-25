#
# Makefile
#
VERSION = snapshot
GHRFLAGS =
.PHONY: build release

default: build

build:
	go build -o pkg/$(VERSION) .

release:
	gh release -t "svq - service query engine" -u sbyware $(GHRFLAGS) svdb pkg/$(VERSION)

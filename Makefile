.PHONY: build, install, clean

build:
	go build .

install:
	go install .

clean:
	go clean -i ./...

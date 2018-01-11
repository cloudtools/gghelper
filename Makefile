.PHONY: build clean release test

build:
	go build .
	go build ./cmd/gghelper

clean:
	rm -f *-setup.tar.gz
	rm -f gghelper gghelper_darwin gghelper_linux

release:
	GOOS=darwin GOARCH=amd64 go build -o gghelper_darwin ./cmd/gghelper
	GOOS=linux GOARCH=amd64 go build -o gghelper_linux ./cmd/gghelper

test:
	go test

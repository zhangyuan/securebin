build:
	go build -ldflags "-s -w"

clean:
	rm -rf securebin
	rm -rf bin/securebin-*

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/securebin_darwin-amd64 -ldflags "-s -w"
	GOOS=darwin GOARCH=arm64 go build -o bin/securebin_darwin-arm64 -ldflags "-s -w"

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/securebin_linux-amd64 -ldflags "-s -w"

build-windows:
    GOOS=windows GOARCH=amd64 go build -o bin/securebin_windows-amd64 -ldflags "-s -w"

build-all: clean build-macos build-linux build-windows

compress-linux:
	upx ./bin/securebin_linux*

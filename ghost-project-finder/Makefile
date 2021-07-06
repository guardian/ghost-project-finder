all: ghost-project-finder

clean:
	rm -f ghost-project-finder
	rm -f ghost-project-finder-mac.zip

ghost-project-finder: main.go scanner.go sender.go ../common/datapacket.go
	GOOS=darwin GOARCH=amd64 go build -o ghost-project-finder

bundle: ghost-project-finder ghost-project-finder-downloads.plist ghost-project-finder-root.plist mac-install.sh
	zip ghost-project-finder-mac.zip ghost-project-finder ghost-project-finder-downloads.plist ghost-project-finder-root.plist mac-install.sh
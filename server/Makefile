all: server.linux64

server.linux64: main.go inputhandler.go helpers.go healthcheck.go ../common/datapacket.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server.linux64

clean:
	rm -f server server.linux64
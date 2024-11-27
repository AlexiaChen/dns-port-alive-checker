build:
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o bin/dns_check main.go
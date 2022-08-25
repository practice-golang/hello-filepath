build:
	go build -ldflags "-w -s -H windowsgui" -o bin/

test:
	go test -v
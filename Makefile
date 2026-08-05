.PHONY: run build test clean

APP_NAME=erm-dokter
MAIN_PATH=cmd/api/main.go

run:
	go run $(MAIN_PATH)

build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

test:
	go test -v ./...

clean:
	rm -rf bin/

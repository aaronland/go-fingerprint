GOMOD=vendor

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/convert cmd/convert/main.go
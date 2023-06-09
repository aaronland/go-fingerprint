GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/convert cmd/convert/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/pdf cmd/pdf/main.go

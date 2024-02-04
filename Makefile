all: test client server

test: reports
	go test -v -cover ./... -coverprofile=reports/report.cover

server: client
	go build -race -o bin/server cmd/server/*.go

demo:
	go run cmd/demo/main.go

coverage:
	go tool cover -html=reports/report.cover -o reports/coverage.html
	@xdg-open reports/coverage.html 2>/dev/null

build:
	go build -o bin/main main.go

client:
	(cd client && npm run build)

clean:
	rm -f bin/*
	rm -f reports/*
	rm -rf client/dist/*

.PHONY: test build clean reports coverage client

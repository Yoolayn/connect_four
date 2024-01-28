test: reports
	go test -v -cover ./... -coverprofile=reports/report.cover

reports:
	mkdir reports -p

server:
	go build -race -o bin/server cmd/server/*.go

demo:
	go run cmd/demo/main.go

coverage:
	go tool cover -html=reports/report.cover -o reports/coverage.html
	@xdg-open reports/coverage.html 2>/dev/null

build:
	go build -o bin/main main.go

clean:
	rm -f bin/*
	rm -f reports/*

.PHONY: test build clean reports coverage

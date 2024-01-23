test: reports
	go test -v -cover ./... -coverprofile=reports/report.cover

reports:
	mkdir reports -p

coverage:
	go tool cover -html=reports/report.cover -o reports/coverage.html
	@xdg-open reports/coverage.html 2>/dev/null

build:
	go build -o bin/main main.go

clean:
	rm -f bin/*
	rm -f reports/*

.PHONY: test build clean reports coverage

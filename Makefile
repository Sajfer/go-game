build:
	go build -o bin/go-game main.go

test:
	go fmt $(go list ./... | grep -v /vendor/)
	go vet $(go list ./... | grep -v /vendor/)
	go test -race $(go list ./... | grep -v /vendor/)

run:
	go run main.go

compile:
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
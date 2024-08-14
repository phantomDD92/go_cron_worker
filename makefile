run:
	go run main.go

tidy:
	go mod tidy

update-all:
	go get -u ./...

install-all:
	go get ./...
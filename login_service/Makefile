BINARY_NAME = login_service

.PHONY: build
build:
	go build -o $(BINARY_NAME) main.go

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
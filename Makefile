.PHONY: build test clean run

build:
	go build -o wizask ./...

test:
	go test ./...

clean:
	rm -f wizask

run:
	go run . $(ARGS)

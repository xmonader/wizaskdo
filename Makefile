.PHONY: build test clean run

build:
	go build -o wizask .
	go build -o wizdo ./cmd/wizdo

test:
	go test ./pkg/...

clean:
	rm -f wizask wizdo

run:
	go run . $(ARGS)

run-do:
	go run ./cmd/wizdo $(ARGS)

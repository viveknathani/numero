build:
	go build -o ./bin/numero main.go

test:
	go test -v ./...

run-dev:
	export LOG_LEVEL=DEBUG && make build && ./bin/numero

run:
	export LOG_LEVEL=INFO && make build && ./bin/numero

format:
	golines -w .

bench:
	go run benchmark/main.go
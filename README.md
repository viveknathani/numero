# numero

[![Go Report Card](https://goreportcard.com/badge/github.com/viveknathani/numero)](https://goreportcard.com/report/github.com/viveknathani/numero) ![test](https://github.com/viveknathani/numero/actions/workflows/test.yaml/badge.svg)
![Go Version](https://img.shields.io/github/go-mod/go-version/viveknathani/numero)
![Last Commit](https://img.shields.io/github/last-commit/viveknathani/numero)

`numero` is software for parsing and evaluating mathematical expressions. It is available as a library and as a web service.

## usage

The library can be used as follows.


Download it:
```bash
go get -u github.com/viveknathani/numero
```

Import it:
```go
import "github.com/viveknathani/numero/nparser"
``` 

Simple example:
```go
expression := "sin(max(2, 3333))"
parser := nparser.New(expression)

result, err := parser.Run()
```

Example with variables:
```go
expression := "x + y"
parser := New(expression)
parser.SetVariable("x", 2)
parser.SetVariable("y", 45)
result, err := parser.Run()
```

The web service can be consumed as follows:

```bash
curl --request POST \
  --url https://numero.vivekn.dev/api/v1/eval \
  --header "Content-Type: application/json" \
  --data '{"expression":"x + sin(max(2, 333))","variables":{"x":100}}'
```

## documentation

**Supported functions**

- `sin`
- `cos`
- `tan`
- `log`
- `ln`
- `sqrt`
- `max`
- `min`

**Supported operators**

- `+`
- `-`
- `*`
- `/`
- `^`

**API**

`POST /api/v1/eval`

Request body parameters (JSON):

- `expression`: the expression to evaluate
- `variables`: a map of variable names to values

Response body:

```json
{
  "data": {
    "result": 99.99117883388611
  },
  "message": "success"
}
```

### benchmarks

This runs a load test for 20 seconds. The test can be found [here](https://github.com/viveknathani/numero/blob/master/benchmark/main.go). The tests were run on a 2021 Macbook Pro with an M1 chip.

```bash
make bench
```
| metric       | value   |
|---------------|---------|
| total requests| 30878   |
| successful requests| 30878|
| failed requests | 688    |
| connection errors| 688   |
| invalid responses| 0     |
| average latency | 1.02 ms|
| min latency     | 0.06 ms|
| max latency     | 36.12 ms|
| requests/second | 1543|

## contributing

I am happy to accept pull requests. No hard rules.

To set up the project for development, we have the following system requirements:

1. git
2. go
3. make

```bash
git clone https://github.com/viveknathani/numero.git
cd numero
make build
make test
make run-dev
```

## motivation

This project started as an exercise in doing some recreational programming. I always knew about the [shunting yard algorithm](https://en.wikipedia.org/wiki/Shunting_yard_algorithm) but never really got to implement it. Lately, I've been writing a lot of code in Go and decided to just do this.

## acknowledgements

created by Vivek Nathani ([@viveknathani_](https://twitter.com/viveknathani_)), licensed under the [MIT License](https://github.com/viveknathani/numero/blob/master/LICENSE).

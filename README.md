# Wallex-Go

This library provides a Go client for Wallex Exchange API described [here](https://wallex-docs.github.io).

## Install

```shell
$ go get -u github.com/wallexchange/wallex-go
```

## Getting Started

```go
import wallex "github.com/wallexchange/wallex-go"

func main() {
    client := wallex.New(wallex.ClientOptions{
        APIKey: "xxx|xxxxx",
    })
    ...
}
```

## TODO

- [ ] Add socket.io endpoints

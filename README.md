# Go Namecheap SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/namecheap/go-namecheap-sdk.svg)](https://pkg.go.dev/github.com/namecheap/go-namecheap-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/namecheap/go-namecheap-sdk)](https://goreportcard.com/report/github.com/namecheap/go-namecheap-sdk)

### Getting

```sh
$ go get github.com/namecheap/go-namecheap-sdk/v2
```

### Usage

```go
import (
    "github.com/namecheap/go-namecheap-sdk/v2"
)

client := NewClient(&ClientOptions{
    UserName:   "UserName",
    ApiUser:    "ApiUser",
    ApiKey:     "ApiKey",
    ClientIp:   "10.10.10.10",
    UseSandbox: false,
})
```

### Contributing

You're welcome to post issues and send your pull requests.

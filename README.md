# NSQ URL Parser

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/bryanaustin/nsqparse)

Parser for storing NSQ connection information in url format. Example:

    nsqd://nsq.server:4150/topic/channel

Example of code usage:
```go
package main

import (
  "github.com/bryanaustin/nsqparse"
  "github.com/nsqio/go-nsq"
)

func main() {
  nu, err := nsqparse.Parse("localhost:4150/topic")
  if err != nil {
    // handle error
  }
  nsqconfig := nsq.NewConfig()
  consumer, err = nu.Consumer(nsqconfig)
  if err != nil {
    // handle error
  }
  consumer.AddHandler(...)
  err = nu.ConnectConsumer(consumer)
  if err != nil {
    // handle error
  }
  // running
}
```

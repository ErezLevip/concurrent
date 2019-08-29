# Concurrent

Concurrent is a simple library for concurrent processing while keeping the original order.

## Installation
go get -u "github.com/erezlevip/concurrent"

## Quick Start

Create a processing function and call Process with a channel of the objects that needs to be processed, the processing function and the limit per goroutine.

Pass channel:
```go
limitPerGoroutine := 5
processFunc := func(input interface{}) (interface{}, error) { return json.Marshal(input) }
jsonResults, err := concurrent.Process(objectsToProcessChannel, processFunc, limitPerGoroutine)
```

Pass slice:
```go
limitPerGoroutine := 5
processFunc := func(input interface{}) (interface{}, error) { return json.Marshal(input) }
jsonResults, err := concurrent.ProcessSlice(objectsToProcessSlice, processFunc, limitPerGoroutine)
```
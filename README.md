# Concurrent

concurrent is a simple library for concurrent processing.

## Installation
go get -u "github.com/erezlevip/concurrent"

## Quick Start

create a processing function and call Process with a channel of the objects that needs to be processed, the processing function and the limit per goroutine.
```go
limitPerGoroutine := 5
processFunc := func(input interface{}) (interface{}, error) { return json.Marshal(input) }
jsonResults, err := concurrent.Process(objectsToProcess, processFunc, limitPerGoroutine)
```

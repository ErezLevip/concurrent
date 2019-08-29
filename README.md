# Concurrent

concurrent is a simple library for concurrent processing.

## Installation
go get -u "github.com/erezlevip/concurrent"

## Quick Start

create a processing function and select a limit.
```go
limitPerGoroutine := 5
processFunc := func(input interface{}) (interface{}, error) { return json.Marshal(input) }
jsonResults, err := concurrent.Process(objectsToProcess, processFunc, limitPerGoroutine)
```

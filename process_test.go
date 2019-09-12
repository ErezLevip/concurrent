package concurrent

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	maxIndexSize := 20
	limit := 3
	inputsCh := loadToStuct(maxIndexSize)

	//to json
	jsonResults, err := Process(inputsCh, process, limit)
	if err != nil {
		panic(err)
	}

	jsonCh := loadFromJson(jsonResults, maxIndexSize)

	//from json
	structResults, err := Process(jsonCh, fromJson, limit)
	if err != nil {
		panic(err)
	}
	match := true
	i := 0
	for v := range loadToStuct(maxIndexSize) {
		fmt.Println(structResults[i].(*NumberHolder).Age, "!=", v.(*NumberHolder).Age, structResults[i].(*NumberHolder).Age != v.(*NumberHolder).Age)
		if structResults[i].(*NumberHolder).Age != v.(*NumberHolder).Age {
			match = false
			break
		}
		i++
	}

	assert.True(t, match, "results dont match the original stage")
}
func TestProcess_empty(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	c := make(chan interface{})
	close(c)

	//to json
	jsonResults, err := Process(c, process, 1)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 0, len(jsonResults), "results dont match the original stage")
}

func TestProcessSlice(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	maxIndexSize := 1000
	inputsSlice := loadToStuctSlice(maxIndexSize)

	//to json
	jsonResults, err := ProcessSlice(inputsSlice, process, 1)
	if err != nil {
		panic(err)
	}

	//from json
	structResults, err := ProcessSlice(jsonResults, fromJson, 1)
	if err != nil {
		panic(err)
	}
	match := true
	for i, v := range loadToStuctSlice(maxIndexSize) {
		fmt.Println(structResults[i].(*NumberHolder).Age, "!=", v.(*NumberHolder).Age, structResults[i].(*NumberHolder).Age != v.(*NumberHolder).Age)
		if structResults[i].(*NumberHolder).Age != v.(*NumberHolder).Age {
			match = false
			break
		}
	}
	assert.True(t, match, "results dont match the original stage")
}

func process(obj interface{}) (interface{}, error) {
	time.Sleep(300* time.Millisecond)
	return json.Marshal(obj)
}
func fromJson(obj interface{}) (interface{}, error) {
	n := NumberHolder{}
	if err := json.Unmarshal(obj.([]byte), &n); err != nil {
		return nil, err
	}
	return &n, nil
}

type NumberHolder struct {
	Age int `json:"age"`
}

func loadFromJson(results []interface{}, maxNum int) chan interface{} {
	ch := make(chan interface{}, maxNum)
	for i := 0; i < maxNum; i++ {
		ch <- results[i]
	}
	close(ch)
	return ch
}
func loadToStuct(maxNum int) chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for i := 0; i < maxNum; i++ {
			ch <- &NumberHolder{Age: i}
		}
	}()
	return ch
}
func loadToStuctSlice(maxNum int) []interface{} {
	s := make([]interface{}, maxNum)
	for i := 0; i < maxNum; i++ {
		s[i] = &NumberHolder{Age: i}
	}
	return s
}

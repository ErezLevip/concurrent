package concurrent

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestProcess(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	inputsCh := loadToStuct()

	//to json
	jsonResults, err := Process(inputsCh, process, 5)
	if err != nil {
		panic(err)
	}

	jsonCh := loadFromJson(jsonResults)

	//from json
	structResults, err := Process(jsonCh, fromJson, 5)
	if err != nil {
		panic(err)
	}
	match := true
	i := 0
	for v := range inputsCh {
		if structResults[i].(*NumberHolder).Age != v.(*NumberHolder).Age {
			match = false
			break
		}
		i++
	}
	assert.True(t, match, "results dont match the original stage")
}

func TestProcessSlice(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	inputsSlice := loadToStuctSlice()

	//to json
	jsonResults, err := ProcessSlice(inputsSlice, process, 5)
	if err != nil {
		panic(err)
	}

	//from json
	structResults, err := ProcessSlice(jsonResults, fromJson, 5)
	if err != nil {
		panic(err)
	}
	match := true
	for i, v := range inputsSlice {
		if structResults[i].(*NumberHolder).Age != v.(*NumberHolder).Age {
			match = false
			break
		}
	}
	assert.True(t, match, "results dont match the original stage")
}

func process(obj interface{}) (interface{}, error) {
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

func loadFromJson(results []interface{}) chan interface{} {
	maxNum := 1000
	ch := make(chan interface{}, maxNum)
	for i := 0; i < maxNum; i++ {
		ch <- results[i]
	}
	close(ch)
	return ch
}
func loadToStuct() chan interface{} {
	maxNum := 1000
	ch := make(chan interface{}, maxNum)
	for i := 0; i < maxNum; i++ {
		ch <- &NumberHolder{Age: i}
	}
	close(ch)
	return ch
}
func loadToStuctSlice() []interface{} {
	maxNum := 1000
	s := make([]interface{}, maxNum)
	for i := 0; i < maxNum; i++ {
		s[i] = &NumberHolder{Age: i}
	}
	return s
}

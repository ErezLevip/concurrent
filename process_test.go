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
	structResults, err := Process(jsonCh, process, 5)
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

func process(obj interface{}) (interface{}, error) {
	return json.Marshal(obj)
}

type NumberHolder struct {
	Age int `json:"age"`
}

func loadFromJson(results []interface{}) chan interface{} {
	maxNum := 1000
	ch := make(chan interface{}, maxNum)
	for i := 0; i < maxNum; i++ {
		var a NumberHolder
		if err := json.Unmarshal(results[i].([]byte), &a); err != nil {
			panic(err)
		}

		ch <- &a
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

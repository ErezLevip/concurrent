package main

import (
	"encoding/json"
	"fmt"
	"github.com/erezlevip/concurrent"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	inputs := loadToStuct()

	//to
	jsonResults, err := concurrent.Process(inputs, process, 5)
	if err != nil {
		panic(err)
	}

	jsonCh := loadFromJson(jsonResults)
	//from
	structResults, err := concurrent.Process(jsonCh, process, 5)
	if err != nil {
		panic(err)
	}
	i:=0
	for v := range inputs {
		if structResults[i].(*NumberHolder).Age != v.(*NumberHolder).Age{
			fmt.Println("i",i,"does not match")
			return
		}
		i++
	}
	fmt.Println("prefect match")
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

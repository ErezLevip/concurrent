package concurrent

import (
	"errors"
)

type processFunc func(interface{}) (interface{}, error)

const (
	NillInputErr       = "input cant be nill"
	NillProcessFuncErr = "process func cant be nill"
	ZeroLimitErr       = "limit must be greater than 0"
)

func Process(input <-chan interface{}, process processFunc, limit int) ([]interface{}, error) {
	if err := validate(input, process, limit); err != nil {
		return nil, err
	}

	mux := newMux(limit, true, process)
	go mux.fanOut(input)

	mux.waitAll()
	if err := mux.errors(); err != nil {
		return nil, err
	}
	return mux.fanIn(), nil
}
func ProcessSlice(input []interface{}, process processFunc, limit int) ([]interface{}, error) {
	if err := validate(input, process, limit); err != nil {
		return nil, err
	}
	mux := newMux(limit, true, process)
	go fanOutSlice(mux, input)
	mux.waitAll()
	if err := mux.errors(); err != nil {
		return nil, err
	}
	return mux.fanIn(), nil
}

func validate(input interface{}, process processFunc, limit int) error {
	if input == nil {
		return errors.New(NillInputErr)
	}
	if process == nil {
		return errors.New(NillProcessFuncErr)
	}
	if limit <= 0 {
		return errors.New(ZeroLimitErr)
	}
	return nil
}

func fanOutSlice(mux *mux, input []interface{}) {
	for i, v := range input {
		mux.getWorker().add(&item{value: v, index: i})
	}
	mux.closeAllInputChannels()
}

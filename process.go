package concurrent

import "errors"

type processFunc func(interface{}) (interface{}, error)

func Process(input <-chan interface{}, process processFunc, limit int) ([]interface{}, error) {
	if err := validate(input, process, limit); err != nil {
		return nil, err
	}

	mux := newMux(limit, true, process)
	prepareFanOutFromChannel(mux, input)
	if err := mux.fanOut(); err != nil {
		return nil, err
	}
	return mux.fanIn(), nil
}
func ProcessSlice(input []interface{}, process processFunc, limit int) ([]interface{}, error) {
	if err := validate(input, process, limit); err != nil {
		return nil, err
	}

	mux := newMux(limit, true, process)
	prepareFanOutFromSlice(mux, input)
	if err := mux.fanOut(); err != nil {
		return nil, err
	}
	return mux.fanIn(), nil
}

func validate(input interface{}, process processFunc, limit int) error {
	if input == nil {
		return errors.New("input cant be nill")
	}
	if process == nil {
		return errors.New("process func cant be nill")
	}
	if limit == 0 {
		return errors.New("limit must be greater than 0")
	}
	return nil
}

func prepareFanOutFromChannel(mux *mux, input <-chan interface{}) {
	i := 0
	for v := range input {
		mux.getInputChan() <- &item{value: v, index: i}
		i++
	}
	mux.closeAllInputChannels()
}
func prepareFanOutFromSlice(mux *mux, input []interface{}) {
	for i, v := range input {
		mux.getInputChan() <- &item{value: v, index: i}
	}
	mux.closeAllInputChannels()
}

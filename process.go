package concurrent

type processFunc func(interface{}) (interface{}, error)

func Process(input <-chan interface{}, process processFunc, limit int) ([]interface{}, error) {
	mux := newMux(limit, true, process)
	prepareFanOutChannels(mux, input)
	return mux.fanOut()
}

func prepareFanOutChannels(mux *mux, input <-chan interface{}) {
	i := 0
	for v := range input {
		mux.getInputChan() <- &item{value: v, index: i}
		i++
	}
	mux.closeAllInputChannels()
}

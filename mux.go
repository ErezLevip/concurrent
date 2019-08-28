package concurrent

import (
	"sync"
)

type mux struct {
	workers         [] *worker
	limit           int
	m               *sync.Mutex
	wg              *sync.WaitGroup
	continueOnError bool
	process         processFunc
}

type worker struct {
	in      chan *item
	out     chan *item
	errChan chan error
}

type item struct {
	index int
	value interface{}
}

func newMux(limit int, continueOnError bool, process processFunc) *mux {
	return &mux{
		workers:         [] *worker{},
		limit:           limit,
		m:               &sync.Mutex{},
		wg:              &sync.WaitGroup{},
		continueOnError: continueOnError,
		process:         process,
	}
}

func (m *mux) getInputChan() chan *item {
	for _, c := range m.workers {
		if m.limit > len(c.in) {
			return c.in
		}
	}
	w := newWorker(m.limit)
	m.m.Lock()
	defer m.m.Unlock()
	m.workers = append(m.workers, w)
	return w.in
}

func (m *mux) fanOut() ([]interface{}, error) {
	for _, w := range m.workers {
		w.startWorker(m.wg, m.limit, m.process, m.continueOnError)
	}
	m.wg.Wait()

	if err := m.checkErrors(); err != nil {
		return nil, err
	}
	return m.aggregateResults(), nil
}

func newWorker(limit int) *worker {
	return &worker{
		in:      make(chan *item, limit),
		errChan: make(chan error, 1),
	}
}

func (w *worker) startWorker(wg *sync.WaitGroup, limit int, process processFunc, continueOnError bool) {
	w.out = make(chan *item, limit)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(w.out)
		defer close(w.errChan)
		for raw := range w.in {
			v, err := process(raw.value)
			if err != nil {
				if continueOnError {
					continue
				}
				w.errChan <- err
				return
			}
			raw.value = v
			w.out <- raw
		}
	}()
}
func (m *mux) checkErrors() error {
	for _, w := range m.workers {
		for err := range w.errChan {
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *mux) closeAllInputChannels() {
	for _, w := range m.workers {
		close(w.in)
	}

}

func (m *mux) aggregateResults() []interface{} {
	results := make([]interface{}, m.countItems())
	i := 0
	for _, w := range m.workers {
		for v := range w.out {
			results[i] = v.value
			i++
		}
	}
	return results
}

func (m *mux) countItems() int {
	count := 0
	for _, w := range m.workers {
		count += len(w.out)
	}
	return count
}

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

type item struct {
	index int
	value interface{}
}

func newMux(limit int, continueOnError bool, process processFunc) *mux {
	m := &mux{
		workers:         []*worker{},
		limit:           limit,
		m:               &sync.Mutex{},
		wg:              &sync.WaitGroup{},
		continueOnError: continueOnError,
		process:         process,
	}

	//the mux will always start with 1 worker that will create and hold the first channel
	m.addWorker()
	return m
}
func (m *mux) addWorker() *worker {
	w := newWorker()
	m.m.Lock()
	defer m.m.Unlock()
	m.workers = append(m.workers, w)
	w.run(m.wg, m.limit, m.process, m.continueOnError)
	return w
}

func (m *mux) getWorker() *worker {
	for _, w := range m.workers {
		if m.limit > w.count {
			return w
		}
	}
	return m.addWorker()
}

func (m *mux) waitAll() {
	m.wg.Wait()
}

func (w *worker) run(wg *sync.WaitGroup, limit int, process processFunc, continueOnError bool) {
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
func (m *mux) errors() error {
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

func (m *mux) fanIn() []interface{} {
	results := make([]interface{}, m.countItems())
	i := 0
	for _, w := range m.workers {
		for v := range w.out {
			results[v.index] = v.value
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

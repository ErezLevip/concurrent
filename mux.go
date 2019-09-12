package concurrent

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type mux struct {
	workers         [] *worker
	limit           int
	m               *sync.Mutex
	wg              *sync.WaitGroup
	continueOnError bool
	process         processFunc
	kill            chan os.Signal
}

type item struct {
	index int
	value interface{}
}

func newMux(limit int, continueOnError bool, process processFunc) *mux {
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
	m := &mux{
		workers:         []*worker{},
		limit:           limit,
		m:               &sync.Mutex{},
		wg:              &sync.WaitGroup{},
		continueOnError: continueOnError,
		process:         process,
		kill:            kill,
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

func (m *mux) fanOut(input <-chan interface{}) {
	defer m.wg.Done()
	defer m.closeAllInputChannels()

	m.wg.Add(1)
	i := 0

	for {
		select {
		case v, canRead := <-input:
			if !canRead {
				return
			}
			m.getWorker().add(&item{value: v, index: i})
			i++
		case <-m.kill:
			return
		}
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

package concurrent

import (
	"github.com/google/uuid"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type worker struct {
	in      chan *item
	out     chan *item
	errChan chan error
	id      string
	count   int
	kill    chan os.Signal
}

func newWorker() *worker {
	id, _ := uuid.NewRandom()
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
	return &worker{
		in:      make(chan *item),
		errChan: make(chan error, 1),
		id:      id.String(),
		kill:    kill,
	}
}

func (w *worker) run(wg *sync.WaitGroup, limit int, process processFunc, continueOnError bool) {
	w.out = make(chan *item, limit)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(w.out)
		defer close(w.errChan)
		for {
			select {
			case raw, canRead := <-w.in:
				if !canRead {
					return
				}
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
			case <-w.kill:
				return
			}
		}
	}()
}

func (w *worker) add(item *item) {
	w.count += 1
	w.in <- item
}

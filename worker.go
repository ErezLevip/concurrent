package concurrent

import (
	"github.com/google/uuid"
)

type worker struct {
	in      chan *item
	out     chan *item
	errChan chan error
	id      string
	count   int
}

func newWorker() *worker {
	id, _ := uuid.NewRandom()
	return &worker{
		in:      make(chan *item),
		errChan: make(chan error, 1),
		id:      id.String(),
	}
}

func (w *worker) add(item *item) {
	w.count += 1
	w.in <- item
}

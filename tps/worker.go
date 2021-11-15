package tps

import (
	"log"
	"sync/atomic"
)

var (
	DefaultDoFunc = func(t Task) error {
		return nil
	}
)

type GeneralWorker interface {
	Do(task Task) error
	Run()
	Close()
}

type Worker struct {
	closing    uint32
	doTaskFunc func(Task) error
}

func NewWorker(doTask func(Task) error) Worker {
	var doTaskFunc func(Task) error
	if doTask != nil {
		doTaskFunc = doTask
	} else {
		doTaskFunc = DefaultDoFunc
	}

	return Worker{doTaskFunc: doTaskFunc}
}

func (w *Worker) Run(queue *Queue) {
	for {
		if atomic.LoadUint32(&w.closing) == 1 {
			break
		}

		task, isEmpty := queue.Shift()
		if isEmpty {
			continue
		}
		if err := w.doTaskFunc(task); err != nil {
			log.Fatal("err doTaskFunc", err)
		}
	}
}

func (w *Worker) Close() {
	atomic.StoreUint32(&w.closing, 1)
}

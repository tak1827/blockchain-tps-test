package tps

import (
	"log"
	"sync/atomic"
)

var (
	DefaultDoFunc = func(t Task, id int) error {
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
	doTaskFunc func(Task, int) error
}

func NewWorker(doTask func(Task, int) error) Worker {
	var doTaskFunc func(Task, int) error
	if doTask != nil {
		doTaskFunc = doTask
	} else {
		doTaskFunc = DefaultDoFunc
	}

	return Worker{
		doTaskFunc: doTaskFunc,
	}
}

func (w *Worker) Run(queue *Queue, id int) {
	for {
		if atomic.LoadUint32(&w.closing) == 1 {
			break
		}

		task, isEmpty := queue.Shift()
		if isEmpty {
			continue
		}
		if err := w.doTaskFunc(task, id); err != nil {
			log.Fatal("err doTaskFunc", err)
		}
	}
}

func (w *Worker) Close() {
	atomic.StoreUint32(&w.closing, 1)
}

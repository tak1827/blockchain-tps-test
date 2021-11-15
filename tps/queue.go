package tps

import (
	"log"
	"sync"
)

const (
	QueusSize = 1 << 24
)

type Queue struct {
	sync.Mutex
	Tasks []Task
}

func NewQueue(size int) Queue {
	if size == 0 {
		size = QueusSize
	}
	return Queue{
		Tasks: make([]Task, 0, size),
	}
}

func (q *Queue) Push(task Task) {
	q.Lock()
	defer q.Unlock()

	if len(q.Tasks)+1 > QueusSize {
		log.Fatal("queue overflow")
	}

	q.Tasks = append(q.Tasks, task)
}

func (q *Queue) Shift() (task Task, isEmpty bool) {
	q.Lock()
	defer q.Unlock()

	if len(q.Tasks) == 0 {
		isEmpty = true
		return
	}

	task, q.Tasks = q.Tasks[0], q.Tasks[1:]
	return
}

func (q *Queue) IsEmpty() bool {
	q.Lock()
	defer q.Unlock()

	return len(q.Tasks) == 0
}

func (q *Queue) CountTasks() int {
	q.Lock()
	defer q.Unlock()

	return len(q.Tasks)
}

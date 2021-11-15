package tps

import (
	"sync"
	"testing"
)

func TestPush(t *testing.T) {
	var (
		queue = Queue{}
		wg    = &sync.WaitGroup{}
		tasks = []Task{
			&BasicTask{},
			&BasicTask{},
			&BasicTask{},
			&BasicTask{},
			&BasicTask{},
			&BasicTask{},
			&BasicTask{},
			&BasicTask{},
			&BasicTask{},
			&BasicTask{},
		}
		pararelCount = 5
	)

	for i := 0; i < pararelCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < len(tasks); j++ {
				queue.Push(tasks[j])
			}
		}()
	}

	wg.Wait()

	if g, w := len(queue.Tasks), len(tasks)*pararelCount; g != w {
		t.Errorf("got: %d, want: %d", g, w)
	}
}

func TestShift(t *testing.T) {
	var (
		queue = Queue{
			Tasks: []Task{
				&BasicTask{},
				&BasicTask{},
				&BasicTask{},
				&BasicTask{},
				&BasicTask{},
				&BasicTask{},
				&BasicTask{},
				&BasicTask{},
				&BasicTask{},
				&BasicTask{},
			},
		}

		wg           = &sync.WaitGroup{}
		pararelCount = 5
	)

	for i := 0; i < pararelCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if _, isEmpty := queue.Shift(); isEmpty {
					break
				}
			}
		}()
	}

	wg.Wait()

	if g, w := len(queue.Tasks), 0; g != w {
		t.Errorf("got: %d, want: %d", g, w)
	}
}

package main

import (
	"math/big"
	"sync"
	"testing"

	"github.com/min-sys/tracking-contract/cli/eth"
)

func TestPush(t *testing.T) {
	addr, _ := eth.GenerateAddr()

	var (
		queue = Queue{}
		wg    = &sync.WaitGroup{}
		tasks = []Task{
			&ETHSendingTask{amount: big.NewInt(1), to: addr},
			&ETHSendingTask{amount: big.NewInt(2), to: addr},
			&ETHSendingTask{amount: big.NewInt(3), to: addr},
			&ETHSendingTask{amount: big.NewInt(4), to: addr},
			&ETHSendingTask{},
			&ETHSendingTask{},
			&ETHSendingTask{},
			&ETHSendingTask{},
			&ETHSendingTask{},
			&ETHSendingTask{},
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
	addr, _ := eth.GenerateAddr()

	var (
		queue = Queue{
			Tasks: []Task{
				&ETHSendingTask{amount: big.NewInt(1), to: addr},
				&ETHSendingTask{amount: big.NewInt(2), to: addr},
				&ETHSendingTask{amount: big.NewInt(3), to: addr},
				&ETHSendingTask{amount: big.NewInt(4), to: addr},
				&ETHSendingTask{},
				&ETHSendingTask{},
				&ETHSendingTask{},
				&ETHSendingTask{},
				&ETHSendingTask{},
				&ETHSendingTask{},
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

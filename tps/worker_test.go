package tps

import (
	"testing"
	"time"
)

func Test_Run_Close(t *testing.T) {
	var (
		w = NewWorker(nil, 1)
		q = NewQueue(1)
	)

	go w.Run(&q)

	time.Sleep(100 * time.Millisecond)

	w.Close()

	if g, w := w.closing, uint32(1); g != w {
		t.Errorf("got: %d, want: %d", g, w)
	}
}

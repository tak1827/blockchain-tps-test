package main

import (
	"testing"
	"time"
)

func Test_Run_Close(t *testing.T) {
	var (
		w = NewWorker(nil)
	)

	go w.Run(nil)

	time.Sleep(100 * time.Millisecond)

	w.Close()

	if g, w := w.closing, uint32(1); g != w {
		t.Errorf("got: %d, want: %d", g, w)
	}
}

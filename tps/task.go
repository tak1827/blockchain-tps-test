package tps

import (
	"fmt"
)

type TaskType int

const (
	ETHSending TaskType = iota

	TaskRetryLimit = 100
)

type Task interface {
	Type() TaskType
	TryCount() int
	IncrementTryCount() error
}

type BasicTask struct {
	tryCount int
}

func (t *BasicTask) Type() TaskType {
	return ETHSending
}

func (t *BasicTask) TryCount() int {
	return t.tryCount
}

func (t *BasicTask) IncrementTryCount() error {
	t.tryCount += 1
	if t.tryCount >= TaskRetryLimit {
		return fmt.Errorf("err task retry limit, tryCount: %d", t.tryCount)
	}
	return nil
}

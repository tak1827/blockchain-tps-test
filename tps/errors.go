package tps

import (
	"errors"
)

var (
	ErrStopTask    = errors.New("stop doing task")
	ErrTaskRetry   = errors.New("task retried")
	ErrTxFailed    = errors.New("transaction is failed")
	ErrWrongNonce  = errors.New("transaction nonce is not correct")
	ErrNotNewBlock = errors.New("new block have not yet mind")
)

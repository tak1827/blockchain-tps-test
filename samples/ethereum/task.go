package main

import (
	"context"
	"fmt"
	"strings"

	// "github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tak1827/blockchain-tps-test/tps"
)

const (
	Eth tps.TaskType = iota

	TaskRetryLimit = 10
)

type EthTask struct {
	to       common.Address
	amount   int64
	input    []byte
	gasLimit uint64

	tryCount int
}

func (t *EthTask) Type() tps.TaskType {
	return Eth
}

func (t *EthTask) TryCount() int {
	return t.tryCount
}

func (t *EthTask) IncrementTryCount() error {
	t.tryCount += 1
	if t.tryCount >= TaskRetryLimit {
		return fmt.Errorf("err task retry limit, tryCount: %d", t.tryCount)
	}
	return nil
}

func (t *EthTask) Do(ctx context.Context, client *EthClient, priv string, nonce uint64, queue *tps.Queue, logger tps.Logger) error {
	res, rootErr := client.SendTx(ctx, priv, nonce, t.to, t.amount, t.input, t.gasLimit)
	if rootErr != nil {
		if strings.Contains(rootErr.Error(), "invalid sequence") {
			logger.Warn(fmt.Sprintf("nonce error, %s", rootErr.Error()))
			return tps.ErrWrongNonce
		}

		logger.Warn(fmt.Sprintf("faild sending, err: %s", rootErr.Error()))
		if err := t.IncrementTryCount(); err != nil {
			return errors.Wrap(rootErr, err.Error())
		}
		queue.Push(t)
		return nil
	}

	logger.Info(fmt.Sprintf("succeed sending tx, hash: %s, nonce: %d", res.Hash(), nonce))
	return nil
}

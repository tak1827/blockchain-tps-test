package main

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	// "github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/min-sys/tracking-contract/cli/eth"
	"github.com/pkg/errors"
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

type ETHSendingTask struct {
	to     common.Address
	amount *big.Int

	tryCount int
}

func (t *ETHSendingTask) Type() TaskType {
	return ETHSending
}

func (t *ETHSendingTask) TryCount() int {
	return t.tryCount
}

func (t *ETHSendingTask) IncrementTryCount() error {
	t.tryCount += 1
	if t.tryCount >= TaskRetryLimit {
		return fmt.Errorf("err task retry limit, tryCount: %d", t.tryCount)
	}
	return nil
}

func (t *ETHSendingTask) Do(client *eth.Client, priv *ecdsa.PrivateKey, nonce uint64, queue *Queue, logger Logger) error {
	var (
		ctx = eth.CtxWithPriv(nil, priv)
	)
	hash, rootErr := client.SendETH(ctx, nonce, t.to, t.amount)
	if rootErr != nil {
		if strings.Contains(rootErr.Error(), "the tx doesn't have the correct nonce") {
			logger.Warn(fmt.Sprintf("nonce error, %s", rootErr.Error()))
			return ErrWrongNonce
		}

		logger.Warn(fmt.Sprintf("faild sending eth, err: %s", rootErr.Error()))
		if err := t.IncrementTryCount(); err != nil {
			return errors.Wrap(rootErr, err.Error())
		}
		queue.Push(t)
	}

	logger.Info(fmt.Sprintf("succeed sending tx, hash: %s, nonce: %d", hash.String(), nonce))
	return nil
}

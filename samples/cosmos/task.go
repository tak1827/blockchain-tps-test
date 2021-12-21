package main

import (
	"context"
	"fmt"
	"strings"

	// "github.com/davecgh/go-spew/spew"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/tak1827/blockchain-tps-test/tps"
)

const (
	Cosm tps.TaskType = iota

	TaskRetryLimit = 10
)

type CosmTask struct {
	to     sdk.AccAddress
	amount int64

	tryCount int
}

func (t *CosmTask) Type() tps.TaskType {
	return Cosm
}

func (t *CosmTask) TryCount() int {
	return t.tryCount
}

func (t *CosmTask) IncrementTryCount() error {
	t.tryCount += 1
	if t.tryCount >= TaskRetryLimit {
		return fmt.Errorf("err task retry limit, tryCount: %d", t.tryCount)
	}
	return nil
}

func (t *CosmTask) Do(ctx context.Context, client *CosmosClient, priv string, nonce uint64, queue *tps.Queue, logger tps.Logger) error {
	res, rootErr := client.SendTx(ctx, priv, nonce, t.to, t.amount)
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

	logger.Info(fmt.Sprintf("succeed sending tx, hash: %s, nonce: %d", res.Hash, nonce))
	return nil
}

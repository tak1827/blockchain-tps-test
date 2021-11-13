package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	// "github.com/davecgh/go-spew/spew"
	"github.com/min-sys/tracking-contract/cli/eth"
	"github.com/pkg/errors"
)

func StartTPSMeasuring(client *eth.Client, closing, idlingDuration *uint32, logger Logger) error {
	var (
		canMesure bool
		startedAd time.Time
		total     uint
		count     uint
		lastBlock uint64
		err       error
	)

	for {
		if atomic.LoadUint32(closing) == 1 {
			break
		}

		if count, lastBlock, err = countTx(client, lastBlock); err != nil {
			if errors.Is(err, ErrNotNewBlock) {
				// sleep a bit
				time.Sleep(1 * time.Second)
				continue
			}
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Warn(fmt.Sprintf("timeout of countTx, setting: %v", client.Setting().Timeout))
				continue
			}
			//TODO: handle timeout error
			return errors.Wrap(err, "err CountTx")
		}

		if !canMesure {
			if count > 0 {
				canMesure = true
				startedAd = time.Now()
			}
			continue
		}

		pendingTx, err := client.TxpoolPendingTxCount(nil)
		if err != nil {
			return errors.Wrap(err, "err TxpoolPendingTxCount")
		}

		NextIdlingDuration(idlingDuration, uint32(count), uint32(pendingTx))

		total += count
		elapsed := time.Now().Sub(startedAd).Seconds()
		fmt.Print("------------------------------------------------------------------------------------\n")
		fmt.Printf("⛓  %d th Block Mind! txs(%d), total txs(%d), TPS(%.2f), pendig txs(%d)\n", lastBlock, count, total, float64(total)/elapsed, pendingTx)
	}

	return nil
}

func countTx(client *eth.Client, lastBlock uint64) (uint, uint64, error) {
	header, err := client.BlockHeader(nil)
	if err != nil {
		return 0, lastBlock, errors.Wrap(err, "err BlockHash")
	}
	if header.Number.Uint64() <= lastBlock {
		return 0, lastBlock, ErrNotNewBlock
	}

	count, err := client.TxCount(nil, header.Hash())
	if err != nil {
		return 0, lastBlock, errors.Wrap(err, "err TxCount")
	}

	return count, header.Number.Uint64(), nil
}

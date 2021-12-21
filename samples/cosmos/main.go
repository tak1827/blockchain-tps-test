package main

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// "github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/tak1827/blockchain-tps-test/tps"
)

var (
	Endpoint = "tcp://localhost:26657"

	PrivKey  = "3b2eb70cf00a779c4bfed132e0fb3f7f982013132d86e344b0e97c7507d0d7a4"
	PrivKey2 = "cb7e1d0611d5c66461822afcbed4d677b19f9188541c62c534338679161e9aa9"
	PrivKey3 = "098cc9ba1d5109b4b81fc06859dc99950617e9d50127ca940e8298bd9fb3c6eb"
	// PrivKey4 = "098cc9ba1d5109b4b81fc06859dc99950617e9d50127ca940e8298bd9fb3c6eb"

	Timeout        = 15 * time.Second
	MaxConcurrency = runtime.NumCPU() - 2
)

func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

func main() {
	var (
		mesuringDuration = 60 * time.Second
		queueSize        = 100
		concurrency      = 2
		queue            = tps.NewQueue(queueSize)
		closing          uint32
		idlingDuration   uint32
		logLevel         = tps.WARN_LEVEL // INFO_LEVEL, WARN_LEVEL, FATAL_LEVEL
		logger           = tps.NewLogger(logLevel)
		privs            = []string{
			PrivKey,
			PrivKey2,
			PrivKey3,
			// PrivKey4,
		}
		testAddrs = createRandomAccounts(100)
	)

	go func() {
		defer atomic.AddUint32(&closing, 1)
		time.Sleep(mesuringDuration)
	}()

	client, err := NewClient(Endpoint)
	if err != nil {
		logger.Fatal("err NewClient: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	addrs := make([]string, len(privs))
	for i := range privs {
		addr, err := AccAddressFromPrivString(privs[i])
		if err != nil {
			logger.Fatal("err AccAddressFromPrivString: ", err)
		}
		addrs[i] = addr

		acc, err := client.Account(ctx, addr)
		if err != nil {
			logger.Fatal("err Account: ", err)
		}

		accNums[addr] = acc.GetAccountNumber()
	}

	wallet, err := tps.NewWallet(ctx, &client, privs, addrs)
	if err != nil {
		logger.Fatal("err NewWallet: ", err)
	}

	taskDo := func(t tps.Task, id int) error {
		task, ok := t.(*CosmTask)
		if !ok {
			return errors.New("unexpected task type")
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		var (
			priv         = wallet.Priv(id)
			currentNonce = wallet.IncrementNonce(priv)
		)
		if err = task.Do(ctx, &client, priv, currentNonce, &queue, logger); err != nil {
			if errors.Is(err, tps.ErrWrongNonce) {
				wallet.RecetNonce(priv, currentNonce)
				task.tryCount = 0
				queue.Push(task)
				return nil
			}
			return errors.Wrap(err, "err Do")
		}

		// time.Sleep(ToDuration(&idlingDuration))

		return nil
	}

	worker := tps.NewWorker(taskDo)

	// performance likely not improved, whene exceed available cpu core
	if concurrency > MaxConcurrency {
		logger.Warn(fmt.Sprintf("concurrency setting is over logical max(%d)", MaxConcurrency))
	}
	for i := 0; i < concurrency; i++ {
		go worker.Run(&queue, i)
	}

	go func() {
		count := 0
		for {
			if atomic.LoadUint32(&closing) == 1 {
				break
			}

			if queue.CountTasks() > queueSize {
				continue
			}

			queue.Push(&CosmTask{
				to:     testAddrs[count%len(testAddrs)],
				amount: 1,
			})
			count++
		}
	}()

	if err = tps.StartTPSMeasuring(context.Background(), &client, &closing, &idlingDuration, logger); err != nil {
		logger.Fatal("err StartTPSMeasuring: ", err)
	}
}

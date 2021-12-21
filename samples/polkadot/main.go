package main

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	// "github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/tak1827/blockchain-tps-test/tps"
)

var (
	// Endpoint = "ws://127.0.0.1:9944"
	Endpoint = "ws://127.0.0.1:9945" // testnet

	Seed0 = "//Alice"
	Seed1 = "//Charlie"
	Seed2 = "//Dave"
	Seed3 = "//Eve"

	Timeout        = 15 * time.Second
	MaxConcurrency = runtime.NumCPU() - 2
)

func createRandomAccounts(accNum int) []types.MultiAddress {
	testAddrs := make([]types.MultiAddress, accNum)
	for i := 0; i < accNum; i++ {
		seed := fmt.Sprintf("//test_account_%d", i)
		key, err := signature.KeyringPairFromSecret(seed, 42)
		if err != nil {
			panic(err)
		}
		testAddrs[i] = types.NewMultiAddressFromAccountID(key.PublicKey)
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
			Seed0,
			Seed1,
			Seed2,
			Seed3,
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
		key, err := signature.KeyringPairFromSecret(privs[i], 42)
		if err != nil {
			panic(err)
		}
		addrs[i] = types.HexEncodeToString(key.PublicKey)
	}

	wallet, err := tps.NewWallet(ctx, &client, privs, addrs)
	if err != nil {
		logger.Fatal("err NewWallet: ", err)
	}

	taskDo := func(t tps.Task, id int) error {
		task, ok := t.(*PolkaTask)
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

			queue.Push(&PolkaTask{
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

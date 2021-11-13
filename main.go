package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/min-sys/tracking-contract/cli/eth"
	"github.com/pkg/errors"
	// "github.com/davecgh/go-spew/spew"
)

var (
	// geth
	Endpoint = "http://127.0.0.1:8545"
	PrivKey  = "58900163de10a0ffe2d4c3224faf8da4d45727cb47114eb43dd754b60ab70cb7"
	// 0xE3b0DE0E4CA5D3CB29A9341534226C4D31C9838f
	PrivKey2 = "d1c71e71b06e248c8dbe94d49ef6d6b0d64f5d71b1e33a0f39e14dadb070304a"
	// 0x26fa9f1a6568b42e29b1787c403B3628dFC0C6FE
	PrivKey3 = "8179ce3d00ac1d1d1d38e4f038de00ccd0e0375517164ac5448e3acc847acb34"

	MaxConcurrency = runtime.NumCPU() - 2
)

func main() {
	var (
		mesuringDuration = 30 * time.Second
		queueSize        = 100
		concurrency      = 2
		queue            = NewQueue(queueSize)
		closing          uint32
		idlingDuration   uint32
		logLevel         = WARN_LEVEL // INFO_LEVEL, WARN_LEVEL, FATAL_LEVEL
		logger           = NewLogger(logLevel)
		keys       = []string{
			PrivKey,
			PrivKey2,
			PrivKey3,
		}
	)

	go func() {
		defer atomic.AddUint32(&closing, 1)
		time.Sleep(mesuringDuration)
	}()

	client, err := eth.NewClient(nil, eth.WithEndpoint(Endpoint), eth.WithPrivKey(PrivKey))
	if err != nil {
		logger.Fatal("err NewClient: ", err)
	}

	wallet, err := NewWallet(&client, keys)
	if err != nil {
		logger.Fatal("err NewWallet: ", err)
	}

	taskDo := func(t Task) error {
		task, ok := t.(*ETHSendingTask)
		if !ok {
			return errors.New("unexpected task type")
		}

		var (
			priv         = wallet.RotatePriv()
			currentNonce = wallet.IncrementNonce(priv)
		)
		if err = task.Do(&client, priv, currentNonce, &queue, logger); err != nil {

			if errors.Is(err, ErrWrongNonce) {
				wallet.RecetNonce(priv, currentNonce)
				task.tryCount = 0
				queue.Push(task)
				return nil
			}

			return errors.Wrap(err, "err Do")
		}

		time.Sleep(ToDuration(&idlingDuration))

		return nil
	}

	worker := NewWorker(taskDo)

	// performance likely not improved, whene exceed available cpu core
	if concurrency > MaxConcurrency {
		logger.Warn(fmt.Sprintf("concurrency setting is over logical max(%d)", MaxConcurrency))
	}
	for i := 0; i < concurrency; i++ {
		go worker.Run(&queue)
	}

	go func() {
		amount := eth.ToWei(1.0, 1)
		to, err := eth.GenerateAddr()
		if err != nil {
			logger.Fatal("err GenerateAddr: ", err)
		}

		for {
			if atomic.LoadUint32(&closing) == 1 {
				break
			}

			if queue.CountTasks() >= queueSize {
				continue
			}

			queue.Push(&ETHSendingTask{
				to:     to,
				amount: amount,
			})
		}
	}()

	if err = StartTPSMeasuring(&client, &closing, &idlingDuration, logger); err != nil {
		logger.Fatal("err StartTPSMeasuring: ", err)
	}
}

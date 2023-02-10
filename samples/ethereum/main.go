package main

import (
	"context"
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tak1827/eth-extended-client/contract"

	// "github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/tak1827/blockchain-tps-test/tps"
)

var (
	Endpoint = "http://localhost:8545"

	PrivKey = "3b2eb70cf00a779c4bfed132e0fb3f7f982013132d86e344b0e97c7507d0d7a4"
	// PrivKey2 = "cb7e1d0611d5c66461822afcbed4d677b19f9188541c62c534338679161e9aa9"
	// PrivKey3 = "098cc9ba1d5109b4b81fc06859dc99950617e9d50127ca940e8298bd9fb3c6eb"
	// PrivKey4 = "098cc9ba1d5109b4b81fc06859dc99950617e9d50127ca940e8298bd9fb3c6eb"

	Timeout        = 15 * time.Second
	MaxConcurrency = runtime.NumCPU() - 2
)

func createRandomAccounts(accNum int) []common.Address {
	testAddrs := make([]common.Address, accNum)
	for i := 0; i < accNum; i++ {
		priv, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		testAddrs[i] = crypto.PubkeyToAddress(priv.PublicKey)
	}

	return testAddrs
}

func main() {
	var (
		ctx              = context.Background()
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
			// PrivKey2,
			// PrivKey3,
			// PrivKey4,
		}
		testAddrs = createRandomAccounts(100)
	)

	go func() {
		defer atomic.AddUint32(&closing, 1)
		time.Sleep(mesuringDuration)
	}()

	client, err := NewClient(ctx, Endpoint)
	if err != nil {
		logger.Fatal("err NewClient: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	addrs := make([]string, len(privs))
	for i := range privs {
		privKey, err := crypto.HexToECDSA(privs[i])
		if err != nil {
			logger.Fatal(err)
		}
		addrs[i] = crypto.PubkeyToAddress(privKey.PublicKey).Hex()
	}

	wallet, err := tps.NewWallet(ctx, &client, privs, addrs)
	if err != nil {
		logger.Fatal("err NewWallet: ", err)
	}

	taskDo := func(t tps.Task, id int) error {
		task, ok := t.(*EthTask)
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

	parsed, _ := abi.JSON(strings.NewReader(contract.IERC721ABI))

	go func() {
		count := 0
		for {
			if atomic.LoadUint32(&closing) == 1 {
				break
			}

			if queue.CountTasks() > queueSize {
				continue
			}

			// calldate of mint nft
			var (
				to       = testAddrs[count%len(testAddrs)]
				tokenId  = big.NewInt(int64(count))
				input, _ = parsed.Pack("safeMint", []interface{}{tokenId, to, ""}...)
				gasLimit = uint64(100000) // gaslimit of mint nft
			)
			queue.Push(&EthTask{
				to:       to,
				amount:   0,
				input:    input,
				gasLimit: gasLimit,
			})
			count++
		}
	}()

	if err = tps.StartTPSMeasuring(context.Background(), &client, &closing, &idlingDuration, logger); err != nil {
		logger.Fatal("err StartTPSMeasuring: ", err)
	}
}

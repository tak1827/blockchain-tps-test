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

	PrivKey  = "10CD36EB1C4D85EA12C4CEB457EE6B87CA4A653E78B0275B93DDE95ECE21AAC0"
	PrivKey2 = "AF3474C24F7F4BCC2E7F01ABCCB245852E0564E2FC4E99133BCB5F64882CF8EB"
	PrivKey3 = "8F0F19CAC1178D1990DD6840FACB3F9CD6E09A9E08EFFF106AC0069EF180092D"
	PrivKey4 = "445C97786D9478F8970876451B52226D9FEE34E750711297B5787192676B39CD"

	ContractAddress = common.HexToAddress("0x89FB319f064cf99a7c5bc7b69d7064ADCFb990e9")

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
		concurrency      = 1
		queue            = tps.NewQueue(queueSize)
		closing          uint32
		idlingDuration   uint32
		logLevel         = tps.WARN_LEVEL // INFO_LEVEL, WARN_LEVEL, FATAL_LEVEL
		logger           = tps.NewLogger(logLevel)
		privs            = []string{
			PrivKey,
			PrivKey2,
			PrivKey3,
			PrivKey4,
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
				account  = testAddrs[count%len(testAddrs)]
				tokenId  = big.NewInt(int64(count))
				input, _ = parsed.Pack("safeMint", []interface{}{tokenId, account, ""}...)
				gasLimit = uint64(100000) // gaslimit of mint nft
			)
			queue.Push(&EthTask{
				to:       ContractAddress,
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

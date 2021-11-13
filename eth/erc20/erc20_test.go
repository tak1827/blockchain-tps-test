package erc20

import (
	"context"
	"math"
	"math/big"
	"sync"
	"testing"

	// "github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/crypto"
	"github.com/min-sys/tracking-contract/cli/constant"
	"github.com/min-sys/tracking-contract/cli/eth"
)

func TestRecord(t *testing.T) {
	var (
		ctx  = context.Background()
		size = 1
	)

	client, err := eth.NewClient(nil, eth.WithEndpoint(constant.TestEndpoint), eth.WithPrivKey(constant.TestPrivKey))
	if err != nil {
		t.Fatal("err create NewClient", err)
	}

	token, err := NewMindenToken(common.HexToAddress(constant.TestMindenERC20Addr), &client)
	if err != nil {
		t.Fatal("err NewMindenToken", err)
	}

	demanders, err := eth.GenerateAddrs(size)
	if err != nil {
		t.Fatal("err GenerateAddrs", err)
	}

	amounts := make([]*big.Int, size)
	eth.GenerateRandInts(size, math.MaxInt64, func(num, index int) {
		amounts[index] = big.NewInt(int64(num))
	})

	datatimes := make([]uint32, size)
	eth.GenerateRandInts(size, math.MaxUint32, func(num, index int) {
		datatimes[index] = uint32(num)
	})

	indexs := make([]uint8, size)
	eth.GenerateRandInts(size, math.MaxUint8, func(num, index int) {
		indexs[index] = uint8(num)
	})

	tx, err := token.Record(ctx, nil, demanders[0], amounts[0], datatimes[0], indexs[0])
	if err != nil {
		t.Error("err Record", err)
	}

	if err = token.Client.ConfirmTx(ctx, tx.Hash()); err != nil {
		t.Error("err confirmTx", err)
	}
}

func TestRecordBatch(t *testing.T) {
	var (
		ctx  = context.Background()
		size = 10
	)

	client, err := eth.NewClient(nil, eth.WithEndpoint(constant.TestEndpoint), eth.WithPrivKey(constant.TestPrivKey))
	if err != nil {
		t.Fatal("err create NewClient", err)
	}

	token, err := NewMindenToken(common.HexToAddress(constant.TestMindenERC20Addr), &client)
	if err != nil {
		t.Fatal("err NewMindenToken", err)
	}

	demanders, err := eth.GenerateAddrs(10)
	if err != nil {
		t.Fatal("err GenerateAddrs", err)
	}

	amounts := make([]*big.Int, size)
	eth.GenerateRandInts(size, math.MaxInt64, func(num, index int) {
		amounts[index] = big.NewInt(int64(num))
	})

	datatimes := make([]uint32, size)
	eth.GenerateRandInts(size, math.MaxUint32, func(num, index int) {
		datatimes[index] = uint32(num)
	})

	indexs := make([]uint8, size)
	eth.GenerateRandInts(size, math.MaxUint8, func(num, index int) {
		indexs[index] = uint8(num)
	})

	tx, err := token.RecordBatch(ctx, nil, demanders, amounts, datatimes, indexs)
	if err != nil {
		t.Error("err RecordBatch", err)
	}

	if err = token.Client.ConfirmTx(ctx, tx.Hash()); err != nil {
		t.Error("err confirmTx", err)
	}
}

func Test_ParallelRecord(t *testing.T) {
	parallelNumber := uint64(10)

	client, err := eth.NewClient(nil, eth.WithEndpoint(constant.TestEndpoint), eth.WithPrivKey(constant.TestPrivKey))
	if err != nil {
		t.Fatal("err create NewClient", err)
	}

	nonce, err := client.Nonce(nil)
	if err != nil {
		t.Fatal("err Nonce", err)
	}

	var (
		ctx  = context.Background()
		size = 10
	)

	token, err := NewMindenToken(common.HexToAddress(constant.TestMindenERC20Addr), &client)
	if err != nil {
		t.Fatal("err NewMindenToken", err)
	}

	demanders, err := eth.GenerateAddrs(10)
	if err != nil {
		t.Fatal("err GenerateAddrs", err)
	}

	amounts := make([]*big.Int, size)
	eth.GenerateRandInts(size, math.MaxInt64, func(num, index int) {
		amounts[index] = big.NewInt(int64(num))
	})

	datatimes := make([]uint32, size)
	eth.GenerateRandInts(size, math.MaxUint32, func(num, index int) {
		datatimes[index] = uint32(num)
	})

	indexs := make([]uint8, size)
	eth.GenerateRandInts(size, math.MaxUint8, func(num, index int) {
		indexs[index] = uint8(num)
	})

	wg := &sync.WaitGroup{}
	for i := uint64(0); i < parallelNumber; i++ {
		wg.Add(1)

		tx, err := token.RecordBatch(ctx, big.NewInt(int64(nonce+i)), demanders, amounts, datatimes, indexs)
		if err != nil {
			t.Error("err RecordBatch", err)
		}

		go func(hash common.Hash) {
			defer wg.Done()
			if err = token.Client.ConfirmTx(ctx, hash); err != nil {
				t.Error("err confirmTx", err)
			}
		}(tx.Hash())
	}
	wg.Wait()
}

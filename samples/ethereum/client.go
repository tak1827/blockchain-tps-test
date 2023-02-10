package main

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/tak1827/blockchain-tps-test/tps"
)

const (
	DefaultGasPrice = int64(0)   // 1 gwai
	DefaultTimeout  = int64(100) // 100 msec
)

var (
	_ tps.Client = (*EthClient)(nil)

	accNums         = make(map[string]uint64)
	timeoutDuration time.Duration
)

type EthClient struct {
	ethclient *ethclient.Client

	GasPrice   *big.Int
	GasTip     *big.Int
	GasBaseFee *big.Int
	chainID    *big.Int
	timeout    int64
}

func NewClient(ctx context.Context, endpoint string) (c EthClient, err error) {
	rpcclient, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		return
	}

	c.ethclient = ethclient.NewClient(rpcclient)
	c.GasPrice = big.NewInt(int64(DefaultGasPrice))
	c.timeout = DefaultTimeout

	supports, err := c.isSupportEIP1559(ctx)
	if err != nil {
		return
	}
	if supports {
		c.GasTip, err = c.ethclient.SuggestGasTipCap(ctx)
		if err != nil {
			return
		}
	}

	header, err := c.ethclient.HeaderByNumber(ctx, nil)
	if err != nil {
		return
	}
	c.GasBaseFee = header.BaseFee

	if c.chainID, err = c.ethclient.ChainID(ctx); err != nil {
		return
	}

	timeoutDuration = time.Duration(time.Duration(c.timeout) * time.Millisecond)

	return
}

func (c EthClient) Nonce(ctx context.Context, address string) (nonce uint64, err error) {
	nonce, err = c.ethclient.NonceAt(ctx, common.HexToAddress(address), nil)
	return
}

func (c EthClient) LatestBlockHeight(ctx context.Context) (uint64, error) {
	return c.ethclient.BlockNumber(ctx)
}

func (c EthClient) CountTx(ctx context.Context, height uint64) (int, error) {
	block, err := c.ethclient.BlockByNumber(ctx, big.NewInt(int64(height)))
	if err != nil {
		return 0, err
	}
	return len(block.Transactions()), nil
}

func (c EthClient) CountPendingTx(ctx context.Context) (int, error) {
	count, err := c.ethclient.PendingTransactionCount(ctx)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (c *EthClient) SendTx(ctx context.Context, priv string, nonce uint64, to common.Address, amount int64, input []byte, gasLimit uint64) (*types.Transaction, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	signedTx, err := c.sinedTx(timeoutCtx, priv, nonce, &to, big.NewInt(amount), input, gasLimit)
	if err != nil {
		return signedTx, errors.Wrap(err, "failed to sign tx")
	}

	if err := c.ethclient.SendTransaction(ctx, signedTx); err != nil {
		return signedTx, errors.Wrap(err, "err SendTransaction")
	}

	return signedTx, nil
}

func (c *EthClient) isSupportEIP1559(ctx context.Context) (bool, error) {
	if _, err := c.ethclient.SuggestGasTipCap(ctx); err != nil {
		if strings.Contains(err.Error(), "eth_maxPriorityFeePerGas does not exist") {
			return false, nil
		} else {
			return false, errors.Wrap(err, "failed to get suggestiion of gas tip cap")
		}
	}
	return true, nil
}

func (c *EthClient) sinedTx(ctx context.Context, priv string, nonce uint64, to *common.Address, amount *big.Int, input []byte, gasLimit uint64) (*types.Transaction, error) {
	privKey, err := crypto.HexToECDSA(priv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get priv")
	}

	var (
		isDynamic bool
		txdata    types.TxData
		signer    = types.NewLondonSigner(c.chainID)
	)
	if c.GasTip != nil {
		isDynamic = true
	}
	if isDynamic {
		gasFee := c.computeGasFee()

		if gasLimit == 0 {
			auth := bind.NewKeyedTransactor(privKey)
			if gasLimit, err = c.estimateGasLimit(ctx, auth.From, to, amount, input, gasFee); err != nil {
				return nil, errors.Wrap(err, "failed to estimate gas")
			}
		}

		txdata = &types.DynamicFeeTx{
			ChainID:    c.chainID,
			Nonce:      nonce,
			GasTipCap:  c.GasTip,
			GasFeeCap:  gasFee,
			Gas:        gasLimit,
			To:         to,
			Value:      amount,
			Data:       input,
			AccessList: nil,
		}
	} else {
		if gasLimit == 0 {
			auth := bind.NewKeyedTransactor(privKey)
			if gasLimit, err = c.estimateGasLimit(ctx, auth.From, to, amount, input, nil); err != nil {
				return nil, errors.Wrap(err, "failed to estimate gas")
			}
		}
		txdata = &types.LegacyTx{
			Nonce:    nonce,
			GasPrice: c.GasPrice,
			Gas:      gasLimit,
			To:       to,
			Value:    amount,
			Data:     input,
		}
	}

	tx, err := types.SignNewTx(privKey, signer, txdata)
	if err != nil {
		return nil, errors.Wrap(err, "at types.SignNewTx")
	}

	return tx, nil
}

func (c *EthClient) computeGasFee() *big.Int {
	// ref: https://github.com/ethereum/go-ethereum/blob/v1.10.17/accounts/abi/bind/base.go#L252
	return new(big.Int).Add(c.GasTip, new(big.Int).Mul(c.GasBaseFee, big.NewInt(2)))
}

func (c *EthClient) estimateGasLimit(ctx context.Context, from common.Address, to *common.Address, value *big.Int, input []byte, gasFee *big.Int) (uint64, error) {
	msg := ethereum.CallMsg{
		From:      from,
		To:        to,
		GasPrice:  c.GasPrice,
		GasTipCap: c.GasTip,
		GasFeeCap: gasFee,
		Value:     value,
		Data:      input,
	}
	return c.ethclient.EstimateGas(ctx, msg)
}

package erc20

import (
	"context"
	"math/big"

	// "github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/min-sys/tracking-contract/cli/eth"
	"github.com/pkg/errors"
)

type MindenToken struct {
	ContractAddress common.Address
	Client          *eth.Client
	Contract        *MindenERC20
}

func NewMindenToken(address common.Address, client *eth.Client) (MindenToken, error) {
	var (
		token = MindenToken{ContractAddress: address, Client: client}
		err   error
	)

	if token.Contract, err = NewMindenERC20(address, token.Client.Ethclient()); err != nil {
		return token, errors.Wrap(err, "err NewMindenERC20")
	}

	return token, nil
}

func (m *MindenToken) Record(ctx context.Context, nonce *big.Int, demander common.Address, amount *big.Int, datetime uint32, index uint8) (*types.Transaction, error) {
	timeoutCtx, cancel := m.Client.Setting().TimeoutContext(ctx)
	defer cancel()

	var (
		setting = m.Client.Setting()
		auth    = bind.NewKeyedTransactor(setting.PrivKey)
		opts    = &bind.TransactOpts{
			From:     auth.From,
			Nonce:    nonce, // nil = use pending state
			Signer:   auth.Signer,
			GasPrice: setting.GasPrice,
			GasLimit: 0, // estimate
			Context:  timeoutCtx,
		}
	)
	return m.Contract.Record(opts, demander, amount, datetime, index)
}

func (m *MindenToken) RecordBatch(ctx context.Context, nonce *big.Int, demanders []common.Address, amounts []*big.Int, datetimes []uint32, indexs []uint8) (*types.Transaction, error) {
	timeoutCtx, cancel := m.Client.Setting().TimeoutContext(ctx)
	defer cancel()

	var (
		setting = m.Client.Setting()
		auth    = bind.NewKeyedTransactor(setting.PrivKey)
		opts    = &bind.TransactOpts{
			From:     auth.From,
			Nonce:    nonce, // nil = use pending state
			Signer:   auth.Signer,
			GasPrice: setting.GasPrice,
			GasLimit: 0, // estimate
			Context:  timeoutCtx,
		}
	)

	return m.Contract.RecordBatch(opts, demanders, amounts, datetimes, indexs)
}

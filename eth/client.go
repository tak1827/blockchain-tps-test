package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	// "github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

const (
	EthSendingGasLimit = 21000
)

type Client struct {
	ethclient *ethclient.Client
	c         *rpc.Client

	setting Setting
}

func newClient(ctx context.Context, s Setting) (Client, error) {
	timeoutCtx, cancel := s.TimeoutContext(ctx)
	defer cancel()

	c, err := rpc.DialContext(timeoutCtx, s.Endpoint)
	if err != nil {
		return Client{}, fmt.Errorf("failed to conecting client, endpoint: %s", s.Endpoint)
	}

	return Client{
		ethclient: ethclient.NewClient(c),
		c:         c,
		setting:   s,
	}, nil
}

func (c *Client) BlockNumer(ctx context.Context) (uint64, error) {
	timeoutCtx, cancel := c.setting.TimeoutContext(ctx)
	defer cancel()

	header, err := c.ethclient.HeaderByNumber(timeoutCtx, nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

func (c *Client) BlockHash(ctx context.Context) (common.Hash, error) {
	var ZeroHash common.Hash

	timeoutCtx, cancel := c.setting.TimeoutContext(ctx)
	defer cancel()

	header, err := c.ethclient.HeaderByNumber(timeoutCtx, nil)
	if err != nil {
		return ZeroHash, err
	}
	return header.Hash(), nil
}

func (c *Client) BlockHeader(ctx context.Context) (*types.Header, error) {
	timeoutCtx, cancel := c.setting.TimeoutContext(ctx)
	defer cancel()

	header, err := c.ethclient.HeaderByNumber(timeoutCtx, nil)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (c *Client) PendingTxCount(ctx context.Context) (uint, error) {
	timeoutCtx, cancel := c.setting.TimeoutContext(ctx)
	defer cancel()

	return c.ethclient.PendingTransactionCount(timeoutCtx)
}

func (c *Client) TxpoolPendingTxCount(ctx context.Context) (uint64, error) {
	timeoutCtx, cancel := c.setting.TimeoutContext(ctx)
	defer cancel()

	var result map[string]hexutil.Uint
	err := c.c.CallContext(timeoutCtx, &result, "txpool_status")
	if err != nil {
		return 0, errors.Wrap(err, "err CallContext")
	}

	return hexutil.DecodeUint64(result["pending"].String())
}

func (c *Client) TxCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	timeoutCtx, cancel := c.setting.TimeoutContext(ctx)
	defer cancel()

	return c.ethclient.TransactionCount(timeoutCtx, blockHash)
}

func (c *Client) Nonce(ctx context.Context) (uint64, error) {
	timeoutCtx, cancel := c.setting.TimeoutContext(ctx)
	defer cancel()

	priv := GetPriv(ctx)

	if priv == nil {
		if c.setting.PrivKey == nil {
			return 0, errors.New("should set private key")
		}
		priv = c.setting.PrivKey
	}

	account := crypto.PubkeyToAddress(priv.PublicKey)
	return c.ethclient.NonceAt(timeoutCtx, account, nil)
}

func (c *Client) SendETH(ctx context.Context, nonce uint64, to common.Address, amount *big.Int) (common.Hash, error) {
	timeoutCtx, cancel := c.setting.TimeoutContext(ctx)
	defer cancel()

	var (
		tx     = types.NewTransaction(nonce, to, amount, uint64(EthSendingGasLimit), c.setting.GasPrice, nil)
		priv = GetPriv(ctx)
		signer = types.HomesteadSigner{}
		hash   common.Hash
	)

	if priv == nil {
		if c.setting.PrivKey == nil {
			return hash, errors.New("should set private key")
		}
		priv = c.setting.PrivKey
	}

	sig, err := crypto.Sign(signer.Hash(tx).Bytes(), priv)
	if err != nil {
		return hash, errors.Wrap(err, "err Sign")
	}

	signedTx, err := tx.WithSignature(signer, sig)
	if err != nil {
		return hash, errors.Wrap(err, "err WithSignature")
	}

	if err = c.ethclient.SendTransaction(timeoutCtx, signedTx); err != nil {
		return hash, errors.Wrap(err, "err SendTransaction")
	}

	return signedTx.Hash(), nil
}

func (c *Client) ConfirmTx(ctx context.Context, hash common.Hash) error {
	timeoutCtx, cancel := c.setting.TimeoutContext(ctx)
	defer cancel()

	return ConfirmTx(timeoutCtx, c.ethclient, hash)
}

func (c *Client) Receipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	timeoutCtx, cancel := c.setting.TimeoutContext(ctx)
	defer cancel()

	return c.ethclient.TransactionReceipt(timeoutCtx, hash)
}

func (c *Client) PrivKey() *ecdsa.PrivateKey {
	return c.setting.PrivKey
}

func (c *Client) GasPrice() *big.Int {
	return c.setting.GasPrice
}

func (c *Client) Ethclient() *ethclient.Client {
	return c.ethclient
}

func (c *Client) Setting() Setting {
	return c.setting
}

// func (c *Client) Record(ctx context.Context, token *erc20.MindenToken, demander common.Address, amount *big.Int, datetime uint32, index uint8) (*types.Transaction, error) {
// 	timeoutCtx, cancel := c.setting.TimeoutContext(ctx)
// 	defer cancel()

// 	return token.Record(timeoutCtx, c.setting.PrivKey, c.setting.GasPrice, nil, demander, amount, datetime, index)
// }

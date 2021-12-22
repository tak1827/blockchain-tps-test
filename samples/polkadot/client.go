package main

import (
	"context"

	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/tak1827/blockchain-tps-test/tps"
	// "github.com/davecgh/go-spew/spew"
)

var (
	_ tps.Client = (*PolkaClient)(nil)

	meta        *types.Metadata
	genesisHash types.Hash
	rv          *types.RuntimeVersion
)

type PolkaClient struct {
	client client.Client
	rpc    *rpc.RPC
}

func NewClient(url string) (c PolkaClient, err error) {
	if url == "" {
		url = config.Default().RPCURL
	}

	c.client, err = client.Connect(url)
	if err != nil {
		return
	}

	c.rpc, err = rpc.NewRPC(c.client)
	if err != nil {
		return
	}

	meta, err = c.rpc.State.GetMetadataLatest()
	if err != nil {
		return
	}

	genesisHash, err = c.rpc.Chain.GetBlockHash(0)
	if err != nil {
		return
	}

	rv, err = c.rpc.State.GetRuntimeVersionLatest()
	if err != nil {
		return
	}

	return
}

func (c PolkaClient) LatestBlockHeight(ctx context.Context) (uint64, error) {
	res, err := c.rpc.Chain.GetBlockLatest()
	if err != nil {
		return 0, err
	}
	return uint64(res.Block.Header.Number), nil
}

func (c PolkaClient) CountTx(ctx context.Context, height uint64) (int, error) {
	hash, err := c.rpc.Chain.GetBlockHash(height)
	if err != nil {
		return 0, err
	}

	res, err := c.rpc.Chain.GetBlock(hash)
	if err != nil {
		return 0, err
	}

	return len(res.Block.Extrinsics), nil
}

func (c PolkaClient) CountPendingTx(ctx context.Context) (int, error) {
	res, err := c.rpc.Author.PendingExtrinsics()
	if err != nil {
		return 0, err
	}
	return len(res), nil
}

func (c PolkaClient) Nonce(ctx context.Context, pubKey string) (uint64, error) {
	key, err := types.CreateStorageKey(meta, "System", "Account", types.MustHexDecodeString(pubKey), nil)
	if err != nil {
		return 0, err
	}

	var accountInfo types.AccountInfo
	ok, err := c.rpc.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		return 0, err
	}

	return uint64(accountInfo.Nonce), nil
}

func (c *PolkaClient) SendTx(ctx context.Context, seed string, nonce uint64, to types.MultiAddress, amount int64) (hash types.Hash, err error) {
	call, err := types.NewCall(meta, "Balances.transfer", to, types.NewUCompactFromUInt(uint64(amount)))
	if err != nil {
		return
	}

	ext := types.NewExtrinsic(call)

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(nonce),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	key, err := signature.KeyringPairFromSecret(seed, 0)
	if err != nil {
		return
	}

	if err = ext.Sign(key, o); err != nil {
		return
	}

	hash, err = c.rpc.Author.SubmitExtrinsic(ext)
	return
}

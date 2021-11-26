package main

import (
	"context"

	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc"
	"github.com/tak1827/blockchain-tps-test/tps"
)

var (
	_ tps.Client = (*PolkaClient)(nil)

	accNums = make(map[string]uint64, 3)
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

func (c PolkaClient) Nonce(ctx context.Context, address string) (uint64, error) {
	return 0, nil
}

package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	txtypes "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	// "github.com/davecgh/go-spew/spew"
	"github.com/tak1827/blockchain-tps-test/tps"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"google.golang.org/grpc"
)

const (
	DefalultRPCURI = "tcp://localhost:26657"
)

var (
	_ tps.Client = (*CosmosClient)(nil)

	ChainID              = os.Getenv("CHAINID")
	Denom                = os.Getenv("DENOM")
	AccountAddressPrefix = os.Getenv("ACCOUNT_PREFIX")

	accNums = make(map[string]uint64)
)

type CosmosClient struct {
	conn *grpc.ClientConn

	clientHTTP *rpchttp.HTTP
	authClient authtypes.QueryClient

	cdc      *codec.ProtoCodec
	txConfig client.TxConfig
}

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountAddressPrefix+sdk.PrefixPublic)
	config.Seal()
}

func NewClient(rpcURI string) (CosmosClient, error) {
	var (
		c      = CosmosClient{}
		encCfg = simapp.MakeTestEncodingConfig()
		err    error
	)

	if rpcURI == "" {
		rpcURI = DefalultRPCURI
	}

	if c.clientHTTP, err = rpchttp.New(rpcURI, "/websocket"); err != nil {
		return c, err
	}

	if c.conn, err = grpc.Dial("127.0.0.1:9090", grpc.WithInsecure()); err != nil {
		return c, err
	}

	c.authClient = authtypes.NewQueryClient(c.conn)

	c.cdc = codec.NewProtoCodec(encCfg.InterfaceRegistry)
	c.txConfig = txtypes.NewTxConfig(c.cdc, txtypes.DefaultSignModes)

	return c, nil
}

func (c CosmosClient) LatestBlockHeight(ctx context.Context) (uint64, error) {
	res, err := c.clientHTTP.Block(ctx, nil)
	if err != nil {
		return 0, err
	}
	return uint64(res.Block.Header.Height), nil
}

func (c CosmosClient) CountTx(ctx context.Context, height uint64) (int, error) {
	h := int64(height)
	res, err := c.clientHTTP.Block(ctx, &h)
	if err != nil {
		return 0, err
	}
	return len(res.Block.Data.Txs), nil
}

func (c CosmosClient) CountPendingTx(ctx context.Context) (int, error) {
	res, err := c.clientHTTP.UnconfirmedTxs(ctx, nil)
	if err != nil {
		return 0, err
	}
	return int(res.Total), nil
}

func (c CosmosClient) Nonce(ctx context.Context, address string) (uint64, error) {
	acc, err := c.Account(ctx, address)
	if err != nil {
		return 0, err
	}
	return acc.GetSequence(), nil
}

func (c CosmosClient) Account(ctx context.Context, address string) (acc authtypes.AccountI, err error) {
	req := &authtypes.QueryAccountRequest{Address: address}
	res, err := c.authClient.Account(ctx, req)
	if err != nil {
		return
	}

	if err = c.cdc.UnpackAny(res.GetAccount(), &acc); err != nil {
		return
	}

	return
}

func (c CosmosClient) Close() {
	c.conn.Close()
}

func PrivFromString(privStr string) (priv cryptotypes.PrivKey, err error) {
	priBytes, err := hex.DecodeString(privStr)
	if err != nil {
		return
	}
	priv = &secp256k1.PrivKey{Key: priBytes}
	return
}

func AccAddressFromPriv(priv cryptotypes.PrivKey) sdk.AccAddress {
	return sdk.AccAddress(priv.PubKey().Address().Bytes())
}

func AccAddressFromPrivString(privStr string) (string, error) {
	priv, err := PrivFromString(privStr)
	if err != nil {
		return "", err
	}
	return AccAddressFromPriv(priv).String(), nil
}

func (c *CosmosClient) BuildTx(msg sdk.Msg, priv cryptotypes.PrivKey, accSeq uint64) (authsigning.Tx, error) {
	var (
		txBuilder = c.txConfig.NewTxBuilder()
		accNum    = accNums[AccAddressFromPriv(priv).String()]
	)

	err := txBuilder.SetMsgs(msg)
	if err != nil {
		return nil, err
	}
	txBuilder.SetGasLimit(uint64(flags.DefaultGasLimit))

	// First round: we gather all the signer infos. We use the "set empty signature" hack to do that.
	if err = txBuilder.SetSignatures(signing.SignatureV2{
		PubKey: priv.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  c.txConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: accSeq,
	}); err != nil {
		return nil, err
	}

	// Second round: all signer infos are set, so each signer can sign.
	signerData := xauthsigning.SignerData{
		ChainID:       ChainID,
		AccountNumber: accNum,
		Sequence:      accSeq,
	}
	sigV2, err := tx.SignWithPrivKey(
		c.txConfig.SignModeHandler().DefaultMode(), signerData,
		txBuilder, priv, c.txConfig, accSeq)
	if err != nil {
		return nil, err
	}
	if err = txBuilder.SetSignatures(sigV2); err != nil {
		return nil, err
	}

	return txBuilder.GetTx(), nil
}

func (c *CosmosClient) SendTx(ctx context.Context, privStr string, seq uint64, to sdk.AccAddress, amount int64) (*ctypes.ResultBroadcastTx, error) {
	priv, err := PrivFromString(privStr)
	if err != nil {
		return nil, err
	}

	var (
		from  = AccAddressFromPriv(priv)
		coins = sdk.NewCoins(sdk.NewInt64Coin(Denom, amount))
		msg   = banktypes.NewMsgSend(from, to, coins)
	)
	tx, err := c.BuildTx(msg, priv, seq)
	if err != nil {
		return nil, err
	}

	txBytes, err := c.txConfig.TxEncoder()(tx)
	if err != nil {
		return nil, err
	}

	res, err := c.clientHTTP.BroadcastTxSync(ctx, txBytes)
	// Note: In async case, response is returnd before TxCheck
	// res, err := c.clientHTTP.BroadcastTxAsync(ctx, txBytes)
	if errRes := client.CheckTendermintError(err, txBytes); errRes != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, fmt.Errorf("code: %d, log: %s, codespace: %s\n", res.Code, res.Log, res.Codespace)
	}

	return res, nil
}

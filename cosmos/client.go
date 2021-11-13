package main

import (
	// "errors"
	"encoding/hex"
	"fmt"
	"context"

	// secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/davecgh/go-spew/spew"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	coregrpc "github.com/tendermint/tendermint/rpc/grpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	// "github.com/cosmos/cosmos-sdk/codec/legacy"
)

type Client interface {
	LatestBlockHeight(context.Context) (int, error)
	CountTx(context.Context, int) (int, error)
	// GetNonce(context.Context) (int, error)
	// SendTx() error
}

var (
	_ Client = (*CosmosClient)(nil)
)

type CosmosClient struct {
	clientHTTP *rpchttp.HTTP
	clientGRPC coregrpc.BroadcastAPIClient
}

const (
	DefalultRPCURI = "tcp://localhost:26657"
	DefaultGRPCURI = "tcp://localhost:36656"
	HomeDir = "/Users/tak/Documents/minden/cosmminden/.chaindata"
	ChainID = "mchain"
	KeyringBackend = "test"
)


func New(rpcURI string, grpcURI string) (CosmosClient, error) {
	// s, ok := setting.(Setting); !ok {
	// 	return cosmClient, fmt.Errorf("unexpected setting type, setting: %+v", s)
	// }
	// initClientCtx := client.Context{}.
	// 	WithJSONMarshaler(encodingConfig.Marshaler).
	// 	WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
	// 	WithTxConfig(encodingConfig.TxConfig).
	// 	WithLegacyAmino(encodingConfig.Amino).
	// 	WithInput(os.Stdin).
	// 	WithAccountRetriever(types.AccountRetriever{}).
	// 	WithBroadcastMode(flags.BroadcastBlock).
	// 	WithHomeDir(defaultNodeHome).
	// 	WithViper("")
	var (
		c = CosmosClient{}
		err error
	)


	if rpcURI == "" {
		rpcURI = DefalultRPCURI
	}

	if c.clientHTTP, err = rpchttp.New(rpcURI, "/websocket"); err != nil {
		return c, err
	}

	// clientCtx := DefaultContext()

	// c.ctx = clientCtx.WithClient(c.clientHTTP)

	if grpcURI == "" {
		grpcURI = DefaultGRPCURI
	}

	c.clientGRPC = coregrpc.StartGRPCClient(grpcURI)

	return c, nil
}

func (c CosmosClient) LatestBlockHeight(ctx context.Context) (int, error) {
	res, err := c.clientHTTP.Block(ctx, nil)
	if err != nil {
		return 0, err
	}
	return int(res.Block.Header.Height), nil
}

func (c CosmosClient) CountTx(ctx context.Context, height int) (int, error) {
	query := fmt.Sprintf("tx.height = %d", height)
	res, err := c.clientHTTP.TxSearch(ctx, query, false, nil, nil, "asc")
	if err != nil {
		return 0, err
	}
	return int(res.TotalCount), nil
}

// func (c CosmosClient) SendCoin(to sdk.AccAddress, amount sdk.Coins, nonce uint64) error {
// 	msg := sdk.NewMsgCreateConsumer(c.ctx.GetFromAddress(), to, amount)
// 	if err := msg.ValidateBasic(); err != nil {
// 		return err
// 	}

// 	txf := c.TxFactory(nonce)

// 	if err = tx.BroadcastTx(c.ctx, txf, msg); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func DefaultContext(from string) client.Context {
// 	cmd := &cobra.Command{}
// 	cmd.Flags().String(flags.HomeDir, HomeDir, "")
// 	cmd.Flags().String(flags.FlagChainID, ChainID, "")
// 	cmd.Flags().String(flags.FlagKeyringBackend, KeyringBackend, "")
// 	cmd.Flags().String(flags.FlagBroadcastMode, flags.BroadcastAsync, "")
// 	cmd.Flags().Bool(flags.FlagSkipConfirmation, true, "")
// 	cmd.Flags().String(flags.From, from, "")

// 	return clinet.GetClientTxContext(cmd)
// }

func (c CosmosClient) SendTx() error {
	priStr := "3b2eb70cf00a779c4bfed132e0fb3f7f982013132d86e344b0e97c7507d0d7a4"
  priBytes, err := hex.DecodeString(priStr)
  if err != nil {
  	return err
  }
  var priv cryptotypes.PrivKey
  priv = &secp256k1.PrivKey{Key: priBytes}
  // // priv2, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), priBytes)
  // spew.Dump("priv2", priv2)
  // priv, err := legacy.PrivKeyFromBytes(priBytes)
  // if err != nil {
  // 	return err
  // }
  // spew.Dump("priv", priv)
  // priv2, _ := priv.(cryptotypes.PrivKey)
  privs := []cryptotypes.PrivKey{priv}
  accNums:= []uint64{1} // The accounts' account numbers
  accSeqs:= []uint64{4} // The accounts' sequence numbers

  // pub, err := cryptocodec.FromTmPubKeyInterface(priv.PubKey())
  // if err != nil {
  // 	return err
  // }

  config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("minden", "minden"+sdk.PrefixPublic)
	config.Seal()

  from := sdk.AccAddress(priv.PubKey().Address().Bytes())

  toStr := "minden1hp5cugfmh8e57pqggaukmjhrsfx2lxt9nthf83"
 //  to, err := sdk.AccAddressFromBech32(toStr)
	// if err != nil {
	// 	return err
	// }
	bz, err := sdk.GetFromBech32(toStr, "minden")
	if err != nil {
		return err
	}

	to := sdk.AccAddress(bz)

	encCfg := simapp.MakeTestEncodingConfig()
  txBuilder := encCfg.TxConfig.NewTxBuilder()
  msg := banktypes.NewMsgSend(from, to, sdk.NewCoins(sdk.NewInt64Coin("mtc", 100)))

  err = txBuilder.SetMsgs(msg)
  if err != nil {
      return err
  }
  txBuilder.SetGasLimit(1000000)
  // txBuilder.SetFeeAmount(sdk.NewInt64Coin("mtc", 1000))

  // First round: we gather all the signer infos. We use the "set empty
  // signature" hack to do that.
  var sigsV2 []signing.SignatureV2
  for i, priv := range privs {
      sigV2 := signing.SignatureV2{
          PubKey: priv.PubKey(),
          Data: &signing.SingleSignatureData{
              SignMode:  encCfg.TxConfig.SignModeHandler().DefaultMode(),
              Signature: nil,
          },
          Sequence: accSeqs[i],
      }

      sigsV2 = append(sigsV2, sigV2)
  }
  err = txBuilder.SetSignatures(sigsV2...)
  if err != nil {
      return err
  }

  // Second round: all signer infos are set, so each signer can sign.
  sigsV2 = []signing.SignatureV2{}
  for i, priv := range privs {
      signerData := xauthsigning.SignerData{
          ChainID:       ChainID,
          AccountNumber: accNums[i],
          Sequence:      accSeqs[i],
      }
      sigV2, err := tx.SignWithPrivKey(
          encCfg.TxConfig.SignModeHandler().DefaultMode(), signerData,
          txBuilder, priv, encCfg.TxConfig, accSeqs[i])
      if err != nil {
          return err
      }

      sigsV2 = append(sigsV2, sigV2)
  }
  err = txBuilder.SetSignatures(sigsV2...)
  if err != nil {
      return err
  }

  // Generated Protobuf-encoded bytes.
  txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
  if err != nil {
      return err
  }

  // // Generate a JSON string.
  // txJSONBytes, err := encCfg.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
  // if err != nil {
  //     return err
  // }
  // txJSON := string(txJSONBytes)

  // res, err := c.clientHTTP.BroadcastTxAsync(context.Background(), txBytes)
  res, err := c.clientHTTP.BroadcastTxSync(context.Background(), txBytes)
  if errRes := client.CheckTendermintError(err, txBytes); errRes != nil {
  	spew.Dump(errRes)
  	panic("hoge!")
	}

	spew.Dump(res)

	return nil
}


// func (c CosmosClient) TxFactory(accountNumber uint64) tx.Factory {
// 	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
// 	fs.Uint16(flags.FlagAccountNumber, accountNumber, "")
// 	return tx.NewFactoryCLI(c.ctx, fs)
// }


func main() {
	// ctx := context.Background()
	c, err := New("", "")
	if err != nil {
		panic(err.Error())
	}
	// res, err := c.LatestBlockHeight(ctx)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// res2, err := c.CountTx(ctx, res)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// spew.Dump(res2)
	err = c.SendTx()
	if err != nil {
		panic(err.Error())
	}
}

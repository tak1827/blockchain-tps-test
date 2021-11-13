package eth

import (
	"math/big"
	"testing"

	// "github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/min-sys/tracking-contract/cli/constant"
)

func TestBlockNumber(t *testing.T) {
	client, err := newClient(nil, defaultSetting(t))
	if err != nil {
		t.Fatal("err create NewClient", err)
	}

	if _, err := client.BlockNumer(nil); err != nil {
		t.Error("err BlockNumer", err)
	}
}

func TestTxCount(t *testing.T) {
	client, err := newClient(nil, defaultSetting(t))
	if err != nil {
		t.Fatal("err create NewClient", err)
	}

	hash, err := client.BlockHash(nil)
	if err != nil {
		t.Error("err BlockNumer", err)
	}

	if _, err := client.TxCount(nil, hash); err != nil {
		t.Error("err TxCount", err)
	}
}

func TestTxpoolPendingTxCount(t *testing.T) {
	client, err := newClient(nil, defaultSetting(t))
	if err != nil {
		t.Fatal("err create NewClient", err)
	}

	if _, err = client.TxpoolPendingTxCount(nil); err != nil {
		t.Error("err TxpoolPendingTxCount", err)
	}
}

func TestNonce(t *testing.T) {
	client, err := newClient(nil, defaultSetting(t))
	if err != nil {
		t.Fatal("err NewClient", err)
	}

	if _, err := client.Nonce(nil); err != nil {
		t.Error("err Nonce", err)
	}
}

func TestSendETH(t *testing.T) {
	client, err := newClient(nil, defaultSetting(t))
	if err != nil {
		t.Fatal("err create NewClient", err)
	}

	to, err := GenerateAddr()
	if err != nil {
		t.Fatal("err GenerateAddr", err)
	}

	amount := ToWei(1.0, 18)

	nonce, err := client.Nonce(nil)
	if err != nil {
		t.Fatal("err Nonce", err)
	}

	txHash, err := client.SendETH(nil, nonce, to, amount)
	if err != nil {
		t.Error("err SendETH", err)
	}

	receipt, err := client.Receipt(nil, txHash)
	if err != nil {
		t.Error("err SendETH", err)
	}
	if receipt.Status != 1 {
		t.Errorf("transaction is failed, receipt: %+v", receipt)
	}
}

func defaultSetting(t *testing.T) Setting {
	priv, err := crypto.HexToECDSA(constant.TestPrivKey)
	if err != nil {
		t.Fatal("err HexToECDSA", err)
	}
	return Setting{
		Endpoint: constant.TestEndpoint,
		Timeout:  constant.DefaultTimeout,
		PrivKey:  priv,
		GasPrice: big.NewInt(int64(constant.DefaultGasPrice)),
	}
}

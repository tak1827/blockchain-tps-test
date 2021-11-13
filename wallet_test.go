package main

import (
	"testing"

	"github.com/min-sys/tracking-contract/cli/eth"
)

var (
	TestEndpoint = "http://127.0.0.1:8545"
	TestPrivKey  = "58900163de10a0ffe2d4c3224faf8da4d45727cb47114eb43dd754b60ab70cb7"
	// 0xE3b0DE0E4CA5D3CB29A9341534226C4D31C9838f
	TestPrivKey2 = "d1c71e71b06e248c8dbe94d49ef6d6b0d64f5d71b1e33a0f39e14dadb070304a"
)

func TestRotatePriv(t *testing.T) {
	var (
		keys = []string{
			PrivKey,
			PrivKey2,
		}
	)
	client, err := eth.NewClient(nil, eth.WithEndpoint(TestEndpoint), eth.WithPrivKey(TestPrivKey))
	if err != nil {
		t.Fatal("err NewClient: ", err)
	}

	wallet, err := NewWallet(&client, keys)
	if err != nil {
		t.Fatal("err NewWallet: ", err)
	}

	if g, w := wallet.RotatePriv(), wallet.privs[0]; g != w {
		t.Errorf("got: %v, want: %d", g, w)
	}

	if g, w := wallet.RotatePriv(), wallet.privs[1]; g != w {
		t.Errorf("got: %v, want: %d", g, w)
	}

	if g, w := wallet.RotatePriv(), wallet.privs[0]; g != w {
		t.Errorf("got: %v, want: %d", g, w)
	}
}

func TestIncrementNonce(t *testing.T) {
	var (
		keys = []string{
			PrivKey,
			PrivKey2,
		}
	)
	client, err := eth.NewClient(nil, eth.WithEndpoint(TestEndpoint), eth.WithPrivKey(TestPrivKey))
	if err != nil {
		t.Fatal("err NewClient: ", err)
	}

	wallet, err := NewWallet(&client, keys)
	if err != nil {
		t.Fatal("err NewWallet: ", err)
	}

	currentNonce := wallet.CurrentNonce(wallet.privs[0])

	if g, w := wallet.IncrementNonce(wallet.privs[0]), currentNonce; g != w {
		t.Errorf("got: %v, want: %d", g, w)
	}

	if g, w := wallet.IncrementNonce(wallet.privs[0]), currentNonce+1; g != w {
		t.Errorf("got: %v, want: %d", g, w)
	}
}

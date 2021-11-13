package main

import (
	"testing"

	"github.com/min-sys/tracking-contract/cli/constant"
	"github.com/min-sys/tracking-contract/cli/eth"
	// "github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func Test_ResetNonce(t *testing.T) {
	client, err := eth.NewClient(nil, eth.WithEndpoint(constant.TestEndpoint), eth.WithPrivKey(constant.TestPrivKey))
	if err != nil {
		t.Fatal(err, "err create NewClient")
	}

	queue := NewQueue(0)

	to, err := eth.GenerateAddr()
	if err != nil {
		t.Fatal("err GenerateAddr", err)
	}

	nonce, err := client.Nonce(nil)
	if err != nil {
		t.Fatal("err Nonce", err)
	}

	amount := eth.ToWei(1.0, 18)

	task := ETHSendingTask{
		to:     to,
		amount: amount,
	}

	if err = task.Do(&client, nonce, &queue, Logger{}); err != nil {
		t.Error("err Do", err)
	}

	if err = task.Do(&client, nonce, &queue, Logger{}); err != nil {
		if !errors.Is(err, ErrWrongNonce) {
			t.Errorf("unexpected error, err: %+v", err)
		}
	}

	if err = task.Do(&client, nonce+2, &queue, Logger{}); err != nil {
		if !errors.Is(err, ErrWrongNonce) {
			t.Errorf("unexpected error, err: %+v", err)
		}
	}

	if err = task.Do(&client, nonce+1, &queue, Logger{}); err != nil {
		t.Error("err Do", err)
	}
}

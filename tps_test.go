package main

// import (
// 	"sync/atomic"
// 	"testing"
// 	"time"

// 	"github.com/min-sys/tracking-contract/cli/eth"
// )

// const (
// 	RopstenEndpoint = "https://ropsten.infura.io/v3/b3dd59dcade64d8d9d7b5dbfe403c152"
// )

// func TestStartTPSMeasuring(t *testing.T) {
// 	client, err := eth.NewClient(nil, eth.WithEndpoint(RopstenEndpoint))
// 	if err != nil {
// 		t.Fatal(err, "err create NewClient")
// 	}

// 	var closing uint32

// 	go func() {
// 		time.Sleep(60 * time.Second)
// 		atomic.AddUint32(&closing, 1)
// 	}()

// 	if err = StartTPSMeasuring(&client, &closing); err != nil {
// 		t.Error(err, "err StartTPSMeasuring")
// 	}
// }

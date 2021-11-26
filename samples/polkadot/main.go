package main

import (
	"github.com/davecgh/go-spew/spew"
)

func main() {
	client, err := NewClient("")
	if err != nil {
		panic(err)
	}

	height, err := client.LatestBlockHeight(nil)
	if err != nil {
		panic(err)
	}
	spew.Dump(height)

	coutx, err := client.CountTx(nil, height)
	if err != nil {
		panic(err)
	}
	spew.Dump(coutx)

	count, err := client.CountPendingTx(nil)
	if err != nil {
		panic(err)
	}
	spew.Dump(count)
}

package tps

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/pkg/errors"
)

type Wallet struct {
	privs      []string
	nonces     map[string]*Nonce
	roteteSlot uint32
}

func NewWallet(ctx context.Context, client Client, privs []string, addrs []string) (w Wallet, err error) {
	if len(privs) != len(addrs) {
		err = fmt.Errorf("the length of privs and addrs should be same")
		return
	}

	w.privs = privs
	w.nonces = make(map[string]*Nonce, len(privs))

	for i := range privs {
		var nonce Nonce
		if nonce, err = NewNonce(ctx, client, addrs[i]); err != nil {
			err = errors.Wrap(err, "err NewNonce")
			return
		}
		w.nonces[privs[i]] = &nonce
	}

	return
}

func (w *Wallet) RotatePriv() string {
	slot := atomic.LoadUint32(&w.roteteSlot)
	atomic.AddUint32(&w.roteteSlot, 1)
	return w.privs[slot%uint32(len(w.privs))]
}

func (w *Wallet) Priv(index int) string {
	i := index % len(w.privs)
	return w.privs[i]
}

func (w *Wallet) IncrementNonce(priv string) uint64 {
	return w.nonces[priv].Increment()
}

func (w *Wallet) CurrentNonce(priv string) uint64 {
	return w.nonces[priv].Current()
}

func (w *Wallet) RecetNonce(priv string, nonce uint64) {
	w.nonces[priv].Reset(nonce)
}

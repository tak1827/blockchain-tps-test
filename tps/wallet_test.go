package tps

import (
	"testing"
)

func TestRotatePriv(t *testing.T) {
	var (
		privs    = []string{"key1", "key2"}
		nonces   = []Nonce{Nonce{}, Nonce{}}
		nonceMap = map[string]*Nonce{
			"key1": &nonces[0],
			"key2": &nonces[1],
		}
		wallet = Wallet{
			privs:  privs,
			nonces: nonceMap,
		}
	)

	if g, w := wallet.RotatePriv(), wallet.privs[0]; g != w {
		t.Errorf("got: %v, want: %s", g, w)
	}

	if g, w := wallet.RotatePriv(), wallet.privs[1]; g != w {
		t.Errorf("got: %v, want: %s", g, w)
	}

	if g, w := wallet.RotatePriv(), wallet.privs[0]; g != w {
		t.Errorf("got: %v, want: %s", g, w)
	}
}

func TestIncrementNonce(t *testing.T) {
	var (
		privs    = []string{"key1", "key2"}
		nonces   = []Nonce{Nonce{}, Nonce{}}
		nonceMap = map[string]*Nonce{
			"key1": &nonces[0],
			"key2": &nonces[1],
		}
		wallet = Wallet{
			privs:  privs,
			nonces: nonceMap,
		}
	)

	currentNonce := wallet.CurrentNonce(wallet.privs[0])

	if g, w := wallet.IncrementNonce(wallet.privs[0]), currentNonce; g != w {
		t.Errorf("got: %v, want: %d", g, w)
	}

	if g, w := wallet.IncrementNonce(wallet.privs[0]), currentNonce+1; g != w {
		t.Errorf("got: %v, want: %d", g, w)
	}
}

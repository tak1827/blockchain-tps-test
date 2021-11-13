package eth

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

type ClientOption interface {
	Apply(*Setting) error
}

type EndpointOpt string

func (o EndpointOpt) Apply(s *Setting) error {
	s.Endpoint = string(o)
	return nil
}
func WithEndpoint(endpoint string) ClientOption {
	return EndpointOpt(endpoint)
}

type TimeoutOpt time.Duration

func (o TimeoutOpt) Apply(s *Setting) error {
	s.Timeout = time.Duration(o)
	return nil
}
func WithTimeout(timeout time.Duration) ClientOption {
	return TimeoutOpt(timeout)
}

type PrivKeyOpt string

func (o PrivKeyOpt) Apply(s *Setting) error {
	priv, err := crypto.HexToECDSA(string(o))
	if err != nil {
		return errors.Wrap(err, "err HexToECDSA")
	}
	s.PrivKey = priv
	return nil
}
func WithPrivKey(priv string) PrivKeyOpt {
	return PrivKeyOpt(priv)
}

type GasPriceOpt int64

func (o GasPriceOpt) Apply(s *Setting) error {
	s.GasPrice = big.NewInt(int64(o))
	return nil
}
func WithGasPrice(gasPrice int64) GasPriceOpt {
	return GasPriceOpt(gasPrice)
}

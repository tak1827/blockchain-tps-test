package eth

import (
	"context"
	"os"

	"github.com/min-sys/tracking-contract/cli/constant"
	"github.com/pkg/errors"
)

func NewClient(ctx context.Context, clinetOptions ...ClientOption) (Client, error) {

	var opts []ClientOption
	if addr := os.Getenv("ETH_ENDPOINT_HOST"); addr != "" {
		opts = append(opts, WithEndpoint(addr))
	}

	if privKey := os.Getenv("ETH_PRIVATE_KEY"); privKey != "" {
		opts = append(opts, WithPrivKey(privKey))
	}

	opts = append(opts, WithTimeout(constant.DefaultTimeout), WithGasPrice(constant.DefaultGasPrice))
	opts = append(opts, clinetOptions...)

	var setting Setting
	for _, o := range opts {
		if err := o.Apply(&setting); err != nil {
			return Client{}, errors.Wrap(err, "err Apply")
		}
	}

	return newClient(ctx, setting)
}

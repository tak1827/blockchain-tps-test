module github.com/tak1827/blockchain-tps-test/samples/polkadot

go 1.16

require (
	github.com/centrifuge/go-substrate-rpc-client/v4 v4.0.0
	github.com/cosmos/cosmos-sdk v0.44.5
	github.com/davecgh/go-spew v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/tak1827/blockchain-tps-test v0.0.0-20211126013655-bd8892b030ee
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

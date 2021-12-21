module github.com/tak1827/blockchain-tps-test/samples/cosmos

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.44.3
	github.com/pkg/errors v0.9.1
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/tak1827/blockchain-tps-test v0.0.0-20211221091226-2e9468abc68f
	github.com/tendermint/tendermint v0.34.14
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa // indirect
	golang.org/x/net v0.0.0-20211108170745-6635138e15ea // indirect
	golang.org/x/sys v0.0.0-20211108224332-cbcd623f202e // indirect
	google.golang.org/genproto v0.0.0-20211104193956-4c6863e31247 // indirect
	google.golang.org/grpc v1.42.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

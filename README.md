# cosmos-chain-tps-test
TPS testing tool of blockchain, like Cosmos, Polkadot and Ethereum.

# How to use
The tps module is the core. Able to mesure tps such as chain wihch satisfy bellow interface.
```go
type Client interface {
	LatestBlockHeight(context.Context) (uint64, error)
	CountTx(context.Context, uint64) (int, error)
	CountPendingTx(context.Context) (int, error)
	Nonce(context.Context, string) (uint64, error)
}
```
The sample codes are in`samples` directory.


# Results
## Cosmos
**Condition**
- single node
Using the auto generate project via [starport](https://github.com/tendermint/starport)
- Mesured by Macbook pro 2020
  - CPU: 1.4 GHz Quad-Core Intel Core i5
  - MEM: 16 GB
**Result**
TPS: around `750`
```sh
gtime -f '%P %Uu %Ss %er %MkB %C' "time" ./recorder
------------------------------------------------------------------------------------
⛓  5013 th Block Mind! txs(4329), total txs(4329), TPS(704.25), pendig txs(5515)
------------------------------------------------------------------------------------
⛓  5014 th Block Mind! txs(4257), total txs(8586), TPS(714.43), pendig txs(5446)
------------------------------------------------------------------------------------
⛓  5015 th Block Mind! txs(4476), total txs(13062), TPS(729.69), pendig txs(5872)
------------------------------------------------------------------------------------
⛓  5016 th Block Mind! txs(4365), total txs(17427), TPS(767.86), pendig txs(5343)
------------------------------------------------------------------------------------
⛓  5017 th Block Mind! txs(4524), total txs(21951), TPS(767.80), pendig txs(5539)
------------------------------------------------------------------------------------
⛓  5018 th Block Mind! txs(4329), total txs(26280), TPS(759.34), pendig txs(5538)
------------------------------------------------------------------------------------
⛓  5019 th Block Mind! txs(4194), total txs(30474), TPS(752.33), pendig txs(5531)
------------------------------------------------------------------------------------
⛓  5020 th Block Mind! txs(4486), total txs(34960), TPS(772.76), pendig txs(5403)
------------------------------------------------------------------------------------
⛓  5021 th Block Mind! txs(3992), total txs(38952), TPS(760.24), pendig txs(4946)
       60.17 real        96.64 user         5.43 sys
169% 96.64u 5.45s 60.18r 64188kB time ./recorder
```

## Polkadot
**Condition**
- single node
Using [template node](https://github.com/substrate-developer-hub/substrate-front-end-template)
- Mesured by Macbook pro 2020
  - CPU: 1.4 GHz Quad-Core Intel Core i5
  - MEM: 16 GB
**Result**
TPS: around `250`
```sh
2021/12/10 10:06:35 Connecting to ws://127.0.0.1:9944...
------------------------------------------------------------------------------------
⛓  14 th Block Mind! txs(1725), total txs(1725), TPS(242.51), pendig txs(324)
------------------------------------------------------------------------------------
⛓  15 th Block Mind! txs(1752), total txs(3477), TPS(279.21), pendig txs(134)
------------------------------------------------------------------------------------
⛓  16 th Block Mind! txs(1652), total txs(5129), TPS(269.28), pendig txs(240)
------------------------------------------------------------------------------------
⛓  17 th Block Mind! txs(1591), total txs(6720), TPS(275.43), pendig txs(109)
       31.26 real        57.74 user         1.41 sys
```

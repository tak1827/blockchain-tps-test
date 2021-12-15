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
  - Using the auto generate project via [starport](https://github.com/tendermint/starport)
- sending simple payment transaction
  - transaction sending response is returned when the transaction is taken in mempool
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
  - Using [template node](https://github.com/substrate-developer-hub/substrate-front-end-template)
- sending simple payment transaction
  - transaction sending response is returned when the transaction is taken in mempool
- Mesured by Macbook pro 2020
  - CPU: 1.4 GHz Quad-Core Intel Core i5
  - MEM: 16 GB

**Result**

TPS: around `500`
```sh
gtime -f '%P %Uu %Ss %er %MkB %C' "time" ./recorder
------------------------------------------------------------------------------------
⛓  2582 th Block Mind! txs(1470), total txs(11987), TPS(498.98), pendig txs(775)
------------------------------------------------------------------------------------
⛓  2583 th Block Mind! txs(1462), total txs(13449), TPS(510.68), pendig txs(416)
------------------------------------------------------------------------------------
⛓  2584 th Block Mind! txs(1433), total txs(14882), TPS(500.49), pendig txs(632)
------------------------------------------------------------------------------------
⛓  2585 th Block Mind! txs(1441), total txs(16323), TPS(492.75), pendig txs(802)
------------------------------------------------------------------------------------
⛓  2586 th Block Mind! txs(1435), total txs(17758), TPS(501.20), pendig txs(490)
------------------------------------------------------------------------------------
⛓  2587 th Block Mind! txs(1461), total txs(19219), TPS(494.94), pendig txs(679)
------------------------------------------------------------------------------------
⛓  2588 th Block Mind! txs(1451), total txs(20670), TPS(502.19), pendig txs(348)
------------------------------------------------------------------------------------
⛓  2589 th Block Mind! txs(1378), total txs(22048), TPS(494.98), pendig txs(480)
       31.26 real        57.74 user         1.41 sys
189% 57.74u 1.41s 31.27r 18632kB time ./recorder
```

# CryptoQuantGo

Re-write everything in Go

```
    +-----------------------+
    |   WebSocket Feeds     |  (BTC/USD, ETH/USD, etc.)
    +-----------------------+
               │
               ▼
    +----------------------+
    | Price Data Channel   |  (Raw price updates)
    +----------------------+
               │
               ▼
    +----------------------+
    | Data Processing      |  (Detects trade signals)
    +----------------------+
      │                   │
      ▼                   ▼
  +----------------+     +----------------+
  | Trade Signal   |     | Data Store     | (Postgres Logging)
  | Channel        |     | Channel        |
  +----------------+     +----------------+
            │
            ▼
    +----------------+
    |  Order Channel |
    +----------------+
            │
            ▼
    +----------------+
    |  Execution     |  (Executes orders)
    +----------------+
```

## Structure 

```
.
├── README.md
├── config
│   ├── exchange.go
│   ├── future_exchange_config.go
│   └── future_exchange_config_test.go
├── data
│   ├── binance_future_market_rest.go
│   ├── binance_future_market_rest_test.go
│   ├── db.go
│   └── db_test.go
├── engine
│   ├── engine.go
│   ├── engine_test.go
│   └── trade.go
├── go.mod
├── go.sum
├── internal
│   ├── channels_core.go
│   ├── types_core.go
│   ├── types_future_exchange.go
│   ├── types_future_exchange_test.go
│   ├── types_future_rest.go
│   ├── types_future_stream.go
│   ├── types_spot_exchange.go
│   └── types_spot_exchange_test.go
├── main.go
├── strategy
│   ├── calculation
│   │   ├── cointegration.go
│   │   ├── cointegration_test.go
│   │   ├── hurst.go
│   │   ├── hurst_test.go
│   │   ├── spread_kalman.go
│   │   ├── spread_ols.go
│   │   └── spread_ols_test.go
│   └── pairs
│       ├── README.md
│       ├── type_assets.go
│       ├── type_pairs.go
│       └── update_pairs.go
├── streams
│   ├── kline.go
│   └── kline_test.go
└── utils
    ├── float_queue.go
    ├── logger.go
    └── timeutils.go
```

`strategy/` is hidden from the master branch. It's never wise to reveal your strategy.
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
├── README.md
├── config/               # Config settings (API keys, DB config, etc.)
│   ├── config.go
├── data/                 # Database (Postgres storage) & data models
│   ├── models.go         # ✅ Define core structs (PriceData, TradeSignal, Order)
│   ├── store.go          # ✅ Postgres interaction (saving WebSocket data)
├── engine/               # Trading execution logic
│   ├── trade.go          # ✅ Handles trade execution
├── go.mod
├── main.go               # ✅ Initializes channels & starts goroutines
├── strategy/             # Trading strategy logic
│   ├── pairs.go          # ✅ Pair trading logic (uses TradeSignal channel)
├── streams/              # WebSocket streaming
│   ├── subscribe.go      # ✅ Connects to WebSockets, sends data to channels
├── utils/                # Utility functions (logs, timestamps, etc.)
│   ├── logger.go
│   ├── timeutils.go
└── internal/             # ✅ New package for shared data & channels
    ├── channels.go       # ✅ Defines all channels
    ├── types.go          # ✅ Defines PriceData, TradeSignal, Order structs
```

`strategy/` is hidden from the master branch. It's never wise to reveal your strategy.
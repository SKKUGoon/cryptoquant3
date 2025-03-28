# CryptoQuantGo

Re-write everything in Go

## Structure 

<details>
<summary>System Architecture Diagram</summary>

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

### 1. Data Management
- **Binance Futures Integration**: Direct connection to Binance Futures API for market data retrieval
- **Database Integration**: Local database connection for data persistence
- **Caching System**: Intelligent caching of market data to minimize API calls
- **Rate Limiting**: Built-in rate limit management for API requests

### 2. Asset Management
The engine maintains individual trading assets with the following features:
- Real-time price tracking
- Historic data management
- Automatic data updates (when data is older than 1 minute)
- Concurrent processing of multiple assets

### 3. Stream Management
Efficient websocket handling for real-time data:
- Grouped symbol connections (chunks of 5)
- Channel-based data distribution
- Automatic connection management
- Error recovery mechanisms

`strategy/` is hidden from the master branch. 

> It's never wise to reveal your strategy. - <i>Some Foolish men, 2025</i>

</details>

## Technical Details

### How to Begin

1. Prepare environment file. This file will be feeded into the program using `init`.

    ```bash:.env
    # # Database - Postgres
    USERNAME=...
    PASSWORD=...
    PG_HOST=...
    PG_PORT=...
    PG_NAME=...
    ```

2. Run the program
    ```bash
    go run .
    ```

### Key Components
- `Engine`: Main controller managing all operations
- `PairAssets`: Individual asset handlers
- `Pairs`: Trading pair managers
- `ChMap`: Channel mapping for data streams

### Data Flow
1. Websocket streams receive real-time market data
2. Data is distributed to relevant asset handlers
3. Assets process and broadcast updates to subscribed pairs
4. Pairs analyze data and generate trading signals

### Configuration
- Test mode support for development
- Configurable trading pairs
- Adjustable data update intervals
- Rate limit configurations

## Usage

```go
// Initialize the engine
engine := NewEngine()

// Optional: Enable test mode
engine.SetTestMode(true)

// Start the system
engine.StartStreamCh()    // Initialize data channels
engine.StartAssets()      // Start asset handlers
engine.StartPairs()       // Initialize trading pairs
engine.StartStream()      // Begin real-time data streaming
```

## Status Monitoring

The engine provides status monitoring through:
```go
status := engine.GetStatus()  // Returns current system status
```

## Dependencies
- Binance Futures API
- Local database system
- Internal trading strategy modules
- Market data streaming services

## Performance Considerations
- Concurrent processing using goroutines
- Buffered channels for data handling
- Optimized API call management
- Memory-efficient data structures

## Development Status
This is an active project under continuous development. Future enhancements may include:
- Buy/Sell Order
- Real time user stream
- Additional trading strategies
- Enhanced risk management
- Extended market coverage
- Advanced analytics capabilities
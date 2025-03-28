package engine

import (
	"log"
	"strings"

	"cryptoquant.com/m/config"
	"cryptoquant.com/m/data"
	"cryptoquant.com/m/internal"
	strategy "cryptoquant.com/m/strategy/pairs"
	"cryptoquant.com/m/streams"
)

type Engine struct {
	// Engine configuration
	FutureConfig *config.FutureTradeConfig

	// Data
	FutureMarket *data.BinanceFutureMarketData
	Database     *data.Database

	// Streaming data
	Pairs map[string]strategy.Pair
	ChMap map[string]chan internal.KlineData // Channels for kline data. Key is symbol

	// Engine status
	isTest        bool
	isStreaming   bool
	isCalculating bool
}

func NewEngine() *Engine {
	db, err := data.ConnectDB()
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	log.Println("Connected to database")

	exchangeRestApi := data.NewBinanceFutureMarketData()
	log.Println("Connected to binance-futures data source")

	// Future exchange information
	config, err := config.NewFutureTradeConfig()
	if err != nil {
		panic("Failed to connect to binance-futures exchange: " + err.Error())
	}
	log.Println("Connected to binance-futures exchange")

	// Setup the internal Engine's attributes
	exchangeRestApi.UpdateRateLimit(config.ExchangeInfo.GetRequestRateLimit())

	engine := &Engine{
		FutureConfig: config,
		Database:     db,
		FutureMarket: exchangeRestApi,
		Pairs:        make(map[string]strategy.Pair),
		ChMap:        make(map[string]chan internal.KlineData),
		isTest:       false,
	}

	return engine
}

func (e *Engine) SetTestMode(testMode bool) {
	e.isTest = testMode
	e.FutureConfig.SetTestMode(testMode)
}

func (e *Engine) StartChMap() {
	symbols := e.FutureConfig.GetAvailableSymbols()
	for _, symbol := range symbols {
		e.ChMap[symbol] = make(chan internal.KlineData, 5)
	}
}

func (e *Engine) PreparePairCalculation() {
	pairs := e.FutureConfig.CreatePair(e.isTest)

	for _, pair := range pairs {
		assets := strings.Split(pair, "-")
		asset1 := assets[0]
		asset2 := assets[1]

		// Create a new pair
		e.Pairs[pair] = *strategy.NewPair(asset1, asset2)

		// Set the channels for the pair
		e.Pairs[pair].Asset1.SetChannel(e.ChMap[asset1])
		e.Pairs[pair].Asset2.SetChannel(e.ChMap[asset2])

		e.FutureMarket.GetKlineData(asset1, "1m", 100)
	}
}

func (e *Engine) StartStream(done chan struct{}) {
	symbols := e.FutureConfig.GetAvailableSymbols()
	// Group symbols into chunks of 5
	const chunkSize = 5
	for i := 0; i < len(symbols); i += chunkSize {
		end := min(i+chunkSize, len(symbols))
		symbolGroup := symbols[i:end]

		// Process each group of symbols - Multiple queires
		go streams.SubscribeKlineMulti(symbolGroup, "1m", e.ChMap, done)
	}

	e.isStreaming = true
}

func (e *Engine) StopStream() {
	e.isStreaming = false

	for _, ch := range e.ChMap {
		close(ch)
	}
}

// ListenPair the engine. ListenPair all the pair streams
func (e *Engine) ListenPair(done chan struct{}) {
	for _, pair := range e.Pairs {
		go pair.Run(done)
	}
}

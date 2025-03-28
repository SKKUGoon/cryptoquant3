package engine

import (
	"fmt"
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
	PairAssets map[string]*strategy.PairAsset
	Pairs      map[string]*strategy.Pair
	ChMap      map[string]chan internal.KlineData // Channels for kline data. Key is symbol

	// Engine status
	done          chan struct{}
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
		FutureConfig:  config,
		Database:      db,
		FutureMarket:  exchangeRestApi,
		PairAssets:    make(map[string]*strategy.PairAsset),
		Pairs:         make(map[string]*strategy.Pair),
		ChMap:         make(map[string]chan internal.KlineData),
		done:          make(chan struct{}),
		isTest:        false,
		isStreaming:   false,
		isCalculating: false,
	}

	return engine
}

func (e *Engine) SetTestMode(testMode bool) {
	e.isTest = testMode
	e.FutureConfig.SetTestMode(testMode)
}

// StartStreamCh starts the stream channels for the available symbols.
// It transfers data from the websocket connections to PairAssets.
func (e *Engine) StartStreamCh() {
	symbols := e.FutureConfig.GetAvailableSymbols()
	for _, symbol := range symbols {
		e.ChMap[symbol] = make(chan internal.KlineData, 5)
	}
}

// StartPairAssets starts the PairAssets for the available symbols.
// Receives data from the stream channels and transfers and broadcasts to the Pairs.
func (e *Engine) StartPairAssets() {
	symbols := e.FutureConfig.GetAvailableSymbols()
	for _, symbol := range symbols {
		fmt.Println("Starting pair asset: ", symbol)
		e.PairAssets[symbol] = strategy.NewPairAsset(symbol)
		e.PairAssets[symbol].SetChannel(e.ChMap[symbol])

		go e.PairAssets[symbol].Run(e.done)
	}
}

func (e *Engine) StopPairAssets() {
	symbols := e.FutureConfig.GetAvailableSymbols()
	for _, symbol := range symbols {
		e.PairAssets[symbol].Close()
	}
}

// StartPairs starts the Pairs for the available symbols.
// Subscribes to the PairAssets and receives data from them.
func (e *Engine) StartPairs() {
	pairAssets := e.FutureConfig.CreatePair()

	for _, pair := range pairAssets {
		symbols := strings.Split(pair, "-")

		asset1, ok := e.PairAssets[symbols[0]]
		if !ok {
			log.Println("Asset1 not found: ", symbols[0])
			continue
		}
		asset2, ok := e.PairAssets[symbols[1]]
		if !ok {
			log.Println("Asset2 not found: ", symbols[1])
			continue
		}

		e.Pairs[pair] = strategy.NewPair()
		e.Pairs[pair].Subscribe(asset1, asset2)

		go e.Pairs[pair].Run(e.done)
	}
}

func (e *Engine) StopPairs() {

}

// StartStream starts the stream for the available symbols.
// It subscribes to the stream channels and receives data from them.
// Should be called after StartStreamCh and StartPairAssets.
func (e *Engine) StartStream() {
	symbols := e.FutureConfig.GetAvailableSymbols()
	// Group symbols into chunks of 5
	const chunkSize = 5
	for i := 0; i < len(symbols); i += chunkSize {
		end := min(i+chunkSize, len(symbols))
		symbolGroup := symbols[i:end]

		// Process each group of symbols - Multiple queires
		go streams.SubscribeKlineMulti(symbolGroup, "1m", e.ChMap, e.done)
	}

	e.isStreaming = true
}

func (e *Engine) StopStream() {
}

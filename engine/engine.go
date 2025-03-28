package engine

import (
	"fmt"
	"log"
	"strings"
	"time"

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
	PairAssets map[string]*strategy.Asset
	Pairs      map[string]*strategy.Pair
	ChMap      map[string]chan internal.KlineDataStream // Channels for kline data. Key is symbol

	// Engine status
	done            chan struct{}
	numberOfAssets  int
	numberOfStreams int
	numberOfPairs   int
	isTest          bool
	isStreaming     bool
	isCalculating   bool
}

func NewEngine() *Engine {
	db, err := data.ConnectDB()
	if err != nil {
		panic("[engine] Failed to connect to database: " + err.Error())
	}
	log.Println("[engine] Connected to database")

	exchangeRestApi := data.NewBinanceFutureMarketData()
	log.Println("[engine] Connected to binance-futures data source")

	// Future exchange information
	config, err := config.NewFutureTradeConfig()
	if err != nil {
		panic("[engine] Failed to connect to binance-futures exchange: " + err.Error())
	}
	log.Println("[engine] Connected to binance-futures exchange")

	// Setup the internal Engine's attributes
	exchangeRestApi.UpdateRateLimit(config.ExchangeInfo.GetRequestRateLimit())

	engine := &Engine{
		FutureConfig:    config,
		Database:        db,
		FutureMarket:    exchangeRestApi,
		PairAssets:      make(map[string]*strategy.Asset),
		Pairs:           make(map[string]*strategy.Pair),
		ChMap:           make(map[string]chan internal.KlineDataStream),
		done:            make(chan struct{}),
		isTest:          false,
		isStreaming:     false,
		isCalculating:   false,
		numberOfAssets:  0,
		numberOfStreams: 0,
		numberOfPairs:   0,
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
		e.ChMap[symbol] = make(chan internal.KlineDataStream, 5)
	}
}

// StartAssets starts the PairAssets for the available symbols.
// Receives data from the stream channels and transfers and broadcasts to the Pairs.
func (e *Engine) StartAssets() {
	symbols := e.FutureConfig.GetAvailableSymbols()
	for _, symbol := range symbols {
		fmt.Println("[engine] Starting individual asset: ", symbol)
		e.PairAssets[symbol] = strategy.NewAsset(symbol)
		e.PairAssets[symbol].SetChannel(e.ChMap[symbol])

		go e.PairAssets[symbol].Run(e.done)
		e.numberOfAssets++
	}
}

func (e *Engine) StopAssets() {
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

		// Get and validate both assets
		assetY, assetX, err := e.getAndValidateAssets(symbols)
		if err != nil {
			log.Println("[engine] Error getting assets:", err)
			continue
		}

		// Update historic data for both assets if needed
		if err := e.updateHistoricData(assetY); err != nil {
			log.Println("[engine] Error updating historic data for asset Y:", err)
			continue
		}

		if err := e.updateHistoricData(assetX); err != nil {
			log.Println("[engine] Error updating historic data for asset X:", err)
			continue
		}

		// Initialize and start the pair
		log.Println("[engine] Starting pair:", pair)
		if err := e.initializeAndStartPair(pair, assetY, assetX); err != nil {
			log.Println("[engine] Error initializing pair:", err)
			continue
		}
		e.numberOfPairs++
	}
}

func (e *Engine) getAndValidateAssets(symbols []string) (*strategy.Asset, *strategy.Asset, error) {
	assetY, ok := e.PairAssets[symbols[0]]
	if !ok {
		return nil, nil, fmt.Errorf("asset1 not found: %s", symbols[0])
	}

	assetX, ok := e.PairAssets[symbols[1]]
	if !ok {
		return nil, nil, fmt.Errorf("asset2 not found: %s", symbols[1])
	}

	return assetY, assetX, nil
}

func (e *Engine) updateHistoricData(asset *strategy.Asset) error {
	if asset.HistoricKlines == nil {
		return e.fetchAndSetHistoricData(asset)
	}

	timestampNow := time.Now().Unix() * 1000
	closeTime, err := asset.HistoricKlines.GetKlineLatestCloseTime()
	if err != nil {
		return fmt.Errorf("failed to get kline closing time data: %v", err)
	}

	if timestampNow-int64(closeTime) > 60000 { // 1 minute
		return e.fetchAndSetHistoricData(asset)
	}

	return nil
}

func (e *Engine) fetchAndSetHistoricData(asset *strategy.Asset) error {
	klineData, err := e.FutureMarket.GetKlineData(asset.Symbol, "1m", 100)
	if err != nil {
		return fmt.Errorf("failed to get kline data: %v", err)
	}
	asset.SetHistoricPrice(klineData)
	return nil
}

func (e *Engine) initializeAndStartPair(pair string, assetY, assetX *strategy.Asset) error {
	e.Pairs[pair] = strategy.NewPair()
	e.Pairs[pair].Subscribe(assetY, assetX)

	assetYClosePrices, err := assetY.HistoricKlines.GetKlineClosePrices()
	if err != nil {
		return fmt.Errorf("failed to get kline close prices for asset Y: %v", err)
	}

	assetXClosePrices, err := assetX.HistoricKlines.GetKlineClosePrices()
	if err != nil {
		return fmt.Errorf("failed to get kline close prices for asset X: %v", err)
	}

	e.Pairs[pair].WarmUpFilter(assetYClosePrices, assetXClosePrices)
	log.Println("[engine] Warm up filter:", pair)
	go e.Pairs[pair].Run(e.done)
	return nil
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
		e.numberOfStreams++
	}

	e.isStreaming = true
}

func (e *Engine) StopStream() {}

func (e *Engine) GetStatus() string {
	return fmt.Sprintf("Assets: %d, Streams: %d, Pairs: %d", e.numberOfAssets, e.numberOfStreams, e.numberOfPairs)
}

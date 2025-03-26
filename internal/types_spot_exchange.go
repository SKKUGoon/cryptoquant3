package internal

import (
	"encoding/json"
	"io"
	"net/http"
)

// Exchange: Receives exchange information from the API request
type Exchange struct {
	Timezone   string `json:"timezone"`
	ServerTime int64  `json:"serverTime"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	ExchangeFilters []any                  `json:"exchangeFilters"`
	Symbols         []SymbolInfo           `json:"symbols"`
	symbolsMap      map[string]*SymbolInfo `json:"-"` // -: ignore this field during JSON marshaling
}

type SymbolInfo struct {
	Symbol                          string         `json:"symbol"`
	Status                          string         `json:"status"`
	BaseAsset                       string         `json:"baseAsset"`
	BaseAssetPrecision              int            `json:"baseAssetPrecision"`
	QuoteAsset                      string         `json:"quoteAsset"`
	QuotePrecision                  int            `json:"quotePrecision"`      // Deprecated
	QuoteAssetPrecision             int            `json:"quoteAssetPrecision"` // Notify k decimal places
	BaseCommissionPrecision         int            `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision        int            `json:"quoteCommissionPrecision"`
	OrderTypes                      []string       `json:"orderTypes"`
	IcebergAllowed                  bool           `json:"icebergAllowed"`
	OcoAllowed                      bool           `json:"ocoAllowed"`
	OtoAllowed                      bool           `json:"otoAllowed"`
	QuoteOrderQtyMarketAllowed      bool           `json:"quoteOrderQtyMarketAllowed"`
	AllowTrailingStop               bool           `json:"allowTrailingStop"`
	CancelReplaceAllowed            bool           `json:"cancelReplaceAllowed"`
	IsSpotTradingAllowed            bool           `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed          bool           `json:"isMarginTradingAllowed"`
	Filters                         []SymbolFilter `json:"filters"`
	Permissions                     []any          `json:"permissions"`
	PermissionSets                  [][]string     `json:"permissionSets"`
	DefaultSelfTradePreventionMode  string         `json:"defaultSelfTradePreventionMode"`
	AllowedSelfTradePreventionModes []string       `json:"allowedSelfTradePreventionModes"`
}

type SymbolFilter struct {
	FilterType            string `json:"filterType"`
	MinPrice              string `json:"minPrice,omitempty"`
	MaxPrice              string `json:"maxPrice,omitempty"`
	TickSize              string `json:"tickSize,omitempty"`
	MinQty                string `json:"minQty,omitempty"`
	MaxQty                string `json:"maxQty,omitempty"`
	StepSize              string `json:"stepSize,omitempty"`
	Limit                 int    `json:"limit,omitempty"`
	MinTrailingAboveDelta int    `json:"minTrailingAboveDelta,omitempty"`
	MaxTrailingAboveDelta int    `json:"maxTrailingAboveDelta,omitempty"`
	MinTrailingBelowDelta int    `json:"minTrailingBelowDelta,omitempty"`
	MaxTrailingBelowDelta int    `json:"maxTrailingBelowDelta,omitempty"`
	BidMultiplierUp       string `json:"bidMultiplierUp,omitempty"`
	BidMultiplierDown     string `json:"bidMultiplierDown,omitempty"`
	AskMultiplierUp       string `json:"askMultiplierUp,omitempty"`
	AskMultiplierDown     string `json:"askMultiplierDown,omitempty"`
	AvgPriceMins          int    `json:"avgPriceMins,omitempty"`
	MinNotional           string `json:"minNotional,omitempty"`
	ApplyMinToMarket      bool   `json:"applyMinToMarket,omitempty"`
	MaxNotional           string `json:"maxNotional,omitempty"`
	ApplyMaxToMarket      bool   `json:"applyMaxToMarket,omitempty"`
	MaxNumOrders          int    `json:"maxNumOrders,omitempty"`
	MaxNumAlgoOrders      int    `json:"maxNumAlgoOrders,omitempty"`
}

// NewSpotExchange creates a new exchange info using the API response
func NewSpotExchange() (*Exchange, error) {
	var exchangeInfo Exchange

	// Binance API url + endpoint
	url := "https://api.binance.com/api/v3/exchangeInfo"

	// Make the API request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &exchangeInfo)
	if err != nil {
		return nil, err
	}

	// Initialize the map and populate it
	exchangeInfo.symbolsMap = make(map[string]*SymbolInfo)
	for i := range exchangeInfo.Symbols {
		exchangeInfo.symbolsMap[exchangeInfo.Symbols[i].Symbol] = &exchangeInfo.Symbols[i]
	}

	return &exchangeInfo, nil
}

// GetSymbolInfo returns the symbol information for the given symbol
func (e *Exchange) GetSymbolInfo(symbol string) *SymbolInfo {
	return e.symbolsMap[symbol]
}

// GetSymbolFilter returns the filter of specified type for the given symbol
func (s *SymbolInfo) GetSymbolFilter(filterType string) *SymbolFilter {
	for _, f := range s.Filters {
		if f.FilterType == filterType {
			return &f
		}
	}
	return nil
}

// GetSymbolQuotePrecision returns the quote precision for the given symbol
func (s *SymbolInfo) GetSymbolQuotePrecision() int {
	return s.QuotePrecision
}

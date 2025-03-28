package internal

import (
	"fmt"
	"strconv"
)

type FutureKlineData [][]any

// Index of the data in the FutureKlineData
const (
	FutureKlineDataOpenTime                 = 0
	FutureKlineDataOpenPrice                = 1
	FutureKlineDataHighPrice                = 2
	FutureKlineDataLowPrice                 = 3
	FutureKlineDataClosePrice               = 4
	FutureKlineDataVolume                   = 5
	FutureKlineDataCloseTime                = 6
	FutureKlineDataQuoteVolume              = 7
	FutureKlineDataNumTrades                = 8
	FutureKlineDataTakerBuyBaseAssetVolume  = 9
	FutureKlineDataTakerBuyQuoteAssetVolume = 10
	FutureKlineDataIgnore                   = 11
)

func (k *FutureKlineData) GetKlineClosePrices() ([]float64, error) {
	if k == nil || len(*k) == 0 {
		return nil, fmt.Errorf("kline data is nil or empty")
	}

	closePrices := make([]float64, len(*k))
	for i, price := range *k {
		if len(price) <= FutureKlineDataClosePrice {
			return nil, fmt.Errorf("kline data at index %d is missing close price", i)
		}
		closePrice, err := strconv.ParseFloat(price[FutureKlineDataClosePrice].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("kline data at index %d is not a float64", i)
		}
		closePrices[i] = closePrice
	}
	return closePrices, nil
}

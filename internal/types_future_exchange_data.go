package internal

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

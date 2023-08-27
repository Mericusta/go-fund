package searcher

import "time"

type StockDailyData interface {
	Code() string
	Market() string
	Date() time.Time
	Open() float64
	Close() float64
	High() float64
	Low() float64
	Volume() float64
	Amount() float64
}

// func

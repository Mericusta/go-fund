package model

import "time"

type StockBriefData interface {
	Code() string
	Name() string
}

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

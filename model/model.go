package model

import "time"

type StockBriefData interface {
	Name() string
	Code() string
}

type StockDailyData interface {
	StockBriefData
	Market() string
	Date() time.Time
	Open() float64
	Close() float64
	High() float64
	Low() float64
	Volume() float64
	Amount() float64
	ChangePercent() string
}

package legulegu

import (
	"encoding/json"
	"os"
)

type StockData struct {
	date    int     `json:"date"`
	open    float32 `json:"open"`
	close   float32 `json:"close"`
	high    float32 `json:"high"`
	low     float32 `json:"low"`
	volume  float32 `json:"volume"`
	tVolume float32 `json:"tVolume"`
}

func (sd *StockData) Open() float32 {
	return sd.open
}

func (sd *StockData) Close() float32 {
	return sd.close
}

type MockData struct {
	MarketID     string       `json:"marketId"`
	Name         string       `json:"name"`
	MockDataList []*StockData `json:"mockDataList"`
}

func (md *MockData) Parse(p string) {
	c, e := os.ReadFile(p)
	if e != nil {
		panic(e)
	}

	e = json.Unmarshal(c, md)
	if e != nil {
		panic(e)
	}
}

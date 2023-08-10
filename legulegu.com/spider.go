package legulegu

import (
	"encoding/json"
	"go-fund/global"
	"os"
	"path/filepath"
	"time"
)

var (
	resourceRelativePath string = "markdown/note/stock"
)

type StockData struct {
	Date    int     `json:"date"`
	Open    float32 `json:"open"`
	Close   float32 `json:"close"`
	High    float32 `json:"high"`
	Low     float32 `json:"low"`
	Volume  float32 `json:"volume"`
	TVolume float32 `json:"tVolume"`
}

func (sd *StockData) DateTime() time.Time {
	return time.Unix(int64(sd.Date/1000), 0)
}

func (sd *StockData) OpenValue() float32 {
	return sd.Open
}

func (sd *StockData) CloseValue() float32 {
	return sd.Close
}

type MockData struct {
	MarketID      string       `json:"marketId"`
	Name          string       `json:"name"`
	MockDataSlice []*StockData `json:"mockDataList"`
}

func NewMockData() *MockData {
	return &MockData{}
}

func (md *MockData) Parse(p string) {
	c, e := os.ReadFile(filepath.Join(global.PersonalDocumentPath, resourceRelativePath, p))
	if e != nil {
		panic(e)
	}

	e = json.Unmarshal(c, md)
	if e != nil {
		panic(e)
	}
}

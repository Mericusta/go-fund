package main

import (
	"go-fund/legulegu.com"
	"go-fund/searcher"

	"github.com/Mericusta/go-stp"
)

func main() {
	md := &legulegu.MockData{}
	md.Parse("./resource/202308082150.json")

	filterBeginDateTS, filterEndDateTS := 1554220800, 1567958400
	mockDataArray := stp.NewArray(md.MockDataSlice)
	filterStockData := mockDataArray.Filter(func(v *legulegu.StockData, i int) bool {
		if ts := v.Date / 1000; filterBeginDateTS <= ts && ts <= filterEndDateTS {
			return true
		}
		return false
	})

	stockData := make([]searcher.SearchMethod1Stock, 0, len(md.MockDataSlice))
	for _, d := range md.MockDataSlice {
		stockData = append(stockData, d)
	}
	searcher.SearchMethod1(stockData, 2)
}

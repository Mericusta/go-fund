package main

import (
	"go-fund/legulegu.com"
	"go-fund/searcher"

	"github.com/Mericusta/go-stp"
)

func main() {
	md := legulegu.NewMockData()
	md.Parse("./resource/202308082150.json")

	filterBeginDateTS, filterEndDateTS := 1554220800, 1567958400
	stockData := make([]searcher.SearchMethod1Stock, 0, len(md.MockDataSlice))
	mockDataArray := stp.NewArray(md.MockDataSlice)
	mockDataArray.ForEach(func(v *legulegu.StockData, i int) {
		if ts := v.Date / 1000; filterBeginDateTS <= ts && ts <= filterEndDateTS {
			stockData = append(stockData, v)
			return
		}
	})

	searcher.SearchMethod1(stockData, 2)
}

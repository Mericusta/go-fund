package main

import (
	"go-fund/legulegu.com"
	"go-fund/searcher"

	"github.com/Mericusta/go-stp"
)

func main() {
	// stockCodeMap := appfinanceifengcom.GetStockList()
	// for code, name := range stockCodeMap {
	// 	fmt.Printf("name %v, code %v\n", name, code)
	// }
	// tushare.GetDailyData("601688.SH", 0, time.Now().Add(-7*time.Hour*24).Unix(), time.Now().Unix())
}

func search_legulegu() {
	md := legulegu.NewMockData()
	md.Parse("202308080005.json")

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

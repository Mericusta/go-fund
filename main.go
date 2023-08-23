package main

import (
	appfinanceifengcom "go-fund/app.finance.ifeng.com"
	"go-fund/legulegu.com"
	"go-fund/searcher"
	"go-fund/tushare.pro"
	"time"

	"github.com/Mericusta/go-stp"
)

func main() {
	stockCodeMap := appfinanceifengcom.GetStockList()
	appfinanceifengcom.SaveStockList(stockCodeMap)
	appfinanceifengcom.LoadStockList()

	dailyData := tushare.GetDailyData("601688.SH", 0, time.Now().AddDate(-1, 0, 0).Unix(), time.Now().Unix())
	tushare.SaveStockDailyData("601688.SH", dailyData)
	tushare.LoadStockDailyData("601688.SH")
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

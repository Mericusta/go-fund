package main

import (
	"fmt"
	appfinanceifengcom "go-fund/app.finance.ifeng.com"
	"go-fund/filter"
	"go-fund/legulegu.com"
	"go-fund/searcher"
	"go-fund/tushare.pro"
	"time"

	"github.com/Mericusta/go-stp"
)

func main() {
	// ---------------- stock list ----------------
	// fmt.Printf("spider stock list\n")
	// stockCodeMap := appfinanceifengcom.GetStockList()
	// fmt.Printf("save stock list\n")
	// appfinanceifengcom.SaveStockList(stockCodeMap)

	fmt.Printf("load stock list\n")
	stockCodeNameMap := appfinanceifengcom.LoadStockList()

	// ---------------- daily data ----------------
	f := func(stockCodeNameMap map[string]string, suffix string) {
		for code, name := range stockCodeNameMap {
			_code := fmt.Sprintf("%v.%v", code, suffix)
			fmt.Printf("spider stock %v - %v daily data", _code, name)
			dailyData := tushare.GetDailyData(_code, 0, time.Now().AddDate(-1, 0, 0).Unix(), time.Now().Unix())
			fmt.Printf("save stock %v - %v daily data", _code, name)
			tushare.SaveStockDailyData(_code, dailyData)
			// fmt.Printf("load stock %v - %v daily data", _code, name)
			// tushare.LoadStockDailyData(_code)
			time.Sleep(time.Second)
		}
	}

	SHStockCodeNameMap, suffix := filter.SH_StockFilter(stockCodeNameMap)
	f(SHStockCodeNameMap, suffix)
	SZStockCodeNameMap, suffix := filter.SZ_StockFilter(stockCodeNameMap)
	f(SZStockCodeNameMap, suffix)
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

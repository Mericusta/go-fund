package spider

import (
	"fmt"
	appfinanceifengcom "go-fund/app.finance.ifeng.com"
	"go-fund/filter"
	"go-fund/searcher"
	"go-fund/tushare.pro"
	"sync"
	"time"
)

func DownloadStockList() {
	stockCodeMap := appfinanceifengcom.GetStockList()
	appfinanceifengcom.SaveStockList(stockCodeMap)
}

func DownloadStockDailyData(beginDate, endDate time.Time) {
	stockCodeNameMap := appfinanceifengcom.LoadStockList()
	SHStockCodeNameMap, suffix := filter.SH_StockFilter(stockCodeNameMap)
	downloadStockDailyData(SHStockCodeNameMap, suffix, beginDate, endDate)
	SZStockCodeNameMap, suffix := filter.SZ_StockFilter(stockCodeNameMap)
	downloadStockDailyData(SZStockCodeNameMap, suffix, beginDate, endDate)
}

func downloadStockDailyData(stockCodeNameMap map[string]string, suffix string, beginDate, endDate time.Time) {
	fmt.Printf("- spider download stock daily data\n")
	fmt.Printf("\t- market: %v\n", suffix)
	fmt.Printf("\t- duration: %v ~ %v\n", beginDate.Format("20060102"), endDate.Format("20060102"))
	counter := 0
	wg := sync.WaitGroup{}
	wg.Add(len(stockCodeNameMap))
	for code, name := range stockCodeNameMap {
		_code := fmt.Sprintf("%v.%v", code, suffix)
		go func(c, n string) {
			defer func() {
				if _err := recover(); _err != nil {
					fmt.Printf("\t\t- spider stock %v - %v occurs error: %v\n", c, n, _err)
				}
			}()
			dailyData := tushare.GetDailyData(_code, name, 0, beginDate.Unix(), endDate.Unix())
			tushare.SaveStockDailyData(c, n, dailyData)
			wg.Done()
		}(_code, name)
		counter++
		if counter%5 == 0 {
			time.Sleep(time.Second)
		}
	}
	wg.Wait()
	fmt.Printf("- spider download stock daily data done, count %v\n", len(stockCodeNameMap))
}

func LoadStockDailyData(code string, beginDate, endDate time.Time) []searcher.StockDailyData {
	slice := make([]searcher.StockDailyData, 0, 128)
	for _, data := range tushare.LoadStockDailyData(code) {
		tradeDate, err := time.Parse("20060102", data.TS_TradeDate)
		if err != nil {
			panic(err)
		}
		if tradeDate.After(beginDate) && tradeDate.Before(endDate) {
			slice = append(slice, data)
		}
	}
	return slice
}

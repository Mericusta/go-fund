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
	SHStockCodeNameMap, market := filter.SH_StockFilter(stockCodeNameMap)
	downloadStockDailyData(SHStockCodeNameMap, market, beginDate, endDate)
	SZStockCodeNameMap, market := filter.SZ_StockFilter(stockCodeNameMap)
	downloadStockDailyData(SZStockCodeNameMap, market, beginDate, endDate)
}

func downloadStockDailyData(stockCodeNameMap map[string]string, market string, beginDate, endDate time.Time) {
	fmt.Printf("- spider download stock daily data\n")
	fmt.Printf("\t- market: %v\n", market)
	fmt.Printf("\t- duration: %v ~ %v\n", beginDate.Format(tushare.TradeDateLayout), endDate.Format(tushare.TradeDateLayout))
	counter := 0
	wg := sync.WaitGroup{}
	wg.Add(len(stockCodeNameMap))
	for code, name := range stockCodeNameMap {
		_code := fmt.Sprintf("%v.%v", code, market)
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

func AppendStockDailyData() {
	stockCodeNameMap := appfinanceifengcom.LoadStockList()
	SHStockCodeNameMap, market := filter.SH_StockFilter(stockCodeNameMap)
	appendStockDailyData(SHStockCodeNameMap, market)
	SZStockCodeNameMap, market := filter.SZ_StockFilter(stockCodeNameMap)
	appendStockDailyData(SZStockCodeNameMap, market)

}

func appendStockDailyData(stockCodeNameMap map[string]string, market string) {
	fmt.Printf("- spider append stock daily data\n")
	fmt.Printf("\t- market: %v\n", market)
	counter := 0
	wg := sync.WaitGroup{}
	wg.Add(len(stockCodeNameMap))
	for code, name := range stockCodeNameMap {
		_code := fmt.Sprintf("%v.%v", code, market)
		_name := name
		go func(c, n string) {
			defer func() {
				if _err := recover(); _err != nil {
					fmt.Printf("\t\t- spider stock %v - %v occurs error: %v\n", c, n, _err)
				}
			}()

			var beginDate time.Time
			endDate, err := time.Parse(tushare.TradeDateLayout, time.Now().AddDate(0, 0, 1).Format(tushare.TradeDateLayout))
			if err != nil {
				panic(err)
			}
			_dailyData := tushare.LoadStockDailyData(_code)
			if len(_dailyData) == 0 {
				beginDate = endDate.AddDate(0, 0, -1)
			} else {
				beginDate, err = time.Parse(tushare.TradeDateLayout, _dailyData[0].TS_TradeDate)
				if err != nil {
					panic(err)
				}
				beginDate = beginDate.AddDate(0, 0, 1)
			}

			dailyData := tushare.GetDailyData(_code, _name, 0, beginDate.Unix(), endDate.Unix())
			dailyData = append(dailyData, _dailyData...)
			tushare.SaveStockDailyData(c, n, dailyData)
			wg.Done()
		}(_code, _name)
		counter++
		if counter%5 == 0 {
			time.Sleep(time.Second)
		}
	}
	wg.Wait()
	fmt.Printf("- spider append stock daily data done, count %v\n", len(stockCodeNameMap))
}

func LoadStockDailyData(code string, beginDate, endDate time.Time) []searcher.StockDailyData {
	slice := make([]searcher.StockDailyData, 0, 128)
	for _, data := range tushare.LoadStockDailyData(code) {
		tradeDate, err := time.Parse(tushare.TradeDateLayout, data.TS_TradeDate)
		if err != nil {
			panic(err)
		}
		if tradeDate.After(beginDate) && tradeDate.Before(endDate) {
			slice = append(slice, data)
		}
	}
	return slice
}

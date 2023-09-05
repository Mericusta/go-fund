package searcher

import (
	"fmt"
	"go-fund/filter"
	appfinanceifengcom "go-fund/spider/app.finance.ifeng.com"
	"go-fund/spider/tushare.pro"
	"sort"
	"sync"
	"time"

	"github.com/Mericusta/go-stp"
)

// SearchAlgorithm1
func SearchAlgorithm1() {
	fmt.Printf("- searcher algorithm 1\n")
	const (
		level int = 5
	)
	var (
		now        string          = time.Now().Format(tushare.TradeDateLayout())
		searchWG   *sync.WaitGroup = &sync.WaitGroup{}
		reportWG   *sync.WaitGroup = &sync.WaitGroup{}
		resultChan chan string
	)

	stockBriefDataList := appfinanceifengcom.LoadStockBriefList()

	reportWG.Add(1)
	resultChan = make(chan string, len(stockBriefDataList))
	go func(_wg *sync.WaitGroup, _resultChan chan string) {
		defer _wg.Done()
		fmt.Printf("- search result:\n")
		for result := range _resultChan {
			fmt.Printf("\t- code %v\n", result)
		}
	}(reportWG, resultChan)

	searchWG.Add(len(stockBriefDataList))
	for _, stockBriefData := range stockBriefDataList {
		code := stockBriefData.Code()
		switch {
		case filter.SH_StockFilter(code):
			code = tushare.MakeStockTSCode(code, filter.SH_Market)
		case filter.SZ_StockFilter(code):
			code = tushare.MakeStockTSCode(code, filter.SZ_Market)
		default:
			searchWG.Done()
			continue
		}
		go func(_wg *sync.WaitGroup, _level int, _now, _code string, _resultChan chan string) {
			defer _wg.Done()
			stockDailyData := tushare.LoadStockDailyData(_code)
			_nowData := stp.NewArray(stockDailyData).Find(func(v *tushare.TS_StockDailyData, i int) bool {
				return v.TS_TradeDate == _now
			})
			if _nowData == nil {
				return
			}
			levelLowestValueSlice := searchAlgorithmUtility1(stockDailyData, _level)
			index := stp.NewArray(levelLowestValueSlice).FindIndex(func(v float64, i int) bool {
				return _nowData.TS_Low < v
			})
			if index == -1 {
				return
			}
			_resultChan <- _code
		}(searchWG, level, now, code, resultChan)
	}
	searchWG.Wait()

	close(resultChan)

	reportWG.Wait()

	fmt.Printf("- searcher algorithm 1 done\n")
}

// searchAlgorithmUtility1 lowest value from daily data
func searchAlgorithmUtility1(slice []*tushare.TS_StockDailyData, level int) []float64 {
	valueSlice := make([]float64, 0, len(slice))
	for _, sdd := range slice {
		valueSlice = append(valueSlice, sdd.Low())
	}
	sort.Float64s(valueSlice)
	if level > len(valueSlice) {
		return valueSlice
	}
	return valueSlice[:level]
}

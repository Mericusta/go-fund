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
// 搜索 指定时间点的价格 是否在 指定时间点过去的某一段时间内的最低价格构成的有限有序集合 中
// @param1          指定最低价格构成的有序集合的元素的数量
// @param2          指定时间点与当前时间点的相差日
// @param3          指定时间点过去若干年
// @param4          指定时间点过去若干月
// @param5          指定时间点过去若干日
func SearchAlgorithm1(count, nowOffsetDay, year, month, day int) {
	fmt.Printf("- searcher algorithm 1\n")
	var (
		end        time.Time       = time.Now().AddDate(0, 0, nowOffsetDay)
		begin      time.Time       = end.AddDate(year, month, day)
		searchWG   *sync.WaitGroup = &sync.WaitGroup{}
		reportWG   *sync.WaitGroup = &sync.WaitGroup{}
		resultChan chan string
	)
	diffDay, weekday := 0, time.Now().Weekday()
	switch weekday {
	case 0:
		diffDay = -2
	case 6:
		diffDay = -1
	}
	end = end.AddDate(0, 0, diffDay)
	begin = begin.AddDate(0, 0, diffDay)
	fmt.Printf("\t- statistics count %v\n", count)
	fmt.Printf("\t- specify trade date %v ~ %v\n", begin.Format(tushare.TradeDateLayout()), end.Format(tushare.TradeDateLayout()))

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
		go func(_wg *sync.WaitGroup, _count int, _begin, _end time.Time, _code string, _resultChan chan string) {
			defer _wg.Done()
			stockDailyData := tushare.LoadStockDailyData(_code)
			_nowData := stp.NewArray(stockDailyData).Find(func(v *tushare.TS_StockDailyData, i int) bool {
				return v.TS_TradeDate == _end.Format(tushare.TradeDateLayout())
			})
			if _nowData == nil {
				return
			}
			stockDailyData = stp.NewArray(stockDailyData).Filter(func(v *tushare.TS_StockDailyData, i int) bool {
				t, err := time.Parse(tushare.TradeDateLayout(), v.TS_TradeDate)
				if err != nil {
					panic(err)
				}
				return (t.Equal(_begin) || t.After(_begin)) && (t.Before(_end) || t.Equal(_end))
			}).Slice()
			levelLowestValueSlice := searchAlgorithmUtility1(stockDailyData, _count)
			index := stp.NewArray(levelLowestValueSlice).FindIndex(func(v float64, i int) bool {
				return _nowData.TS_Low < v
			})
			if index == -1 {
				return
			}
			_resultChan <- _code
		}(searchWG, count, begin, end, code, resultChan)
	}
	searchWG.Wait()

	close(resultChan)

	reportWG.Wait()

	fmt.Printf("- searcher algorithm 1 done\n")
}

// searchAlgorithmUtility1 lowest value from daily data make up a sorted slice
func searchAlgorithmUtility1(slice []*tushare.TS_StockDailyData, count int) []float64 {
	valueSlice := make([]float64, 0, len(slice))
	for _, sdd := range slice {
		valueSlice = append(valueSlice, sdd.Low())
	}
	sort.Float64s(valueSlice)
	if count > len(valueSlice) {
		return valueSlice
	}
	return valueSlice[:count]
}

package spider

import (
	"fmt"
	"go-fund/filter"
	"go-fund/model"
	appfinanceifengcom "go-fund/spider/app.finance.ifeng.com"
	fundeastmoney "go-fund/spider/fund.eastmoney.com"
	"go-fund/spider/tushare.pro"
	"sync"
	"time"

	"github.com/Mericusta/go-stp"
)

// MultiDownloadStockDailyData 批量下载股票每日行情数据
// - downloadDates 未缺省时，下载 downloadDates 范围内的
// - downloadDates 缺省时，取当日时间
// - 当日时间15点后下载当日的行情数据
// - 当日时间15点前下载昨日的行情数据
func MultiDownloadStockDailyData(stockBriefDataSlice []model.StockBriefData, downloadDates ...time.Time) {
	var beginDate, endDate time.Time
	// make up download dates
	l := len(downloadDates)
	switch {
	case l == 0 || downloadDates[0].Format(tushare.TradeDateLayout()) == time.Now().Format(tushare.TradeDateLayout()):
		downloadDate := time.Now()
		if downloadDate.Hour() < 15 {
			downloadDate = downloadDate.AddDate(0, 0, -1)
		}
		beginDate, endDate = downloadDate, downloadDate
		fmt.Println("- spider download stock daily data on", downloadDate.Format(tushare.TradeDateLayout()))
	case l == 1:
		beginDate, endDate = downloadDates[0], downloadDates[0]
		fmt.Println("- spider download stock daily data on", downloadDates[0].Format(tushare.TradeDateLayout()))
	case l > 1:
		for _, downloadDate := range downloadDates {
			if beginDate.IsZero() || downloadDate.Before(beginDate) {
				beginDate = downloadDate
			}
			if endDate.IsZero() || downloadDate.After(endDate) {
				endDate = downloadDate
			}
		}
		fmt.Println("- spider download stock daily data during", beginDate.Format(tushare.TradeDateLayout()), "to", endDate.Format(tushare.TradeDateLayout()))
	}
	// make up stock brief data
	SHStockBriefDataSlice := make([]model.StockBriefData, 0, len(stockBriefDataSlice)/2)
	SZStockBriefDataSlice := make([]model.StockBriefData, 0, len(stockBriefDataSlice)/2)
	stp.NewArray(stockBriefDataSlice).ForEach(func(v model.StockBriefData, i int) {
		switch {
		case filter.SH_StockFilter(v.Code()) || filter.SH_ETF_FundFilter(v.Code()):
			SHStockBriefDataSlice = append(SHStockBriefDataSlice, v)
		case filter.SZ_StockFilter(v.Code()) || filter.SZ_ETF_FundFilter(v.Code()):
			SZStockBriefDataSlice = append(SZStockBriefDataSlice, v)
		}
	})
	// download
	if len(SHStockBriefDataSlice) > 0 {
		downloadStockDailyData(SHStockBriefDataSlice, filter.SH_Market, beginDate, endDate)
	}
	if len(SZStockBriefDataSlice) > 0 {
		downloadStockDailyData(SZStockBriefDataSlice, filter.SZ_Market, beginDate, endDate)
	}
}

// downloadStockDailyData 下载股票日行情数据
func downloadStockDailyData(stockBriefDataSlice []model.StockBriefData, market string, beginDate, endDate time.Time) {
	fmt.Println("\t- spider download stock daily data in market", market, "total", len(stockBriefDataSlice))
	counter, wg := 0, sync.WaitGroup{}
	wg.Add(len(stockBriefDataSlice))
	for _, stockBriefData := range stockBriefDataSlice {
		name, code := stockBriefData.Name(), fmt.Sprintf("%v.%v", stockBriefData.Code(), market)
		go func(_code, _name string) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("\t\t- spider stock %v - %v occurs error: %v\n", code, name, err)
				}
			}()
			dailyData := tushare.DownloadDailyData(_code, _name, 0, beginDate.Unix(), endDate.Unix())
			if len(dailyData) > 0 {
				tushare.SaveStockDailyData(_code, _name, dailyData)
			}
			wg.Done()
		}(code, name)
		counter++
		if counter%5 == 0 {
			time.Sleep(time.Second)
		}
	}
	wg.Wait()
	fmt.Println("\t- spider download stock daily data in market", market, "done, count", len(stockBriefDataSlice))
}

// MultiLoadStockDailyData 批量加载股票日行情数据
func MultiLoadStockDailyData(stockBriefDataSlice []model.StockBriefData) [][]model.StockDailyData {
	// make up stock brief data
	SHStockBriefDataSlice := make([]model.StockBriefData, 0, len(stockBriefDataSlice)/2)
	SZStockBriefDataSlice := make([]model.StockBriefData, 0, len(stockBriefDataSlice)/2)
	stp.NewArray(stockBriefDataSlice).ForEach(func(v model.StockBriefData, i int) {
		switch {
		case filter.SH_StockFilter(v.Code()):
			SHStockBriefDataSlice = append(SHStockBriefDataSlice, v)
		case filter.SZ_StockFilter(v.Code()):
			SZStockBriefDataSlice = append(SZStockBriefDataSlice, v)
		}
	})
	// load
	stockDailyDataSlice := make([][]model.StockDailyData, 0, 32)
	if len(SHStockBriefDataSlice) > 0 {
		stockDailyDataSlice = append(stockDailyDataSlice, loadStockDailyData(SHStockBriefDataSlice, filter.SH_Market)...)
	}
	if len(SZStockBriefDataSlice) > 0 {
		stockDailyDataSlice = append(stockDailyDataSlice, loadStockDailyData(SZStockBriefDataSlice, filter.SZ_Market)...)
	}
	return stockDailyDataSlice
}

type spiderData struct {
	model.StockBriefData
	*tushare.TS_StockDailyData
}

func (sd *spiderData) Name() string { return sd.StockBriefData.Name() }
func (sd *spiderData) Code() string { return sd.StockBriefData.Code() }

// loadStockDailyData 本地加载股票每日行情数据
func loadStockDailyData(stockBriefDataSlice []model.StockBriefData, market string) [][]model.StockDailyData {
	counter, wg, stockDailyDataSlice := 0, sync.WaitGroup{}, make([][]model.StockDailyData, len(stockBriefDataSlice))
	wg.Add(len(stockBriefDataSlice))
	for index, stockBriefData := range stockBriefDataSlice {
		name, code := stockBriefData.Name(), fmt.Sprintf("%v.%v", stockBriefData.Code(), market)
		// stockDailyDataSlice = append(stockDailyDataSlice, nil)
		go func(_code, _name string) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("\t- spider stock %v - %v occurs error: %v\n", _code, _name, err)
				}
			}()
			_stockDailyDataSlice := make([]model.StockDailyData, 0, 128)
			stp.NewArray(tushare.LoadStockDailyData(_code, _name)).ForEach(func(v *tushare.TS_StockDailyData, i int) {
				_stockDailyDataSlice = append(_stockDailyDataSlice, &spiderData{StockBriefData: stockBriefData, TS_StockDailyData: v})
				stockDailyDataSlice[index] = _stockDailyDataSlice
			})
			wg.Done()
		}(code, name)
		counter++
		if counter%5 == 0 {
			time.Sleep(time.Second)
		}
	}
	wg.Wait()
	return stockDailyDataSlice
}

// ----------------------------------------------------------------

// DownloadStockBriefData 下载股票简略数据
func DownloadStockBriefData() {
	stockBriefSlice := appfinanceifengcom.DownloadStockSlice()
	appfinanceifengcom.SaveStockList(stockBriefSlice)
}

// OutputStockBriefStatistics 输出股票简略数据的统计结果
// - 沪市股票数量
// - 深市股票数量
func OutputStockBriefStatistics() {
	stockBriefDataSlice := appfinanceifengcom.LoadStockBriefList()
	SHStockCount, SZStockCount := 0, 0
	stp.NewArray(stockBriefDataSlice).ForEach(func(v *appfinanceifengcom.AFI_StockBriefData, i int) {
		switch {
		case filter.SH_StockFilter(v.Code()):
			SHStockCount++
		case filter.SZ_StockFilter(v.Code()):
			SZStockCount++
		}
	})

	formatter := func(market string, count int) {
		fmt.Printf("\t- market: %v, count: %v\n", market, count)
	}

	fmt.Printf("- stock statistics\n")
	formatter(filter.SH_Market, SHStockCount)
	formatter(filter.SZ_Market, SZStockCount)
}

// // AppendStockDailyData 根据现有历史每日行情数据追加每日行情数据直至当前时间
// func AppendStockDailyData() {
// 	stockBriefDataSlice := appfinanceifengcom.LoadStockBriefList()
// 	SHStockBriefData := make([]model.StockBriefData, 0, len(stockBriefDataSlice)/2)
// 	SZStockBriefData := make([]model.StockBriefData, 0, len(stockBriefDataSlice)/2)
// 	stp.NewArray(stockBriefDataSlice).ForEach(func(v *appfinanceifengcom.AFI_StockBriefData, i int) {
// 		switch {
// 		case filter.SH_StockFilter(v.Code()):
// 			SHStockBriefData = append(SHStockBriefData, v)
// 		case filter.SZ_StockFilter(v.Code()):
// 			SZStockBriefData = append(SZStockBriefData, v)
// 		}
// 	})
// 	appendStockDailyData(SHStockBriefData, filter.SH_Market)
// 	appendStockDailyData(SZStockBriefData, filter.SZ_Market)
// }

// func appendStockDailyData(stockBriefDataSlice []model.StockBriefData, market string) {
// 	fmt.Printf("- spider append stock daily data\n")
// 	fmt.Printf("\t- market: %v\n", market)
// 	counter, wg := 0, sync.WaitGroup{}
// 	wg.Add(len(stockBriefDataSlice))
// 	for _, data := range stockBriefDataSlice {
// 		_name, _code := data.Name(), fmt.Sprintf("%v.%v", data.Code(), market)
// 		go func(c, n string) {
// 			defer func() {
// 				if _err := recover(); _err != nil {
// 					fmt.Printf("\t\t- spider stock %v - %v occurs error: %v\n", c, n, _err)
// 				}
// 			}()

// 			var beginDate time.Time
// 			endDate, err := time.Parse(tushare.TradeDateLayout(), time.Now().AddDate(0, 0, 1).Format(tushare.TradeDateLayout()))
// 			if err != nil {
// 				panic(err)
// 			}
// 			_dailyData := tushare.LoadStockDailyData(_code)
// 			if len(_dailyData) == 0 {
// 				beginDate = endDate.AddDate(0, 0, -1)
// 			} else {
// 				beginDate, err = time.Parse(tushare.TradeDateLayout(), _dailyData[0].TS_TradeDate)
// 				if err != nil {
// 					panic(err)
// 				}
// 				beginDate = beginDate.AddDate(0, 0, 1)
// 			}

// 			dailyData := tushare.DownloadDailyData(_code, _name, 0, beginDate.Unix(), endDate.Unix())
// 			dailyData = append(dailyData, _dailyData...)
// 			tushare.SaveStockDailyData(c, n, dailyData)
// 			wg.Done()
// 		}(_code, _name)
// 		counter++
// 		if counter%5 == 0 {
// 			time.Sleep(time.Second)
// 		}
// 	}
// 	wg.Wait()
// 	fmt.Printf("- spider append stock daily data done, count %v\n", len(stockBriefDataSlice))
// }

// // ArchiveStockDailyData 压缩并归档现有历史每日行情数据
// func ArchiveStockDailyData() {
// 	stockBriefDataSlice := appfinanceifengcom.LoadStockBriefList()
// 	stockDailyDataArchivePath := filepath.Join(global.PersonalDocumentPath, global.StockDailyDataArchiveRelativePath)
// 	if stp.IsExist(stockDailyDataArchivePath) {
// 		if err := os.Remove(stockDailyDataArchivePath); err != nil {
// 			panic(err)
// 		}
// 		fmt.Printf("\t\t- spider remove stock daily data archive\n")
// 	}

// 	stockDailyDataArchiveFile, err := os.Create(stockDailyDataArchivePath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer stockDailyDataArchiveFile.Close()

// 	stockDailyDataMap := make(map[string][]*tushare.TS_StockDailyData)
// 	for _, stockBriefData := range stockBriefDataSlice {
// 		var market string
// 		switch {
// 		case filter.SH_StockFilter(stockBriefData.Code()):
// 			market = filter.SH_Market
// 		case filter.SZ_StockFilter(stockBriefData.Code()):
// 			market = filter.SZ_Market
// 		}
// 		if len(market) > 0 {
// 			_code := fmt.Sprintf("%v.%v", stockBriefData.Code(), market)
// 			stockDailyDataSlice := tushare.LoadStockDailyData(_code)
// 			if len(stockDailyDataSlice) == 0 {
// 				panic(fmt.Sprintf("stock %v daily data is empty", _code))
// 			}
// 			stockDailyDataMap[stockBriefData.Code()] = stockDailyDataSlice
// 		}
// 	}

// 	b, err := json.Marshal(stockDailyDataMap)
// 	if err != nil {
// 		panic(err)
// 	}

// 	buffer := &bytes.Buffer{}
// 	writer := gzip.NewWriter(buffer)
// 	writer.Write(b)
// 	defer writer.Close()

// 	_, err = stockDailyDataArchiveFile.Write(buffer.Bytes())
// 	if err != nil {
// 		panic(err)
// 	}
// }

// stock ETF

// DownloadStockETFSlice 下载场内ETF简略数据
func DownloadStockETFSlice() {
	stockETFBriefSlice := fundeastmoney.DownloadStockETFSlice()
	fundeastmoney.SaveStockETFList(stockETFBriefSlice)
}

package spider

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"go-fund/filter"
	"go-fund/global"
	"go-fund/model"
	appfinanceifengcom "go-fund/spider/app.finance.ifeng.com"
	fundeastmoney "go-fund/spider/fund.eastmoney.com"
	"go-fund/spider/tushare.pro"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Mericusta/go-stp"
)

// stock

func DownloadStockDailyData(stockBriefDataSlice []model.StockBriefData) {
	beginDate, err := time.Parse(tushare.TradeDateLayout(), global.Date1)
	if err != nil {
		panic(err)
	}
	endDate := time.Now()
	SHStockBriefData := make([]model.StockBriefData, 0, len(stockBriefDataSlice)/2)
	SZStockBriefData := make([]model.StockBriefData, 0, len(stockBriefDataSlice)/2)
	stp.NewArray(stockBriefDataSlice).ForEach(func(v model.StockBriefData, i int) {
		switch {
		case filter.SH_StockFilter(v.Code()):
			SHStockBriefData = append(SHStockBriefData, v)
		case filter.SZ_StockFilter(v.Code()):
			SZStockBriefData = append(SZStockBriefData, v)
		}
	})
	downloadStockDailyData(SHStockBriefData, filter.SH_Market, beginDate, endDate)
	downloadStockDailyData(SZStockBriefData, filter.SZ_Market, beginDate, endDate)
}

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

// DownloadAllStockDailyData 下载股票历史每日行情数据
func DownloadAllStockDailyData() {
	beginDate, err := time.Parse(tushare.TradeDateLayout(), global.Date1)
	if err != nil {
		panic(err)
	}
	endDate := time.Now()
	stockBriefDataSlice := appfinanceifengcom.LoadStockBriefList()
	SHStockBriefData := make([]model.StockBriefData, 0, len(stockBriefDataSlice)/2)
	SZStockBriefData := make([]model.StockBriefData, 0, len(stockBriefDataSlice)/2)
	stp.NewArray(stockBriefDataSlice).ForEach(func(v *appfinanceifengcom.AFI_StockBriefData, i int) {
		switch {
		case filter.SH_StockFilter(v.Code()):
			SHStockBriefData = append(SHStockBriefData, v)
		case filter.SZ_StockFilter(v.Code()):
			SZStockBriefData = append(SZStockBriefData, v)
		}
	})
	downloadStockDailyData(SHStockBriefData, filter.SH_Market, beginDate, endDate)
	downloadStockDailyData(SZStockBriefData, filter.SZ_Market, beginDate, endDate)
}

func downloadStockDailyData(stockBriefDataSlice []model.StockBriefData, market string, beginDate, endDate time.Time) {
	fmt.Printf("- spider download stock daily data, total %v\n", len(stockBriefDataSlice))
	fmt.Printf("\t- market: %v\n", market)
	fmt.Printf("\t- duration: %v ~ %v\n", beginDate.Format(tushare.TradeDateLayout()), endDate.Format(tushare.TradeDateLayout()))
	counter, wg := 0, sync.WaitGroup{}
	wg.Add(len(stockBriefDataSlice))
	for _, data := range stockBriefDataSlice {
		_name, _code := data.Name(), fmt.Sprintf("%v.%v", data.Code(), market)
		go func(c, n string) {
			defer func() {
				if _err := recover(); _err != nil {
					fmt.Printf("\t\t- spider stock %v - %v occurs error: %v\n", c, n, _err)
				}
			}()
			dailyData := tushare.DownloadDailyData(_code, _name, 0, beginDate.Unix(), endDate.Unix())
			tushare.SaveStockDailyData(c, n, dailyData)
			wg.Done()
		}(_code, _name)
		counter++
		if counter%5 == 0 {
			time.Sleep(time.Second)
		}
	}
	wg.Wait()
	fmt.Printf("- spider download stock daily data done, count %v\n", len(stockBriefDataSlice))
}

// AppendStockDailyData 根据现有历史每日行情数据追加每日行情数据直至当前时间
func AppendStockDailyData() {
	stockBriefDataSlice := appfinanceifengcom.LoadStockBriefList()
	SHStockBriefData := make([]model.StockBriefData, 0, len(stockBriefDataSlice)/2)
	SZStockBriefData := make([]model.StockBriefData, 0, len(stockBriefDataSlice)/2)
	stp.NewArray(stockBriefDataSlice).ForEach(func(v *appfinanceifengcom.AFI_StockBriefData, i int) {
		switch {
		case filter.SH_StockFilter(v.Code()):
			SHStockBriefData = append(SHStockBriefData, v)
		case filter.SZ_StockFilter(v.Code()):
			SZStockBriefData = append(SZStockBriefData, v)
		}
	})
	appendStockDailyData(SHStockBriefData, filter.SH_Market)
	appendStockDailyData(SZStockBriefData, filter.SZ_Market)
}

func appendStockDailyData(stockBriefDataSlice []model.StockBriefData, market string) {
	fmt.Printf("- spider append stock daily data\n")
	fmt.Printf("\t- market: %v\n", market)
	counter, wg := 0, sync.WaitGroup{}
	wg.Add(len(stockBriefDataSlice))
	for _, data := range stockBriefDataSlice {
		_name, _code := data.Name(), fmt.Sprintf("%v.%v", data.Code(), market)
		go func(c, n string) {
			defer func() {
				if _err := recover(); _err != nil {
					fmt.Printf("\t\t- spider stock %v - %v occurs error: %v\n", c, n, _err)
				}
			}()

			var beginDate time.Time
			endDate, err := time.Parse(tushare.TradeDateLayout(), time.Now().AddDate(0, 0, 1).Format(tushare.TradeDateLayout()))
			if err != nil {
				panic(err)
			}
			_dailyData := tushare.LoadStockDailyData(_code)
			if len(_dailyData) == 0 {
				beginDate = endDate.AddDate(0, 0, -1)
			} else {
				beginDate, err = time.Parse(tushare.TradeDateLayout(), _dailyData[0].TS_TradeDate)
				if err != nil {
					panic(err)
				}
				beginDate = beginDate.AddDate(0, 0, 1)
			}

			dailyData := tushare.DownloadDailyData(_code, _name, 0, beginDate.Unix(), endDate.Unix())
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
	fmt.Printf("- spider append stock daily data done, count %v\n", len(stockBriefDataSlice))
}

// ArchiveStockDailyData 压缩并归档现有历史每日行情数据
func ArchiveStockDailyData() {
	stockBriefDataSlice := appfinanceifengcom.LoadStockBriefList()
	stockDailyDataArchivePath := filepath.Join(global.PersonalDocumentPath, global.StockDailyDataArchiveRelativePath)
	if stp.IsExist(stockDailyDataArchivePath) {
		if err := os.Remove(stockDailyDataArchivePath); err != nil {
			panic(err)
		}
		fmt.Printf("\t\t- spider remove stock daily data archive\n")
	}

	stockDailyDataArchiveFile, err := os.Create(stockDailyDataArchivePath)
	if err != nil {
		panic(err)
	}
	defer stockDailyDataArchiveFile.Close()

	stockDailyDataMap := make(map[string][]*tushare.TS_StockDailyData)
	for _, stockBriefData := range stockBriefDataSlice {
		var market string
		switch {
		case filter.SH_StockFilter(stockBriefData.Code()):
			market = filter.SH_Market
		case filter.SZ_StockFilter(stockBriefData.Code()):
			market = filter.SZ_Market
		}
		if len(market) > 0 {
			_code := fmt.Sprintf("%v.%v", stockBriefData.Code(), market)
			stockDailyDataSlice := tushare.LoadStockDailyData(_code)
			if len(stockDailyDataSlice) == 0 {
				panic(fmt.Sprintf("stock %v daily data is empty", _code))
			}
			stockDailyDataMap[stockBriefData.Code()] = stockDailyDataSlice
		}
	}

	b, err := json.Marshal(stockDailyDataMap)
	if err != nil {
		panic(err)
	}

	buffer := &bytes.Buffer{}
	writer := gzip.NewWriter(buffer)
	writer.Write(b)
	defer writer.Close()

	_, err = stockDailyDataArchiveFile.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func LoadStockDailyData() {
	stockDailyDataArchivePath := filepath.Join(global.PersonalDocumentPath, global.StockDailyDataArchiveRelativePath)
	if !stp.IsExist(stockDailyDataArchivePath) {
		panic(stockDailyDataArchivePath)
	}

	b, err := os.ReadFile(stockDailyDataArchivePath)
	if err != nil {
		panic(err)
	}

	reader, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	buffer := &bytes.Buffer{}
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		panic(err)
	}

	stockDailyDataMap := make(map[string][]*tushare.TS_StockDailyData)
	err = json.Unmarshal(buffer.Bytes(), &stockDailyDataMap)
	if err != nil {
		panic(err)
	}

	fmt.Printf("stockDailyDataMap = %v\n", stockDailyDataMap)
}

// stock ETF

// DownloadStockETFSlice 下载场内ETF简略数据
func DownloadStockETFSlice() {
	stockETFBriefSlice := fundeastmoney.DownloadStockETFSlice()
	fundeastmoney.SaveStockETFList(stockETFBriefSlice)
}

package main

import (
	"go-fund/observer"
	"go-fund/spider"
)

type tempStock struct {
	name string
	code string
}

var (
	tempStockList = []tempStock{
		{name: "海天味业", code: "603288"},
		{name: "中科曙光", code: "603019"},
		{name: "卫星化学", code: "002648"},
		{name: "诺力股份", code: "603611"},
		{name: "北新建材", code: "000786"},
		{name: "华泰证券", code: "601688"},
		{name: "中国黄金", code: "600916"},
		{name: "吉比特", code: "603444"},
		{name: "贵州茅台", code: "600519"},
		{name: "五粮液", code: "000858"},
		{name: "创业环保", code: "600874"},
		{name: "春风动力", code: "603129"},
		{name: "荣泰健康", code: "603579"},
		{name: "仙鹤股份", code: "603733"},
		{name: "双环传动", code: "002472"},
		{name: "老白干酒", code: "600559"},
		{name: "招商银行", code: "600036"},
		{name: "中材科技", code: "002080"},
		{name: "建设机械", code: "600984"},
		{name: "捷荣技术", code: "002855"},
	}
)

func main() {
	// 添加待观察股票
	for _, tempStock := range tempStockList {
		observer.AppendStockToObserveList(tempStock.name, tempStock.code)
	}

	// 查找待观察股票的每日数据
	spider.DownloadStockDailyData(observer.LoadObserveStockBriefList())

	// 创建日志文件
	// append := false
	// now := time.Now()
	// logFileName := now.Format(stp.DateLayout) + ".log"
	// logFilePath := filepath.Join(global.PersonalDocumentPath + "/markdown/note/investment/stock/statistics/" + logFileName)
	// if !stp.IsExist(logFilePath) {
	// 	stp.CreateFile(logFilePath)
	// 	append = true
	// }

	// 每个月1号重新下载所有股票数据
	// if now.Day() == 1 {
	// 	spider.DownloadStockBriefData()
	// 	spider.OutputStockBriefStatistics()
	// 	spider.DownloadStockDailyData()
	// }

	// var lastExecuteTime *time.Time
	// err := stp.ReadFileLineOneByOne(logFilePath, func(s string, i int) bool {
	// 	if i == 0 {
	// 		if len(s) > 0 {
	// 			t, err := time.Parse(tushare.TradeDateLayout(), s)
	// 			if err != nil {
	// 				return false
	// 			}
	// 			lastExecuteTime = &t
	// 		}
	// 	}
	// 	return false
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// if now.Hour() < 15 {
	// 	// 每天15点前搜索时排除当天数据
	// 	searcher.SearchAlgorithm1(3, -1, -1, 0, 0)
	// } else {
	// 	// 每天15点后，根据当日执行日志判断是否需要更新当日数据
	// 	if append && now.Hour() >= 15 {
	// 		spider.AppendStockDailyData()
	// 		// spider.ArchiveStockDailyData()
	// 		// spider.LoadStockDailyData()
	// 		// spider.DownloadStockETFSlice()
	// 	}

	// 	// 每天15点后搜索时包括当天数据
	// 	searcher.SearchAlgorithm1(3, 0, -1, 0, 0)
	// }
}

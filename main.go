package main

import (
	"go-fund/global"
	"go-fund/searcher"
	"go-fund/spider"
	"go-fund/spider/tushare.pro"
	"path/filepath"
	"time"

	"github.com/Mericusta/go-stp"
)

func main() {
	// 创建日志文件
	now := time.Now()
	logFileName := now.Format(stp.DateLayout) + ".log"
	logFilePath := filepath.Join(global.PersonalDocumentPath + "/markdown/note/investment/stock/statistics/" + logFileName)
	if !stp.IsExist(logFilePath) {
		stp.CreateFile(logFilePath)
	}

	// 每个月1号重新下载所有股票数据
	if now.Day() == 1 {
		spider.DownloadStockBriefData()
		spider.OutputStockBriefStatistics()
		spider.DownloadStockDailyData()
	}

	var lastExecuteTime *time.Time
	err := stp.ReadFileLineOneByOne(logFilePath, func(s string, i int) bool {
		if i == 0 {
			if len(s) > 0 {
				t, err := time.Parse(tushare.TradeDateLayout(), s)
				if err != nil {
					return false
				}
				lastExecuteTime = &t
			}
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	if now.Hour() < 15 {
		// 每天15点前搜索时排除当天数据
		searcher.SearchAlgorithm1(3, -1, -1, 0, 0)
	} else {
		// 每天15点后，根据当日执行日志判断是否需要更新当日数据
		if lastExecuteTime == nil || lastExecuteTime.Hour() < 15 {
			// spider.AppendStockDailyData()
			// spider.ArchiveStockDailyData()
			// spider.LoadStockDailyData()
			// spider.DownloadStockETFSlice()
		}

		// 每天15点后搜索时包括当天数据
		searcher.SearchAlgorithm1(3, 0, -1, 0, 0)
	}
}

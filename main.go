package main

import (
	"flag"
	"fmt"
	"go-fund/global"
	"go-fund/observer"
	"go-fund/spider/tushare.pro"
	"strings"
	"time"
)

var date = flag.String("date", time.Now().Format(tushare.TradeDateLayout()), "spider and calculate date")
var splitter = func() { fmt.Println(strings.Repeat("-", 64)) }

type tempStock struct {
	name string
	code string
}

var (
	tempStockList = []tempStock{
		{name: "百川股份", code: "002455"},
		{name: "中国石油", code: "601857"},
		{name: "工业富联", code: "601138"},
		{name: "英维克", code: "002837"},
		{name: "立讯精密", code: "002475"},
		{name: "海天味业", code: "603288"},
		{name: "中科曙光", code: "603019"},
		{name: "诺力股份", code: "603611"},
		{name: "华泰证券", code: "601688"},
		{name: "中国黄金", code: "600916"},
		{name: "吉比特", code: "603444"},
		{name: "贵州茅台", code: "600519"},
		{name: "创业环保", code: "600874"},
		{name: "春风动力", code: "603129"},
		{name: "荣泰健康", code: "603579"},
		{name: "仙鹤股份", code: "603733"},
		{name: "老白干酒", code: "600559"},
		{name: "招商银行", code: "600036"},
		{name: "建设机械", code: "600984"},
		{name: "卫星化学", code: "002648"},
		{name: "北新建材", code: "000786"},
		{name: "五粮液", code: "000858"},
		{name: "双环传动", code: "002472"},
		{name: "中材科技", code: "002080"},
		{name: "捷荣技术", code: "002855"},
		{name: "纳指科技ETF", code: "159509"},
		{name: "道琼斯ETF", code: "513400"},
		{name: "标普ETF", code: "159655"},
	}
)

func init() {
	flag.Parse()
}

// initObserver 初始化待观察数据
func initObserver() {
	// 清空待观察列表中的股票
	observer.ClearObserveStockList()

	// 添加待观察股票
	for _, tempStock := range tempStockList {
		observer.AppendStockToObserveStockList(tempStock.name, tempStock.code)
	}
}

// downloadObserverStockDailyData 下载待观察股票每日行情数据，从 20071016 至 endDate
func downloadObserverStockDailyData(endDateStr string) {
	beginDate, err := time.Parse(tushare.TradeDateLayout(), global.Date1)
	if err != nil {
		panic(err)
	}
	endDate, err := time.Parse(tushare.TradeDateLayout(), endDateStr)
	if err != nil {
		panic(err)
	}
	observer.DownloadObserveStockDailyData(beginDate, endDate)
}

// calculateObserverStockMAData 计算待观察股票的 MA 数据
func calculateObserverStockMAData(calculateDateStr string) {
	calculateDate, err := time.Parse(tushare.TradeDateLayout(), calculateDateStr)
	if err != nil {
		panic(err)
	}
	observer.CalculateObserverStockMAData(calculateDate)
}

func main() {
	// initObserver()
	downloadObserverStockDailyData(*date)
	calculateObserverStockMAData(*date)
}

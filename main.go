package main

import (
	"go-fund/global"
	"go-fund/spider"
	"time"
)

func main() {
	beginDate, err := time.Parse("20060102", global.Date1)
	if err != nil {
		panic(err)
	}
	endDate, err := time.Parse("20060102", "20230823")
	if err != nil {
		panic(err)
	}

	spider.LoadStockDailyData("600536.SH", beginDate, endDate)
}

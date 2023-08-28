package main

import (
	"go-fund/spider"
)

func main() {
	// beginDate, err := time.Parse("20060102", global.Date1)
	// if err != nil {
	// 	panic(err)
	// }
	// endDate, err := time.Parse("20060102", "20230823")
	// if err != nil {
	// 	panic(err)
	// }

	spider.AppendStockDailyData()
}

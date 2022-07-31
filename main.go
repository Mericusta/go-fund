package main

import (
	"fmt"
	fundeastmoney "go-fund/fund.eastmoney.com"
	"go-fund/global"
)

func main() {
	for name, code := range global.FundNameCodeMap {
		date, num := fundeastmoney.GetFundInfoByCode(code)
		pngName := fundeastmoney.GetFuncChartsByCode(code)
		fmt.Printf("name %v date %v num %v png %v\n", name, date, num, pngName)
	}
}

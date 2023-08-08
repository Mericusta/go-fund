package main

import (
	"go-fund/legulegu.com"
	"go-fund/searcher"
)

func main() {
	md := &legulegu.MockData{}
	md.Parse("./resource/202308072334.json")

	stockData := make([]searcher.SearchMethod1Stock, 0, len(md.MockDataList))
	for _, d := range md.MockDataList {
		stockData = append(stockData, d)
	}
	searcher.SearchMethod1(stockData, 2)
}

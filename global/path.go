package global

import "os"

const (
	StockListRelativePath             string = "markdown/note/investment/stock/data/stock_list.json"
	LocalStockListDataPath            string = "../stock_list"
	StockDailyDataRelativePathFormat  string = "markdown/note/investment/stock/data/daily/%v.json"
	StockTradeSimulateRelativePath    string = "markdown/note/investment/stock/trade_simulate"
	StockDailyDataArchiveRelativePath string = "markdown/note/investment/stock/data/archive.json"
)

var (
	PersonalDocumentPath string = os.Getenv("PD")
)

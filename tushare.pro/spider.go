package tushare

import (
	"encoding/json"
	"fmt"
	"go-fund/global"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Mericusta/go-stp"
)

var (
	tokenRelativePath string = "markdown/note/Tushare/token"
	token             string
	url               string = "http://api.tushare.pro"
)

func init() {
	token = LoadToken()
	fmt.Printf("tushare.pro token = |%v|\n", token)
}

func LoadToken() string {
	b, e := os.ReadFile(filepath.Join(global.PersonalDocumentPath, tokenRelativePath))
	if e != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

type postRequest struct {
	ApiName string            `json:"api_name"`
	Token   string            `json:"token"`
	Params  map[string]string `json:"params"`
	Fields  string            `json:"fields"`
}

type postResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Fields []string `json:"fields"`
		Items  [][]any  `json:"items"`
	} `json:"data"`
}

type StockDailyData struct {
	TSCode    string  `json:"ts_code"`    // 股票代码
	TradeDate string  `json:"trade_date"` // 交易日期
	Open      float64 `json:"open"`       // 开盘价格
	Close     float64 `json:"close"`      // 收盘价格
	High      float64 `json:"high"`       // 最高价格
	Low       float64 `json:"low"`        // 最低价格
	PreClose  float64 `json:"pre_close"`  // 昨日收盘（前复权）
	Change    float64 `json:"change"`     // 涨跌额度
	PctChg    float64 `json:"pct_chg"`    // 涨跌幅度
	Vol       float64 `json:"vol"`        // 成交量（手）
	Amount    float64 `json:"amount"`     // 成交额（千元）
}

func GetDailyData(tsCode string, tradeDate, startDate, endDate int64) []*StockDailyData {
	apiName := "daily"
	params := make(map[string]string)
	params["ts_code"] = tsCode
	if tradeDate > 0 {
		params["trade_date"] = time.Unix(tradeDate, 0).Format("20060102")
	}
	if startDate > 0 {
		params["start_date"] = time.Unix(startDate, 0).Format("20060102")
	}
	if endDate > 0 {
		params["end_date"] = time.Unix(endDate, 0).Format("20060102")
	}

	req := &postRequest{
		ApiName: apiName,
		Token:   token,
		Params:  params,
	}
	resp, err := global.HTTPClient.R().SetBody(req).Post(url)
	if err != nil {
		panic(err)
	}

	content := resp.Body()
	rep := &postResponse{}
	err = json.Unmarshal(content, rep)
	if err != nil {
		panic(err)
	}

	return stp.ReflectStructValueSlice[StockDailyData](rep.Data.Fields, rep.Data.Items, "json")
}

var (
	stockDailyDataRelativePathFormat string = "markdown/note/stock/%v.json"
)

func SaveStockDailyData(code string, slice []*StockDailyData) {
	stockDailyDataPath := filepath.Join(global.PersonalDocumentPath, fmt.Sprintf(stockDailyDataRelativePathFormat, code))
	if stp.IsExist(stockDailyDataPath) {
		if err := os.Remove(stockDailyDataPath); err != nil {
			panic(err)
		}
	}
	stockDailyDataFile, err := os.Create(stockDailyDataPath)
	if err != nil {
		panic(err)
	}
	defer stockDailyDataFile.Close()

	b, err := json.Marshal(slice)
	if err != nil {
		panic(err)
	}

	_, err = stockDailyDataFile.Write(b)
	if err != nil {
		panic(err)
	}
}

func LoadStockDailyData(code string) []*StockDailyData {
	stockDailyDataPath := filepath.Join(global.PersonalDocumentPath, fmt.Sprintf(stockDailyDataRelativePathFormat, code))
	if !stp.IsExist(stockDailyDataPath) {
		return nil
	}
	stockDailyData, err := os.ReadFile(stockDailyDataPath)
	if err != nil {
		panic(err)
	}

	slice := make([]*StockDailyData, 0, 1024)
	err = json.Unmarshal(stockDailyData, &slice)
	if err != nil {
		panic(err)
	}

	return slice
}

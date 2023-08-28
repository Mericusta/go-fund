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
	tokenRelativePath                    string = "markdown/note/Tushare/token"
	token                                string
	url                                  string = "http://api.tushare.pro"
	tokenRequestTimesLimitForEverySecond int    = 500  // 每分钟最多请求500次
	tokenDataCountLimitForEveryRequest   int    = 6000 // 每次最多6000条数据（23年交易日历史数据）
	ticker                               time.Ticker
	tickerRequestCount                   int64
	TradeDateLayout                      string = strings.Replace(stp.DateLayout, "-", "", -1)
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

func MakeStockTSCode(code, market string) string {
	return fmt.Sprintf("%v.%v", code, market)
}

// TS_StockDailyData 字段类型顺序必须和 postResponse.Data.Items 保持一致
// 实现 searcher.StockDailyData 接口
type TS_StockDailyData struct {
	TS_Code      string  `json:"ts_code"`    // 股票代码
	TS_TradeDate string  `json:"trade_date"` // 交易日期
	TS_Open      float64 `json:"open"`       // 开盘价格
	TS_Close     float64 `json:"close"`      // 收盘价格
	TS_High      float64 `json:"high"`       // 最高价格
	TS_Low       float64 `json:"low"`        // 最低价格
	TS_PreClose  float64 `json:"pre_close"`  // 昨日收盘（前复权）
	TS_Change    float64 `json:"change"`     // 涨跌额度
	TS_PctChg    float64 `json:"pct_chg"`    // 涨跌幅度
	TS_Vol       float64 `json:"vol"`        // 成交量（手）
	TS_Amount    float64 `json:"amount"`     // 成交额（千元）
}

func (sdd *TS_StockDailyData) Code() string   { return sdd.TS_Code[:strings.Index(sdd.TS_Code, ".")] }
func (sdd *TS_StockDailyData) Market() string { return sdd.TS_Code[strings.Index(sdd.TS_Code, ".")+1:] }
func (sdd *TS_StockDailyData) Date() time.Time {
	t, _ := time.Parse(TradeDateLayout, sdd.TS_TradeDate)
	return t
}
func (sdd *TS_StockDailyData) OpenValue() float64  { return sdd.TS_Open }
func (sdd *TS_StockDailyData) CloseValue() float64 { return sdd.TS_Close }
func (sdd *TS_StockDailyData) HighValue() float64  { return sdd.TS_High }
func (sdd *TS_StockDailyData) LowValue() float64   { return sdd.TS_Low }
func (sdd *TS_StockDailyData) Volume() float64     { return sdd.TS_Vol }
func (sdd *TS_StockDailyData) Amount() float64     { return sdd.TS_Amount }

func GetDailyData(code, name string, tradeDate, startDate, endDate int64) []*TS_StockDailyData {
	fmt.Printf("\t\t- spider get stock %v - %v daily data\n", code, name)
	apiName := "daily"
	params := make(map[string]string)
	params["ts_code"] = code
	if tradeDate > 0 {
		params["trade_date"] = time.Unix(tradeDate, 0).Format(TradeDateLayout)
	}
	if startDate > 0 {
		params["start_date"] = time.Unix(startDate, 0).Format(TradeDateLayout)
	}
	if endDate > 0 {
		params["end_date"] = time.Unix(endDate, 0).Format(TradeDateLayout)
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

	return stp.ReflectStructValueSlice[TS_StockDailyData](rep.Data.Fields, rep.Data.Items, "json")
}

var (
	stockDailyDataRelativePathFormat string = "markdown/note/stock/daily/%v.json"
)

func SaveStockDailyData(code, name string, slice []*TS_StockDailyData) {
	fmt.Printf("\t\t- spider save stock %v - %v daily data\n", code, name)
	stockDailyDataPath := filepath.Join(global.PersonalDocumentPath, fmt.Sprintf(stockDailyDataRelativePathFormat, code))
	if stp.IsExist(stockDailyDataPath) {
		if err := os.Remove(stockDailyDataPath); err != nil {
			panic(err)
		}
		fmt.Printf("\t\t- spider remove stock %v - %v daily data\n", code, name)
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

func AppendDailyData(code, name string, slice []*TS_StockDailyData) {
	fmt.Printf("\t\t- spider append stock %v - %v daily data\n", code, name)
	stockDailyDataPath := filepath.Join(global.PersonalDocumentPath, fmt.Sprintf(stockDailyDataRelativePathFormat, code))
	var _slice []*TS_StockDailyData
	if stp.IsExist(stockDailyDataPath) {
		_slice = LoadStockDailyData(code)
	} else {
		_, err := os.Create(stockDailyDataPath)
		if err != nil {
			panic(err)
		}
	}
	slice = append(slice, _slice...)

	b, err := json.Marshal(slice)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(stockDailyDataPath, b, 0644)
	if err != nil {
		panic(err)
	}
}

func LoadStockDailyData(code string) []*TS_StockDailyData {
	stockDailyDataPath := filepath.Join(global.PersonalDocumentPath, fmt.Sprintf(stockDailyDataRelativePathFormat, code))
	if !stp.IsExist(stockDailyDataPath) {
		return nil
	}
	stockDailyData, err := os.ReadFile(stockDailyDataPath)
	if err != nil {
		panic(err)
	}

	slice := make([]*TS_StockDailyData, 0, 1024)
	err = json.Unmarshal(stockDailyData, &slice)
	if err != nil {
		panic(err)
	}

	return slice
}

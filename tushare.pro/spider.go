package tusharepro

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
	TSCode    string  `json:"ts_code"`
	TradeDate string  `json:"trade_date"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	PreClose  float64 `json:"pre_close"`
	Change    float64 `json:"change"`
	PctChg    float64 `json:"pct_chg"`
	Vol       float64 `json:"vol"`
	Amount    float64 `json:"amount"`
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

	return stp.AssignStructMember[StockDailyData](rep.Data.Fields, rep.Data.Items, "json")
}

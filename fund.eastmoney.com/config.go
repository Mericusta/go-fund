package fundeastmoney

import (
	"strings"
	"time"
)

var (
	fund_REPLACE_KEYWORD_HOME_URL   = "[RP_HOME_URL]"
	fund_REPLACE_KEYWORD_FUND_CODE  = "[RP_FUND_CODE]"
	fund_REPLACE_KEYWORD_CHART_TIME = "[RP_CHART_TIME]"

	// 主页
	homeUrl string = "https://fund.eastmoney.com"

	// 基金信息页
	fundInfoUrlTemplate string = "[RP_HOME_URL]/[RP_FUND_CODE].html"

	// 基金净值估算图
	fundChartUrlTemplate string = "http://j4.dfcfw.com/charts/pic6/[RP_FUND_CODE].png?v=[RP_CHART_TIME]"
)

func getFundInfoUrlByCode(code string) string {
	fundInfoUrl := strings.Replace(fundInfoUrlTemplate, fund_REPLACE_KEYWORD_HOME_URL, homeUrl, -1)
	fundInfoUrl = strings.Replace(fundInfoUrl, fund_REPLACE_KEYWORD_FUND_CODE, code, -1)
	return fundInfoUrl
}

func getFundChartUrlByCode(code string) string {
	fundChartUrl := strings.Replace(fundChartUrlTemplate, fund_REPLACE_KEYWORD_FUND_CODE, code, -1)
	fundChartUrl = strings.Replace(fundChartUrl, fund_REPLACE_KEYWORD_CHART_TIME, time.Now().Format("20060102150405"), -1)
	return fundChartUrl
}

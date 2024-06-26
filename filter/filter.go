package filter

func marketFilter(m map[string]string, checker func(string) bool) map[string]string {
	_m := make(map[string]string)
	for code, name := range m {
		if checker(code) {
			_m[code] = name
		}
	}
	return _m
}

const (
	SH_Market  string = "SH" // 沪市
	SSE_Market string = "SH" // 沪市科创板
	SZ_Market  string = "SZ" // 深市
	GEM_Market string = "SZ" // 深市创业板
)

// SH_StockFilter 沪市主板股票筛选器
func SH_StockFilter(code string) bool {
	return rune(code[0]) == '6' && !SSE_StockFilter(code)
}

// SH_ETF_FundFilter 沪市 ETF 基金
func SH_ETF_FundFilter(code string) bool {
	return code[0:2] == "51"
}

// SSE_StockFilter 沪市科创板股票筛选器
func SSE_StockFilter(code string) bool {
	return code[0:3] == "688"
}

// SZ_StockFilter 深市主板股票筛选器
func SZ_StockFilter(code string) bool {
	return rune(code[0]) == '0'
}

// SZ_ETF_FundFilter 深市 ETF 基金
func SZ_ETF_FundFilter(code string) bool {
	return code[0:2] == "15"
}

// GEM_StockFilter 深市创业板股票筛选器
func GEM_StockFilter(m map[string]string) (map[string]string, string) {
	return marketFilter(m, func(s string) bool { return rune(s[0]) == '3' }), "SZ"
}

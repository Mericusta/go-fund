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

// SH_StockFilter 沪市主板股票筛选器
func SH_StockFilter(m map[string]string) (map[string]string, string) {
	return marketFilter(m, func(s string) bool { return rune(s[0]) == '6' }), "SH"
}

// SZ_StockFilter 深市主板股票筛选器
func SZ_StockFilter(m map[string]string) (map[string]string, string) {
	return marketFilter(m, func(s string) bool { return rune(s[0]) == '0' }), "SZ"
}

// GEM_StockFilter 深市创业板股票筛选器
func GEM_StockFilter(m map[string]string) (map[string]string, string) {
	return marketFilter(m, func(s string) bool { return rune(s[0]) == '3' }), "SZ"
}

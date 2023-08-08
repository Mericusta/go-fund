package searcher

import "time"

// ---------------- search method 1 ----------------

type SearchMethod1Stock interface {
	OpenValue() float32
	CloseValue() float32
}

type SearchMethod1Result struct {
	Day1  SearchMethod1Stock
	Day2  SearchMethod1Stock
	Delta float32
}

// SearchMethod1 第二天的开盘价相比第一天收盘价的差值百分比
// @param1       百分比绝对值
func SearchMethod1(stockData []SearchMethod1Stock, delta float32) []*SearchMethod1Result {
	l := len(stockData)
	if l < 2 {
		return nil
	}
	result := make([]*SearchMethod1Result, 0, l/2)
	for i := 0; i+1 < l; i++ {
		d1, d2 := stockData[i], stockData[i+1]
		if d1 == nil || d2 == nil {
			continue
		}
		var _delta, _v float32
		if d1.CloseValue() > d2.OpenValue() {
			_delta, _v = (d1.CloseValue()/d2.OpenValue()-1)*100, -1
		} else if d1.CloseValue() < d2.OpenValue() {
			_delta, _v = (d2.OpenValue()/d1.CloseValue()-1)*100, 1
		} else {
			continue
		}
		if _delta < delta {
			continue
		}
		result = append(result, &SearchMethod1Result{
			Day1:  d1,
			Day2:  d2,
			Delta: _delta * _v,
		})
	}
	return result
}

// ---------------- search method 2 ----------------

type SearchMethod2Stock interface {
	DateTime() time.Time
}

func SearchMethod2(stockData []SearchMethod2Stock) {

}

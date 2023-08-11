package simulation

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"golang.org/x/exp/constraints"
)

type player struct {
	holdCount           int
	holdValue           int
	cost                int
	release             int
	tradeWhilePercent   int
	tradeCount          int
	tradeTimes          int
	lastTradeDay        int
	maxNonTradeDuration [2]int
	forceClear          bool
}

type tradeRecord struct {
	sell  bool
	count int
	value int
}

type simulateStatistics struct {
	days                  int
	upDays                int
	cumulativeUpPercent   int
	downDays              int
	cumulativeDownPercent int
}

type chartStatistics struct {
	title    string
	subTitle string
	xAxis    []interface{}
	yAxisMap map[string][]interface{}
}

func (cs *chartStatistics) makeChart() {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: cs.title, Subtitle: cs.subTitle}),
	)
	line.SetXAxis(cs.xAxis)
	for series, yData := range cs.yAxisMap {
		yAxis := make([]opts.LineData, 0, len(yData))
		for _, data := range yData {
			yAxis = append(yAxis, opts.LineData{Value: data})
		}
		line.AddSeries(series, yAxis)
	}

	page := components.NewPage()
	page.AddCharts(line)
	f, err := os.Create("page.html")
	if err != nil {
		panic(err)
	}
	err = page.Render(io.MultiWriter(f))
	if err != nil {
		panic(err)
	}
}

type totalStatistics struct {
	locker                sync.Mutex
	wg                    sync.WaitGroup
	count                 int
	equalCount            int
	upCount               int
	upDays                int
	upCumulativePercent   int
	downCount             int
	maxUp                 int
	maxDown               int
	downDays              int
	downCumulativePercent int
	maxNonTradeDuration   []int
}

func getDeltaPercent(deltaPercentRange int) int {
	r := random.Intn(deltaPercentRange) + 1
	if (random.Intn(100)+1)%2 == 0 {
		return r // value up
	}
	return -r // value down
}

func (p *player) trade(deltaPercent, tradeValue int) *tradeRecord {
	if deltaPercent < 0 && deltaPercent <= -p.tradeWhilePercent {
		// TODO: h.holdValue > dayValue
		return p.tradeBuy(tradeValue)
	} else if deltaPercent > 0 && deltaPercent > p.tradeWhilePercent {
		return p.tradeSell(tradeValue)
	} else {
		if deltaPercent == 0 {
			panic("deltaPercent is 0")
		}
	}
	return nil
}

func (p *player) tradeBuy(tradeValue int) *tradeRecord {
	_tradeCount := p.tradeCount
	if p.release < p.tradeCount*tradeValue {
		_tradeCount = p.release / tradeValue
	}
	if _tradeCount > 0 {
		p.cost += _tradeCount * tradeValue
		p.release -= _tradeCount * tradeValue
		p.holdCount += _tradeCount
		if p.holdCount > 0 {
			p.holdValue = p.cost / p.holdCount
		}
		return &tradeRecord{sell: false, count: _tradeCount, value: tradeValue}
	}
	return nil
}

func (p *player) tradeSell(tradeValue int) *tradeRecord {
	_tradeCount := p.tradeCount
	if p.holdCount < p.tradeCount {
		_tradeCount = p.holdCount
	}
	if _tradeCount > 0 {
		p.cost -= _tradeCount * tradeValue
		p.release += _tradeCount * tradeValue
		p.holdCount -= _tradeCount
		if p.holdCount > 0 {
			p.holdValue = p.cost / p.holdCount
		}
		return &tradeRecord{sell: true, count: _tradeCount, value: tradeValue}
	}
	return nil
}

func (p *player) clear(tradeValue int) {
	if p.holdCount < 0 {
		panic("hold count cannot be negative")
	}
	p.cost -= p.holdCount * tradeValue
	p.release += p.holdCount * tradeValue
	p.holdCount -= p.holdCount
	p.holdValue = 0
}

func tick(ts *totalStatistics, output, interactive bool) {
	var (
		day        int = 30
		initValue  int = 10000 // 100.00 * 100
		upWeight   int = 50
		downWeight int = 50
		p              = &player{
			holdCount:           0,
			holdValue:           0,
			cost:                0,
			release:             1000000,
			tradeWhilePercent:   3,
			tradeCount:          10,
			tradeTimes:          0,
			maxNonTradeDuration: [...]int{-1, -1},
			forceClear:          true,
		}
		deltaPercentRange int = 5
		// tradeSlice        []*record = make([]*record, day)
		cs = &chartStatistics{
			title:    fmt.Sprintf("%v days simulation", day),
			subTitle: fmt.Sprintf("%v%% percent up/down %v%%", upWeight*100/(upWeight+downWeight), deltaPercentRange),
			xAxis:    []interface{}{0},
			yAxisMap: map[string][]interface{}{
				"end":          {initValue},
				"deltaValue":   {0},
				"deltaPercent": {0},
			},
		}
		input = bufio.NewScanner(os.Stdin)
	)

	if interactive {
		fmt.Printf("init hold value:")
		p.holdCount = terminalInput[int](input)
		p.holdValue = terminalInput[int](input)
		p.cost = terminalInput[int](input)
		p.release = terminalInput[int](input)
	}

	if output {
		fmt.Printf("---------------- simulate init ----------------\n")
		fmt.Printf("player:\n")
		fmt.Printf("- hold count: %v\n", p.holdCount)
		fmt.Printf("- hold value: %v\n", p.holdValue)
		fmt.Printf("- cost: %v\n", p.cost)
		fmt.Printf("- release: %v\n", p.release)
		fmt.Printf("- result: %v\n", p.holdCount*p.holdValue-p.cost)
	}

	if output {
		fmt.Printf("---------------- simulating ----------------\n")
	}
	dayValue := initValue
	ss := &simulateStatistics{days: day}
	for d := 0; d < day; d++ {
		deltaPercent := getDeltaPercent(deltaPercentRange)
		if deltaPercent > 0 {
			ss.upDays++
			ss.cumulativeUpPercent += deltaPercent
		} else if deltaPercent < 0 {
			ss.downDays++
			ss.cumulativeDownPercent += (-deltaPercent)
		} else {
			panic("deltaPercent is 0")
		}
		deltaValue := dayValue * deltaPercent / 100
		begin, end := dayValue, dayValue+deltaValue
		dayValue = end
		record := p.trade(deltaPercent, dayValue)
		if record == nil {
			if p.maxNonTradeDuration[0] == -1 && p.maxNonTradeDuration[1] == -1 {
				p.maxNonTradeDuration[0] = d
			} else if p.maxNonTradeDuration[0] != -1 && p.maxNonTradeDuration[1] == -1 {
				p.maxNonTradeDuration[1] = d
			} else if p.maxNonTradeDuration[0] != -1 && p.maxNonTradeDuration[1] != -1 {
				if p.maxNonTradeDuration[1]-p.maxNonTradeDuration[0] < d-p.lastTradeDay {
					p.maxNonTradeDuration[0] = p.lastTradeDay + 1
					p.maxNonTradeDuration[1] = d
				}
			}
		} else {
			p.lastTradeDay = d
		}
		if output {
			fmt.Printf("---------------- day %v ----------------\n", d)
			fmt.Printf("- info\n")
			fmt.Printf("\t- begin %v\n", begin)
			fmt.Printf("\t- end %v\n", end)
			fmt.Printf("\t- percent %v\n", deltaPercent)
			fmt.Printf("- trade\n")
			if record == nil {
				fmt.Printf("\t- no trade\n")
			} else {
				op := "buy"
				if record.sell {
					op = "sell"
				}
				fmt.Printf("\t- %v count %v with value %v\n", op, record.count, record.value)
			}
			fmt.Printf("- player\n")
			fmt.Printf("\t- hold count %v\n", p.holdCount)
			fmt.Printf("\t- hold value %v\n", p.holdValue)
			fmt.Printf("\t- day value %v\n", p.holdCount*dayValue)
			fmt.Printf("\t- release %v\n", p.release)
			fmt.Printf("\t- cost %v\n", p.cost)

			// charts
			cs.xAxis = append(cs.xAxis, d+1)
			cs.yAxisMap["end"] = append(cs.yAxisMap["end"], end)
			cs.yAxisMap["deltaValue"] = append(cs.yAxisMap["deltaValue"], deltaValue)
			cs.yAxisMap["deltaPercent"] = append(cs.yAxisMap["deltaPercent"], deltaPercent)
		}
	}

	if p.forceClear {
		if output {
			fmt.Printf("---------------- simulate force clear ----------------\n")
		}
		p.clear(dayValue)
	}

	if output {
		fmt.Printf("---------------- simulate end ----------------\n")
		fmt.Printf("player:\n")
		fmt.Printf("- hold count: %v\n", p.holdCount)
		fmt.Printf("- hold value: %v\n", p.holdValue)
		fmt.Printf("- cost: %v\n", p.cost)
		fmt.Printf("- release: %v\n", p.release)
		fmt.Printf("- result: %v\n", p.holdCount*p.holdValue-p.cost)
		fmt.Printf("---------------- simulate statistics ----------------\n")
		fmt.Printf("records: %v days\n", ss.days)
		fmt.Printf("- up: %v days, %v%%\n", ss.upDays, ss.cumulativeUpPercent)
		fmt.Printf("- down: %v days, %v%%\n", ss.downDays, ss.cumulativeDownPercent)
		cs.makeChart()
	}

	ts.locker.Lock()
	if delta := p.holdCount*p.holdValue - p.cost; delta > 0 {
		ts.upCount++
		if ts.maxUp < delta {
			ts.maxUp = delta
		}
	} else if delta < 0 {
		ts.downCount++
		if ts.maxDown > delta {
			ts.maxDown = delta
		}
	} else {
		ts.equalCount++
	}
	ts.upDays += ss.upDays
	ts.upCumulativePercent += ss.cumulativeUpPercent
	ts.downDays += ss.downDays
	ts.downCumulativePercent += ss.cumulativeDownPercent
	if maxNonTradeDuration := p.maxNonTradeDuration[1] - p.maxNonTradeDuration[0]; maxNonTradeDuration > 0 {
		ts.maxNonTradeDuration = append(ts.maxNonTradeDuration, maxNonTradeDuration)
		if maxNonTradeDuration >= 28 {
			panic("here")
		}
	}
	ts.locker.Unlock()
	ts.wg.Done()
}

func terminalInput[T constraints.Integer](input *bufio.Scanner) T {
	var s string
	if input.Scan() {
		s = strings.TrimSpace(input.Text())
	}

	var v T
	if len(s) == 0 {
		return v
	}

	f := func() interface{} {
		_v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(err)
		}
		return _v
	}

	return f().(T)
}

var random *rand.Rand

func simulate() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))

	var (
		output        = false
		interactive   = false
		simulateCount = 100000
		ts            = &totalStatistics{
			count:  simulateCount,
			locker: sync.Mutex{},
			wg:     sync.WaitGroup{},
		}
	)

	if output || interactive {
		simulateCount = 1
	}

	ts.wg.Add(simulateCount)

	for simulateIndex := 0; simulateIndex < simulateCount; simulateIndex++ {
		go tick(ts, output, interactive)
	}

	ts.wg.Wait()

	fmt.Printf("---------------- simulate statistics ----------------\n")
	fmt.Printf("simulate: %v\n", ts.count)
	fmt.Printf("- up: %v times, max %v, days percents %v%%, cumulative %v, avg %v\n", ts.upCount, ts.maxUp, ts.upDays*100/(ts.upDays+ts.downDays), ts.upCumulativePercent, ts.upCumulativePercent/ts.upDays)
	fmt.Printf("- down: %v times, max %v, days percents %v%%, cumulative %v avg %v\n", ts.downCount, ts.maxDown, ts.downDays*100/(ts.upDays+ts.downDays), ts.downCumulativePercent, ts.downCumulativePercent/ts.downDays)
	fmt.Printf("- equal: %v\n", ts.equalCount)
	s := 0
	over7DaysNonTradeDurationCount := 0
	over14DaysNonTradeDurationCount := 0
	over21DaysNonTradeDurationCount := 0
	over30DaysNonTradeDurationCount := 0
	maxNonTradeDuration := 0
	averageMaxNonTradeDuration := 0
	for _, d := range ts.maxNonTradeDuration {
		s += d
		if 7 <= d && d < 14 {
			over7DaysNonTradeDurationCount++
		}
		if 14 <= d && d < 21 {
			over14DaysNonTradeDurationCount++
		}
		if 21 <= d && d < 30 {
			over21DaysNonTradeDurationCount++
		}
		if 30 <= d {
			over30DaysNonTradeDurationCount++
		}
		if maxNonTradeDuration < d {
			maxNonTradeDuration = d
		}
	}
	if len(ts.maxNonTradeDuration) > 0 {
		averageMaxNonTradeDuration = s / len(ts.maxNonTradeDuration)
	}

	fmt.Printf("- non-trade\n")
	fmt.Printf("\t- max duration: %v\n", maxNonTradeDuration)
	fmt.Printf("\t- average duration: %v\n", averageMaxNonTradeDuration)
	fmt.Printf("\t- over 7 days duration: %v\n", over7DaysNonTradeDurationCount)
	fmt.Printf("\t- over 14 days duration: %v\n", over14DaysNonTradeDurationCount)
	fmt.Printf("\t- over 21 days duration: %v\n", over21DaysNonTradeDurationCount)
	fmt.Printf("\t- over 30 days duration: %v\n", over30DaysNonTradeDurationCount)
}

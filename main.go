package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type player struct {
	holdCount         int
	holdValue         int
	cost              int
	release           int
	tradeWhilePercent int
	tradeCount        int
	forceClear        bool
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
}

func getDeltaPercent(deltaPercentRange int) int {
	r := rand.Intn(deltaPercentRange) + 1
	if (rand.Intn(100)+1)%2 == 0 {
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

func simulate(ts *totalStatistics, output bool) {
	var (
		day       int = 30
		initValue int = 10000 // 100.00 * 100
		p             = &player{
			holdCount:         0,
			holdValue:         0,
			cost:              0,
			release:           1000000,
			tradeWhilePercent: 5,
			tradeCount:        10,
			forceClear:        true,
		}
		deltaPercentRange int = 10
		// tradeSlice        []*record = make([]*record, day)
	)

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
	ts.locker.Unlock()
	ts.wg.Done()
}

func main() {
	// for name, code := range global.FundNameCodeMap {
	// 	date, num := fundeastmoney.GetFundInfoByCode(code)
	// 	pngName := fundeastmoney.GetFuncChartsByCode(code)
	// 	fmt.Printf("name %v date %v num %v png %v\n", name, date, num, pngName)
	// }

	rand.Seed(time.Now().UnixNano())

	var (
		output        = true
		simulateCount = 100000
		ss            = &totalStatistics{
			count:  simulateCount,
			locker: sync.Mutex{},
			wg:     sync.WaitGroup{},
		}
	)

	if output {
		simulateCount = 1
	}

	ss.wg.Add(simulateCount)

	for simulateIndex := 0; simulateIndex < simulateCount; simulateIndex++ {
		go simulate(ss, output)
	}

	ss.wg.Wait()

	fmt.Printf("---------------- simulate statistics ----------------\n")
	fmt.Printf("simulate: %v\n", ss.count)
	fmt.Printf("- up: %v times, max %v, days percents %v%%, cumulative %v, avg %v\n", ss.upCount, ss.maxUp, ss.upDays*100/(ss.upDays+ss.downDays), ss.upCumulativePercent, ss.upCumulativePercent/ss.upDays)
	fmt.Printf("- down: %v times, max %v, days percents %v%%, cumulative %v avg %v\n", ss.downCount, ss.maxDown, ss.downDays*100/(ss.upDays+ss.downDays), ss.downCumulativePercent, ss.downCumulativePercent/ss.downDays)
	fmt.Printf("- equal: %v\n", ss.equalCount)
}

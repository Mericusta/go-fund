package observer

import (
	"encoding/json"
	"fmt"
	"go-fund/global"
	"go-fund/model"
	"go-fund/spider"
	"go-fund/spider/tushare.pro"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Mericusta/go-stp"
)

var (
	observeStockListFileRelPath string = "markdown/note/investment/stock/observe_list.json"
)

type observeData struct {
	ObserveSlice []*observeStock `json:"observe_slice"`
	AbandonSlice []*observeStock `json:"abandon_slice"`
}

type observeStock struct {
	OS_name string `json:"os_name"`
	OS_code string `json:"os_code"`
}

func (os *observeStock) Name() string { return os.OS_name }
func (os *observeStock) Code() string { return os.OS_code }

// AppendStockToObserveStockList 添加股票到观察列表中
func AppendStockToObserveStockList(name, code string) {
	observeStockListFilePath := filepath.Join(global.PersonalDocumentPath, observeStockListFileRelPath)
	err := stp.WriteFileByOverwriting(observeStockListFilePath, func(oldContentBytes []byte) ([]byte, error) {
		duplicated, observeData := false, &observeData{}
		if len(oldContentBytes) != 0 {
			err := json.Unmarshal(oldContentBytes, observeData)
			if err != nil {
				return nil, err
			}
			stp.NewArray(observeData.ObserveSlice).ForEach(func(v *observeStock, i int) {
				if !duplicated {
					duplicated = v.OS_name == name || v.OS_code == code
				}
			})
		}
		if !duplicated {
			observeData.ObserveSlice = append(observeData.ObserveSlice, &observeStock{
				OS_name: name, OS_code: code,
			})
		}
		newContentBytes, err := json.Marshal(observeData)
		if err != nil {
			return nil, err
		}
		return newContentBytes, nil
	})
	if err != nil {
		panic(err)
	}
}

// AbandonStockFromObserveList 从观察列表中移除股票
func AbandonStockFromObserveList(name, code string) {
	observeStockListFilePath := filepath.Join(global.PersonalDocumentPath, observeStockListFileRelPath)
	err := stp.WriteFileByOverwriting(observeStockListFilePath, func(oldContentBytes []byte) ([]byte, error) {
		duplicated, observeData := false, &observeData{}
		if len(oldContentBytes) != 0 {
			err := json.Unmarshal(oldContentBytes, observeData)
			if err != nil {
				return nil, err
			}
			stp.NewArray(observeData.AbandonSlice).ForEach(func(v *observeStock, i int) {
				if !duplicated {
					duplicated = v.OS_name == name || v.OS_code == code
				}
			})
		}
		if !duplicated {
			observeData.AbandonSlice = append(observeData.AbandonSlice, &observeStock{
				OS_name: name, OS_code: code,
			})
		}
		index := stp.NewArray(observeData.ObserveSlice).FindIndex(func(v *observeStock, i int) bool {
			return v.OS_name == name && v.OS_code == code
		})
		if index != -1 {
			observeData.ObserveSlice = append(observeData.ObserveSlice[:index], observeData.ObserveSlice[index+1:]...)
		}
		newContentBytes, err := json.Marshal(observeData)
		if err != nil {
			return nil, err
		}
		return newContentBytes, nil
	})
	if err != nil {
		panic(err)
	}
}

// ClearObserveStockList 清空待观察列表中的股票
func ClearObserveStockList() {
	observeStockListFilePath := filepath.Join(global.PersonalDocumentPath, observeStockListFileRelPath)
	err := stp.WriteFileByOverwriting(observeStockListFilePath, func(oldContentBytes []byte) ([]byte, error) {
		observeData := &observeData{}
		if len(oldContentBytes) != 0 {
			err := json.Unmarshal(oldContentBytes, observeData)
			if err != nil {
				return nil, err
			}
		}
		if len(observeData.ObserveSlice) > 0 {
			stp.NewArray(observeData.ObserveSlice).ForEach(func(v *observeStock, i int) {
				abandonIndex := stp.NewArray(observeData.AbandonSlice).FindIndex(func(_v *observeStock, i int) bool {
					return _v.OS_name == v.OS_name || _v.OS_code == v.OS_code
				})
				if abandonIndex == -1 {
					observeData.AbandonSlice = append(observeData.AbandonSlice, v)
				}
			})
			observeData.ObserveSlice = nil
		}
		newContentBytes, err := json.Marshal(observeData)
		if err != nil {
			return nil, err
		}
		return newContentBytes, nil
	})
	if err != nil {
		panic(err)
	}
}

// DownloadObserveStockDailyData 下载待观察股票的每日行情数据
func DownloadObserveStockDailyData(downloadDates ...time.Time) {
	observeStockBriefList := LoadObserveStockBriefList()
	if len(observeStockBriefList) == 0 {
		return
	}
	// 查找待观察股票的每日数据
	spider.MultiDownloadStockDailyData(observeStockBriefList, downloadDates...)
}

// LoadObserveStockBriefList 加载观察股票列表
func LoadObserveStockBriefList() []model.StockBriefData {
	var observeStockBriefList []model.StockBriefData
	observeStockListFilePath := filepath.Join(global.PersonalDocumentPath, observeStockListFileRelPath)
	observeStockSlice, err := stp.ReadFile(observeStockListFilePath, func(b []byte) ([]model.StockBriefData, error) {
		observeData := &observeData{}
		err := json.Unmarshal(b, observeData)
		if err != nil {
			return nil, err
		}
		stp.NewArray(observeData.ObserveSlice).ForEach(func(v *observeStock, i int) {
			observeStockBriefList = append(observeStockBriefList, v)
		})
		return observeStockBriefList, nil
	})
	if err != nil {
		panic(err)
	}
	return observeStockSlice
}

type MAPView struct {
	Value float64
	View  string
}

func (v *MAPView) Less() {

}

// CalculateObserverStockMAData 计算待观察股票的 MA 数据
func CalculateObserverStockMAData(calculateDates ...time.Time) {
	for _, calculateDate := range calculateDates {
		fmt.Println("- observer calculate stock MA data on", calculateDate.Format(tushare.TradeDateLayout()))
		// TODO: 这里需要保证 daily data 里面的数据没有间隔，保证交易日数据是连续的
		stockDailyDataMap := spider.MultiLoadStockDailyData(LoadObserveStockBriefList())
		for _, stockDailyData := range stockDailyDataMap {
			calculateIndex := stp.NewArray(stockDailyData).FindIndex(func(v model.StockDailyData, i int) bool {
				return v.Date().Format(tushare.TradeDateLayout()) == calculateDate.Format(tushare.TradeDateLayout())
			})
			if calculateIndex == -1 {
				fmt.Println("\t- can not find calculate date daily data", calculateDate.Format(tushare.TradeDateLayout()))
				return
			}
			dailyDataCount := len(stockDailyData)
			ma5Total, ma10Total, ma20Total := float64(0), float64(0), float64(0)
			for index := 0; index < 20; index++ {
				relIndex := calculateIndex + index
				if relIndex < 0 || relIndex >= dailyDataCount {
					fmt.Println("\t\t- not enough stock daily data")
					break
				}
				dailyData := stockDailyData[relIndex]
				switch {
				case index < 5:
					ma5Total += dailyData.Close()
					ma10Total += dailyData.Close()
					ma20Total += dailyData.Close()
				case index < 10:
					ma10Total += dailyData.Close()
					ma20Total += dailyData.Close()
				default:
					ma20Total += dailyData.Close()
				}
			}
			// fmt.Println("\t\t- date", calculateDate.Format(tushare.TradeDateLayout()), "ma5", ma5Total/5)
			// fmt.Println("\t\t- date", calculateDate.Format(tushare.TradeDateLayout()), "ma10", ma10Total/10)
			// fmt.Println("\t\t- date", calculateDate.Format(tushare.TradeDateLayout()), "ma20", ma20Total/20)
			MAPViewSlice := make([]*MAPView, 0, 6)
			MAPViewSlice = append(MAPViewSlice,
				&MAPView{Value: ma5Total / 5, View: "5"},
				&MAPView{Value: ma10Total / 10, View: "10"},
				&MAPView{Value: ma20Total / 20, View: "20"},
				&MAPView{Value: stockDailyData[calculateIndex].Low(), View: "L"},
				&MAPView{Value: stockDailyData[calculateIndex].Close(), View: "P"},
				&MAPView{Value: stockDailyData[calculateIndex].High(), View: "H"},
			)
			sort.Slice(MAPViewSlice, func(i, j int) bool { return MAPViewSlice[i].Value < MAPViewSlice[j].Value })
			viewSlice := make([]string, 0, 6)
			stp.NewArray(MAPViewSlice).ForEach(func(v *MAPView, i int) { viewSlice = append(viewSlice, v.View) })
			fmt.Println("\t- calculate stock", stockDailyData[calculateIndex].Name(), "MA data", "Change", stockDailyData[calculateIndex].ChangePercent(), "MA/P", strings.Join(viewSlice, "<"), "Table", stockDailyData[calculateIndex].ChangePercent(), "|", strings.Join(viewSlice, "<"))
		}
	}
}

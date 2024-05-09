package observer

import (
	"encoding/json"
	"go-fund/global"
	"go-fund/model"
	"path/filepath"

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

// AppendStockToObserveList 添加股票到观察列表中
func AppendStockToObserveList(name, code string) {
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

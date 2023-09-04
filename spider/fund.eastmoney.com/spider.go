package fundeastmoney

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-fund/global"
	"os"
	"path/filepath"

	"github.com/Mericusta/go-stp"
	"github.com/PuerkitoBio/goquery"
)

var (
	url             string = "http://fund.eastmoney.com"
	stockETFListUrl string = url + "/cnjy_jzzzl.html"
)

type FE_StockETFBriefData struct {
	FE_Code string `json:"code"`
	FE_Name string `json:"name"`
}

func (sbd *FE_StockETFBriefData) Code() string { return sbd.FE_Code }
func (sbd *FE_StockETFBriefData) Name() string { return sbd.FE_Name }

func DownloadStockETFSlice() []*FE_StockETFBriefData {
	fmt.Printf("- spider download stock ETF brief data from %v\n", url)

	codeNameMap := make(map[string]string)
	FEStockETFBriefSlice := make([]*FE_StockETFBriefData, 0, 1024)

	resp, err := global.HTTPClient.R().Get(stockETFListUrl)
	if err != nil {
		panic(err)
	}
	content := resp.Body()

	_content, err := stp.ToUtf8[stp.GBK](content)
	if err != nil {
		panic(err)
	}

	contentReader := bytes.NewReader(_content)
	doc, err := goquery.NewDocumentFromReader(contentReader)
	if err != nil {
		panic(err)
	}

	tableNodeTR := doc.Find("#oTable tr[id]")
	if tableNodeTR.Length() < 1 {
		panic(tableNodeTR.Length())
	}

	tableNodeTR.Each(func(i int, s *goquery.Selection) {
		trNodeCheckbox := s.Find("input[id]")
		if trNodeCheckbox.Length() < 1 {
			panic(trNodeCheckbox.Length())
		}
		code, has := trNodeCheckbox.Attr("id")
		if !has {
			panic(has)
		}
		trNodeNobrA := s.Find("nobr").Find("a").Get(0)
		if trNodeNobrA == nil {
			panic(trNodeNobrA)
		}
		name := trNodeNobrA.FirstChild.Data
		if _, has := codeNameMap[code]; !has {
			codeNameMap[code] = name
			FEStockETFBriefSlice = append(FEStockETFBriefSlice, &FE_StockETFBriefData{
				FE_Code: code, FE_Name: name,
			})
		}
	})

	fmt.Printf("- spider download stock ETF brief data count %v\n", len(codeNameMap))
	return FEStockETFBriefSlice
}

// func convertFEStockETFBriefDataSlice(stockETFBriefDataSlice []searcher.StockBriefData) []*FE_StockETFBriefData {
// 	stockETFCodeBriefMap := make(map[string]searcher.StockBriefData)
// 	for _, stockETFBrief := range stockETFBriefDataSlice {
// 		stockETFCodeBriefMap[stockETFBrief.Code()] = stockETFBrief
// 	}
// 	keySlice := stp.Key(stockETFCodeBriefMap)
// 	sort.Strings(keySlice)

// 	FEStockBriefDataSlice := make([]*FE_StockETFBriefData, 0, len(stockETFBriefDataSlice))
// 	for _, key := range keySlice {
// 		FEStockBriefDataSlice = append(FEStockBriefDataSlice, stockETFCodeBriefMap[key].(*FE_StockETFBriefData))
// 	}

// 	return FEStockBriefDataSlice
// }

// func revertFEStockETFBriefDataSlice(FEStockBriefSlice []*FE_StockETFBriefData) []searcher.StockBriefData {
// 	stockBriefDataSlice := make([]searcher.StockBriefData, 0, len(FEStockBriefSlice))
// 	for _, d := range FEStockBriefSlice {
// 		stockBriefDataSlice = append(stockBriefDataSlice, d)
// 	}
// 	return stockBriefDataSlice
// }

func SaveStockETFList(stockETFBriefSlice []*FE_StockETFBriefData) {
	fmt.Printf("- spider save stock ETF brief data to personal document %v\n", global.StockETFListRelativePath)

	// slice := convertFEStockETFBriefDataSlice(stockETFBriefSlice)
	// if len(slice) != len(stockETFBriefSlice) {
	// 	panic("length not equal")
	// }

	stockETFListPath := filepath.Join(global.PersonalDocumentPath, global.StockETFListRelativePath)
	if stp.IsExist(stockETFListPath) {
		if err := os.Remove(stockETFListPath); err != nil {
			panic(err)
		}
		fmt.Printf("\t- spider remove old stock etf list file\n")
	}
	stockETFListFile, err := os.Create(stockETFListPath)
	if err != nil {
		panic(err)
	}
	defer stockETFListFile.Close()

	b, err := json.Marshal(stockETFBriefSlice)
	if err != nil {
		panic(err)
	}

	_, err = stockETFListFile.Write(b)
	if err != nil {
		panic(err)
	}
}

func LoadStockETFList() []*FE_StockETFBriefData {
	fmt.Printf("- spider load stock etf brief data from personal document %v\n", global.StockETFListRelativePath)

	stockETFListPath := filepath.Join(global.PersonalDocumentPath, global.StockETFListRelativePath)
	if !stp.IsExist(stockETFListPath) {
		fmt.Printf("\t- stock etf brief data file %v not exists in personal document\n", global.StockETFListRelativePath)
		return nil
	}
	stockETFList, err := os.ReadFile(stockETFListPath)
	if err != nil {
		panic(err)
	}

	slice := make([]*FE_StockETFBriefData, 0, 8192)
	err = json.Unmarshal(stockETFList, &slice)
	if err != nil {
		panic(err)
	}

	return slice
}

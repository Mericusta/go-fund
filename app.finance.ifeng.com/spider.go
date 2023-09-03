package appfinanceifengcom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-fund/global"
	"go-fund/searcher"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/Mericusta/go-stp"
	"github.com/PuerkitoBio/goquery"
)

const (
	url                string = "https://app.finance.ifeng.com"
	stockListUrlFormat string = "/list/stock.php?t=hs&f=chg_pct&o=desc&p=%v"
)

type AFI_StockBriefData struct {
	AFI_Code string `json:"code"`
	AFI_Name string `json:"name"`
}

func (sbd *AFI_StockBriefData) Code() string { return sbd.AFI_Code }
func (sbd *AFI_StockBriefData) Name() string { return sbd.AFI_Name }

func DownloadStockSlice() []searcher.StockBriefData {
	fmt.Printf("- spider get stock brief data from %v\n", url)

	targetHeaderIndexMap := map[string]int{"代码": -1, "名称": -1}
	targetIndexHeaderMap := map[int]string{}
	codeNameMap := make(map[string]string)
	AFIStockBriefSlice := make([]*AFI_StockBriefData, 0, 8192)

	for page := 1; true; page++ {
		_url := fmt.Sprintf(stockListUrlFormat, page)
		resp, err := global.HTTPClient.R().Get(_url)
		if err != nil {
			panic(err)
		}

		content := resp.Body()
		contentReader := bytes.NewReader(content)
		doc, err := goquery.NewDocumentFromReader(contentReader)
		if err != nil {
			panic(err)
		}

		tableNodeTR := doc.Find(".tab01").Find("table").Find("tr")
		if tableNodeTR.Length() < 3 {
			break
		}
		tableNodeTR.Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				s.Find("th").Each(func(j int, s *goquery.Selection) {
					if index, has := targetHeaderIndexMap[s.Text()]; index == -1 && has {
						targetHeaderIndexMap[s.Text()] = j
						if title, has := targetIndexHeaderMap[j]; len(title) == 0 || !has {
							targetIndexHeaderMap[j] = s.Text()
						}
					}
				})
			} else {
				var code, name string
				tableNodeTD := s.Find("td")
				if tableNodeTD.Length() > 2 {
					tableNodeTD.Each(func(i int, s *goquery.Selection) {
						if title, has := targetIndexHeaderMap[i]; len(title) > 0 && has {
							switch title {
							case "代码":
								code = s.Text()
							case "名称":
								name = s.Text()
							}
						}
					})
					if _, has := codeNameMap[code]; !has {
						codeNameMap[code] = name
						AFIStockBriefSlice = append(AFIStockBriefSlice, &AFI_StockBriefData{
							AFI_Code: code, AFI_Name: name,
						})
					}
				}
			}
		})

		fmt.Printf("\t- handle page %v done, stock count %v\n", page, len(codeNameMap))
		time.Sleep(time.Second)
		for header := range targetHeaderIndexMap {
			targetHeaderIndexMap[header] = -1
		}
		targetIndexHeaderMap = make(map[int]string)
	}

	fmt.Printf("- spider get stock brief data count %v\n", len(codeNameMap))
	return revertAFIStockBriefDataSlice(AFIStockBriefSlice)
}

func convertAFIStockBriefDataSlice(stockBriefDataSlice []searcher.StockBriefData) []*AFI_StockBriefData {
	stockNameBriefMap := make(map[string]searcher.StockBriefData)
	for _, stockBrief := range stockBriefDataSlice {
		stockNameBriefMap[stockBrief.Code()] = stockBrief
	}
	keySlice := stp.Key(stockNameBriefMap)
	sort.Strings(keySlice)

	AFIStockBriefDataSlice := make([]*AFI_StockBriefData, 0, len(stockBriefDataSlice))
	for _, key := range keySlice {
		AFIStockBriefDataSlice = append(AFIStockBriefDataSlice, stockNameBriefMap[key].(*AFI_StockBriefData))
	}

	return AFIStockBriefDataSlice
}

func revertAFIStockBriefDataSlice(AFIStockBriefSlice []*AFI_StockBriefData) []searcher.StockBriefData {
	stockBriefDataSlice := make([]searcher.StockBriefData, 0, len(AFIStockBriefSlice))
	for _, d := range AFIStockBriefSlice {
		stockBriefDataSlice = append(stockBriefDataSlice, d)
	}
	return stockBriefDataSlice
}

func SaveStockList(stockBriefSlice []searcher.StockBriefData) {
	fmt.Printf("- spider save stock brief data to personal document %v\n", global.StockListRelativePath)

	slice := convertAFIStockBriefDataSlice(stockBriefSlice)
	if len(slice) != len(stockBriefSlice) {
		panic("length not equal")
	}

	stockListPath := filepath.Join(global.PersonalDocumentPath, global.StockListRelativePath)
	if stp.IsExist(stockListPath) {
		if err := os.Remove(stockListPath); err != nil {
			panic(err)
		}
		fmt.Printf("\t- spider remove old stock list file\n")
	}
	stockListFile, err := os.Create(stockListPath)
	if err != nil {
		panic(err)
	}
	defer stockListFile.Close()

	b, err := json.Marshal(slice)
	if err != nil {
		panic(err)
	}

	_, err = stockListFile.Write(b)
	if err != nil {
		panic(err)
	}
}

func LoadStockList() []searcher.StockBriefData {
	fmt.Printf("- spider load stock list from personal document %v\n", global.StockListRelativePath)

	stockListPath := filepath.Join(global.PersonalDocumentPath, global.StockListRelativePath)
	if !stp.IsExist(stockListPath) {
		fmt.Printf("\t- stock list %v not exists in personal document\n", global.StockListRelativePath)
		return nil
	}
	stockList, err := os.ReadFile(stockListPath)
	if err != nil {
		panic(err)
	}

	slice := make([]*AFI_StockBriefData, 0, 8192)
	err = json.Unmarshal(stockList, &slice)
	if err != nil {
		panic(err)
	}

	return revertAFIStockBriefDataSlice(slice)
}

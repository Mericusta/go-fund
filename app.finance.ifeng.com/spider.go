package appfinanceifengcom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-fund/global"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Mericusta/go-stp"
	"github.com/PuerkitoBio/goquery"
)

func GetStockList() map[string]string {
	urlFormat := "https://app.finance.ifeng.com/list/stock.php?t=hs&f=chg_pct&o=desc&p=%v"

	targetHeaderIndexMap := map[string]int{"代码": -1, "名称": -1}
	targetIndexHeaderMap := map[int]string{}
	data := make(map[string]string)

	for page := 1; true; page++ {
		url := fmt.Sprintf(urlFormat, page)
		resp, err := global.HTTPClient.R().Get(url)
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
					data[code] = name
				}
			}
		})

		fmt.Printf("- handle page %v done, stock count %v\n", page, len(data))
		time.Sleep(time.Second)
		for header := range targetHeaderIndexMap {
			targetHeaderIndexMap[header] = -1
		}
		targetIndexHeaderMap = make(map[int]string)
	}
	return data
}

var (
	stockListRelativePath  string = "markdown/note/stock/stock_list.json"
	localStockListDataPath string = "../stock_list"
)

func ConvertStockList(stockNameCodeMap map[string]string) {
	if len(stockNameCodeMap) == 0 {
		stockNameCodeMap = make(map[string]string)
		stp.ReadFileLineOneByOne(localStockListDataPath, func(s string, i int) bool {
			slice := strings.Split(s, "|")
			if len(slice) < 4 {
				return true
			}
			stockNameCodeMap[slice[3]] = slice[1]
			return true
		})
	}

	keySlice := stp.Key(stockNameCodeMap)
	sort.Strings(keySlice)

	stockSlice := make([]struct {
		Code string
		Name string
	}, 0, len(stockNameCodeMap))
	for _, key := range keySlice {
		stockSlice = append(stockSlice, struct {
			Code string
			Name string
		}{
			Code: stockNameCodeMap[key],
			Name: key,
		})
	}

	SaveStockList(stockSlice)
}

func SaveStockList(s []struct {
	Code string
	Name string
}) {
	stockListPath := filepath.Join(global.PersonalDocumentPath, stockListRelativePath)
	if stp.IsExist(stockListPath) {
		if err := os.Remove(stockListPath); err != nil {
			panic(err)
		}
	}
	stockListFile, err := os.Create(stockListPath)
	if err != nil {
		panic(err)
	}
	defer stockListFile.Close()

	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	_, err = stockListFile.Write(b)
	if err != nil {
		panic(err)
	}
}

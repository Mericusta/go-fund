package appfinanceifengcom

import (
	"bytes"
	"fmt"
	"go-fund/global"
	"time"

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
		if tableNodeTR.Length() == 0 {
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

		time.Sleep(time.Second)
		for header := range targetHeaderIndexMap {
			targetHeaderIndexMap[header] = -1
		}
		targetIndexHeaderMap = make(map[int]string)
	}

	fmt.Printf("data = %v\n", data)
	return nil
}

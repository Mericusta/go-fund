package fundeastmoney

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func Spider(fundNameCodeMap map[string]string) {
	for name, code := range fundNameCodeMap {
		date, num := GetFundInfoByCode(code)
		pngName := GetFuncChartsByCode(code)
		fmt.Printf("name %v date %v num %v png %v\n", name, date, num, pngName)
	}
}

func GetFundInfoByCode(code string) (string, string) {
	res, err := http.Get(getFundInfoUrlByCode(code))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	dataOfFundNodeSelection := doc.Find(".dataOfFund")
	date, err := dataOfFundNodeSelection.Find(".dataItem01").Find("#gz_gztime").Html()
	if err != nil {
		panic(err)
	}
	// fmt.Printf("date %v\n", date)

	num, err := dataOfFundNodeSelection.Find(".dataItem02").Find(".dataNums").Children().First().Html()
	if err != nil {
		panic(err)
	}
	// fmt.Printf("num %v\n", num)
	return date, num
}

func GetFuncChartsByCode(code string) string {
	res, err := http.Get(getFundChartUrlByCode(code))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	contentBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	pngName := fmt.Sprintf("%v_%v.png", code, time.Now().Format("20060102150405"))
	err = ioutil.WriteFile(pngName, contentBytes, 0644)
	if err != nil {
		panic(err)
	}

	return pngName
}

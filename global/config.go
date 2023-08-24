package global

import (
	"os"

	"github.com/go-resty/resty/v2"
)

var (
	PersonalDocumentPath string
	HTTPClient           *resty.Client
	FundNameCodeMap      map[string]string
	FundCodeNameMap      map[string]string
	Date1                string = "20071016" // 上证指数涨至历史第一高点 6124.04
	Date2                string = "20210228" // 上证指数涨至历史第二高点 3731.69
)

func init() {
	PersonalDocumentPath = os.Getenv("PD")
	HTTPClient = resty.New()

	FundNameCodeMap = make(map[string]string)
	FundCodeNameMap = make(map[string]string)

	FundNameCodeMap["招商国证生物医药指数(LOF)A"] = "161726"
	FundNameCodeMap["招商中证白酒指数(LOF)A"] = "161725"
}
